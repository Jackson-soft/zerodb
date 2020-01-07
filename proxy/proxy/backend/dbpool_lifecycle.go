package backend

import (
	"sync"
	"time"

	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"github.com/pkg/errors"
)

const (
	Up = iota
	Cleaning
	Cleaned
	Down

	InitConnCount           = 16
	DefaultMaxConnNum       = 1024
	PingPeroid        int64 = 30 // sec
)

// cobar -> mysql 数据库连接池
type DBPool struct {
	sync.RWMutex

	status int8

	hostGroup string
	poolType  string

	addr     string
	user     string
	password string
	dbPool   string

	maxConnNum  int
	InitConnNum int

	// 活的空闲连接 + 死的连接 + 被借走的连接 = maxConnNum
	// 有效后端连接 = 活的空闲连接 + 被借走的连接 = maxConnNum - 死的连接
	// 没有与MySQL激活，只是一个死亡连接，或者叫一个壳
	deadConns chan *MysqlConn
	// 已经与MySQL激活的的空闲连接
	idleConns chan *MysqlConn
	// 正在被使用的连接
	ActiveConns sync.Map // map[int64]*MysqlConn

	// 心跳
	detectConn    *MysqlConn // 单独用于心跳的连接
	NeedHeartbeat bool

	maxIdleLiveTime int64 // sec
}

func InitDBPool(addr string, hostGroup string, poolType string, user string, password string, dbName string, initConnNum int, maxConnNum int) (*DBPool, error) {
	var err error
	pool := new(DBPool)
	pool.hostGroup = hostGroup
	pool.poolType = poolType
	pool.addr = addr
	pool.user = user
	pool.password = password
	pool.dbPool = dbName
	pool.maxIdleLiveTime = 10
	pool.NeedHeartbeat = true

	if 0 <= maxConnNum {
		pool.maxConnNum = maxConnNum
		pool.InitConnNum = initConnNum
	} else {
		pool.maxConnNum = DefaultMaxConnNum
		pool.InitConnNum = InitConnCount
	}

	//modified by huajia
	pool.detectConn, err = pool.createPooledMySQLConn()
	if err == nil {
		//pool.CloseResource()
		//return nil, err
		pool.deadConns = make(chan *MysqlConn, pool.maxConnNum)
		pool.idleConns = make(chan *MysqlConn, pool.maxConnNum)
		for i := 0; i < pool.maxConnNum; i++ {
			if i < pool.InitConnNum {
				conn, err := pool.createPooledMySQLConn()
				if err != nil {
					pool.CloseIdleBackendConns()
					pool.ClosePingDetectConn()
					pool.CloseResource()
					return nil, err
				}
				conn.pushTimestamp = time.Now().Unix()
				pool.idleConns <- conn
			} else {
				conn := NewConnCache()
				pool.deadConns <- conn
			}
		}
		pool.status = Up
	} else {
		pool.detectConn = nil
		pool.status = Down
	}
	//end modified

	go pool.WatermarkCheckTask()
	go pool.IdleConnEvictTask()
	go pool.IdleConnSupplyTask()
	go pool.IdleConnKeepAliveTask()

	return pool, nil
}

//added by huajia
func (pool *DBPool) ReConnect() error {
	defer pool.Unlock()
	pool.Lock()
	if pool.status != Down {
		return nil
	}
	pool.deadConns = make(chan *MysqlConn, pool.maxConnNum)
	pool.idleConns = make(chan *MysqlConn, pool.maxConnNum)

	detectConn, err := pool.createPooledMySQLConn()
	if err != nil {
		pool.CloseIdleBackendConns()
		pool.ClosePingDetectConn()
		pool.CloseResource()
		return err
	}
	pool.detectConn = detectConn

	pool.status = Up

	return nil
}
//end added

// 重新初始化，防止chan被delete了
func (pool *DBPool) ReInit() error {
	defer pool.Unlock()
	pool.Lock()
	if pool.status != Cleaned {
		return nil
	}

	pool.deadConns = make(chan *MysqlConn, pool.maxConnNum)
	pool.idleConns = make(chan *MysqlConn, pool.maxConnNum)

	detectConn, err := pool.createPooledMySQLConn()
	if err != nil {
		pool.CloseIdleBackendConns()
		pool.ClosePingDetectConn()
		pool.CloseResource()
		return err
	}
	pool.detectConn = detectConn

	pool.status = Up

	return nil
}

// 预热  TODO nanxing 需要修改逻辑
func (pool *DBPool) WarmUpPool() error {
	defer pool.Unlock()
	pool.Lock()

	return nil
}

func (pool *DBPool) CloseResource() error {
	pool.Lock()
	pool.deadConns = nil
	pool.idleConns = nil
	pool.status = Down
	pool.Unlock()
	return nil
}

func (pool *DBPool) GetName() string {
	return pool.hostGroup + "-" + pool.poolType + ":[" + pool.addr + "]"
}

func (pool *DBPool) IdleConnCount() int {
	pool.RLock()
	defer pool.RUnlock()
	return len(pool.idleConns)
}

