package frontend

import (
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var baseConnId uint32 = 0

//client <-> proxy
type FrontendConn struct {
	sync.Mutex

	proxy *ProxyEngine

	pkg        *mysql.PacketIO
	c          net.Conn
	capability uint32
	collation  mysql.CollationId
	charset    string
	user       string
	logicDb    string // proxy中的数据库名称
	physicalDb string // 真实MySQL中的数据库名称，有分库后缀
	salt       []byte

	status uint16 // 集成了许多种类的状态

	// transaction 事务
	// 用于在执行事务SQL时锁定 BEGIN -> COMMIT 之间的ShardingNode/BackendConn的1对1关系
	txConns map[int]*backend.BackendConn // map[TableIndex]*backend.BackendConn

	closed       bool
	lastInsertId int64
	affectedRows int64
	stmtId       uint32
	stmts        map[uint32]*Stmt //prepare相关,client端到proxy的stmt

	// performance 性能
	handlerArena *util.FastAllocator
	executeArena *util.FastYetSafeAllocator

	// stat 统计
	connectionId uint32
	bornTime     int64 // 产生时间戳
}

// 将普通的TCP连接包装成FrontendConn
func NewFrontendConn(co net.Conn, s *ProxyEngine) *FrontendConn {
	c := new(FrontendConn)
	tcpConn := co.(*net.TCPConn)

	tcpConn.SetNoDelay(false)
	tcpConn.SetReadBuffer(16 * 1024)
	tcpConn.SetWriteBuffer(16 * 1024)

	c.c = tcpConn

	c.pkg = mysql.NewPacketIO(tcpConn)
	c.proxy = s

	c.pkg.Sequence = 0

	c.connectionId = atomic.AddUint32(&baseConnId, 1)

	c.status = mysql.SERVER_STATUS_AUTOCOMMIT

	// TODO nanxing 上线需要去掉
	c.status |= mysql.SERVER_STATUS_PERMIT_MULTI_ROUTE

	c.salt, _ = mysql.RandomBuf(20)

	c.txConns = make(map[int]*backend.BackendConn)

	c.closed = false

	c.charset = mysql.DEFAULT_CHARSET
	c.collation = mysql.DEFAULT_COLLATION_ID

	c.stmtId = 0
	c.stmts = make(map[uint32]*Stmt)

	c.bornTime = time.Now().UnixNano()

	c.handlerArena = util.NewFastAllocator(16 * 1024) // 16kb
	c.executeArena = util.NewFastYetSafeAllocator(8 * 1024)

	return c
}

func (fc *FrontendConn) Start() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			glog.Glog.Errorf("run error. err:%v, stack:%s", err, string(buf))
		}

		fc.Close(false)
	}()

	glog.Glog.Infof("connection comes, connectionID:[%v], schema:[%s], IP:[%s]", fc.connectionId, fc.logicDb, fc.c.RemoteAddr().String())
	monitor.Monitor.IncrFrontendConns(fc.logicDb)

	for {
		fc.handlerArena.Reset()
		fc.executeArena.Reset()

		data, err := fc.readPacket()

		if err != nil {
			return
		}

		// proxy service可能在重启，所以直接退出协程
		if !fc.proxy.running {
			return
		}

		if err := fc.dispatch(data); err != nil {
			glog.Glog.Errorf("connectionID:[%v], command executing failed. error: %s",
				fc.connectionId,
				err.Error(),
			)
			fc.writeError(err)
		}

		if fc.closed {
			return
		}

		fc.pkg.Sequence = 0
	}
}

// actively: 是否是主动由client端调用
func (fc *FrontendConn) Close(actively bool) error {
	if fc.closed {
		return nil
	}
	fc.c.Close()

	monitor.Monitor.DecrFrontendConns(fc.logicDb)

	glog.Glog.Infof("connection leaves, connectionID:[%v], schema:[%s], IP:[%s], survival time: %vms, active:%v",
		fc.connectionId,
		fc.logicDb,
		fc.GetRemoteIP(),
		(time.Now().UnixNano()-fc.bornTime)/int64(time.Millisecond),
		actively)

	fc.closed = true
	return nil
}

func (fc *FrontendConn) GetRemoteIP() string {
	return fc.c.RemoteAddr().String()
}

func (fc *FrontendConn) permitMultiRoute(status uint16) bool {
	return status&mysql.SERVER_STATUS_PERMIT_MULTI_ROUTE > 0
}