func (pool *DBPool) CloseActiveBackendConns() {
	pool.ActiveConns.Range(func(key, value interface{}) bool {
		if value != nil {
			if v, ok := value.(*MysqlConn); ok {
				pool.returnConnToDeadChan(v)
			}
		}
		return true
	})
}

func (pool *DBPool) CloseIdleBackendConns() {
	defer pool.Unlock()
	pool.Lock()

	pool.status = Cleaning
	close(pool.idleConns)
	for conn := range pool.idleConns {
		pool.returnConnToDeadChan(conn)
	}
	pool.status = Cleaned
}

func (pool *DBPool) ClosePingDetectConn() {
	pool.detectConn.Close()
}

func (pool *DBPool) Ping() error {
	var err error
	if pool.detectConn == nil {
		pool.detectConn, err = pool.createPooledMySQLConn()
		if err != nil {
			pool.detectConn = nil
			return err
		}
	}
	err = pool.detectConn.PingTimeout()
	if err != nil {
		pool.detectConn.Close()
		pool.detectConn = nil
		return err
	}
	return nil
}

//added by huajia
func (pool *DBPool) IsActive() bool {
	return pool.detectConn != nil
}
//end added

func (pool *DBPool) createPooledMySQLConn() (*MysqlConn, error) {
	co := NewConnCache()

	if err := co.ConnectMysql(pool.addr, pool.user, pool.password, pool.dbPool); err != nil {
		return nil, errors.WithMessage(err, "connecting mysql "+pool.addr+" failed")
	}

	return co, nil
}

func (pool *DBPool) returnConnToDeadChan(co *MysqlConn) error {
	if co != nil {
		co.Close()
		if pool.deadConns != nil {
			select {
			case pool.deadConns <- co:
				return nil
			default:
				return nil
			}
		}
	}
	return nil
}

/**
既然是已经放回pool的连接，有可能是用过的，也有可能是新创建的。
对于用过的并且归还的连接，需要一个重置的步骤
*/
func (pool *DBPool) resetConn(co *MysqlConn) error {
	var err error
	//reuse Connection
	if co.IsInTransaction() {
		//we can not reuse a connection in transaction status
		err = co.Rollback()
		if err != nil {
			return err
		}
	}

	if !co.IsAutoCommit() {
		//we can not  reuse a connection not in autocomit
		_, err = co.exec("set autocommit = 1")
		if err != nil {
			return err
		}
	}

	//connection may be set names early
	//we must use default utf8
	if co.GetCharset() != mysql.DEFAULT_CHARSET {
		err = co.SetCharset(mysql.DEFAULT_CHARSET, mysql.DEFAULT_COLLATION_ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pool *DBPool) popConn() (*MysqlConn, error) {
	var co *MysqlConn
	var err error

	if pool.idleConns == nil || pool.deadConns == nil {
		//test huajia
		glog.Glog.Errorf("db pool %s is closed....", pool.hostGroup)
		return nil, errcode.ErrDBPoolClosed
	}
	co = pool.getConnFromIdle(pool.idleConns)
	if co == nil {
		glog.Glog.Warnf("run out of idle conn in pool[%s] active: %d, idle: %d, dead: %d, try borrow and activate conn from dead.",
			pool.GetName(),
			pool.ActiveSize(),
			pool.RemainIdleSize(),
			pool.RemainDeadSize())
		co, err = pool.getConnFromDead()
		if err != nil {
			return nil, err
		}
	}

	err = pool.resetConn(co)
	if err != nil {
		pool.returnConnToDeadChan(co)
		return nil, err
	}

	return co, nil
}

func (pool *DBPool) getConnFromIdle(idleConns chan *MysqlConn) *MysqlConn {
	var co *MysqlConn
	var err error
	for 0 < len(idleConns) {
		co = <-idleConns
		// TODO nanxing 重新获取的时候，是否需要测试存活性，对性能非常有影响
		if co != nil && PingPeroid < time.Now().Unix()-co.pushTimestamp {
			err = co.Ping()
			if err != nil {
				pool.returnConnToDeadChan(co)
				co = nil
			}
		}
		if co != nil {
			break
		}
	}
	return co
}

func (pool *DBPool) getConnFromDead() (*MysqlConn, error) {
	var co *MysqlConn
	var err error

	co = <-pool.deadConns

	err = co.ConnectMysql(pool.addr, pool.user, pool.password, pool.dbPool)
	if err != nil {
		pool.returnConnToDeadChan(co)
		return nil, err
	}
	return co, nil
}

func (pool *DBPool) returnConnToLiveChan(co *MysqlConn, err error) {
	if co == nil {
		return
	}
	if pool.idleConns == nil {
		co.Close()
		return
	}
	if err != nil {
		pool.returnConnToDeadChan(co)
		return
	}
	co.pushTimestamp = time.Now().Unix()
	select {
	case pool.idleConns <- co:
		return
	default:
		pool.returnConnToDeadChan(co)
		return
	}
}
