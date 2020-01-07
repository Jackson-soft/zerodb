package frontend

import (
	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"net"
	"sync"
)

type ProxyEngine struct {
	// auth
	addr     string
	user     string
	password string

	// logging
	slowLogTime int
	LogSQL      bool

	// config
	Conf *config.Config

	// sharding
	HostGroupNodes    sync.Map //map[string]*backend.HostGroupNode
	schemaNodes       sync.Map //map[string]*backend.SchemaNode
	HostGroupClusters sync.Map // map[string]*config.HostGroupCluster

	router *Router

	// net
	listener net.Listener
	running  bool
	state    string

	// lock
	confApplyLock sync.Mutex

	//added by huajia, for cobar
	proxyClusters []ProxyCluster
}

func NewProxyEngine(addr string, charset string, logSQL bool) (*ProxyEngine, error) {
	s := new(ProxyEngine)

	s.addr = addr
	s.LogSQL = logSQL

	if len(charset) == 0 {
		charset = mysql.DEFAULT_CHARSET //utf8mb4
	}
	cid, ok := mysql.CharsetIds[charset]
	if !ok {
		return nil, errcode.BuildError(errcode.InvalidCharset, charset)
	}
	//change the default charset
	mysql.DEFAULT_CHARSET = charset
	mysql.DEFAULT_COLLATION_ID = cid
	mysql.DEFAULT_COLLATION_NAME = mysql.Collations[cid]
	glog.Glog.Infof("MySQL engine using charset:[%s]", charset)
	return s, nil
}

func (e *ProxyEngine) ServeInBlocking(listened func()) {
	if e.running {
		panic("proxy engine is running, do not serve repeatedly")
	}

	e.running = true
	e.warmUpHostGroups()

	var err error
	netProto := "tcp"

	e.listener, err = net.Listen(netProto, e.addr)

	if err != nil {
		panic(err)
	}

	glog.Glog.Infof("Start listening at %s", e.addr)

	glog.Glog.Infof("MySQL engine is on.")
	for e.running {
		conn, err := e.listener.Accept()
		if err != nil {
			glog.Glog.Errorf("accpet error: %v", err)
			continue
		}
		listened()

		go e.newConn(conn) // TCP读写
	}
}

func (e *ProxyEngine) CloseProxyEngine() {
	if !e.running {
		return
	}
	e.running = false
	glog.Glog.Infof("Shutting down MySQL engine.")
	e.cleanUpFrontendConnections()
	e.cleanUpBackendConnections()
	glog.Glog.Infof("MySQL engine is off.")
}

func (e *ProxyEngine) cleanUpFrontendConnections() {
	if e.listener != nil {
		e.listener.Close()
	}
	glog.Glog.Infof("Frontend connections have been cleaned up.")
}

// clean up all backend connections except the ping detector
func (e *ProxyEngine) cleanUpBackendConnections() {
	e.HostGroupNodes.Range(func(key, value interface{}) bool {
		if hostGroup, ok := value.(*backend.HostGroupNode); ok {
			for _, wPool := range hostGroup.Write {
				wPool.CloseActiveBackendConns()
				wPool.CloseIdleBackendConns()
			}

			for _, rPool := range hostGroup.Read {
				rPool.CloseActiveBackendConns()
				rPool.CloseIdleBackendConns()
			}
		}

		return true
	})
	glog.Glog.Infof("Backend connections have been cleaned up.")
}

func (e *ProxyEngine) GetState() string {
	return e.state
}
func (e *ProxyEngine) SetState(state string) {
	e.state = state
}

// 主要是初始化HostGroup的连接池，从而激活idle连接的补充
func (e *ProxyEngine) warmUpHostGroups() {
	e.HostGroupNodes.Range(func(key, value interface{}) bool {
		if hostGroup, ok := value.(*backend.HostGroupNode); ok {
			for _, wPool := range hostGroup.Write {
				wPool.ReInit()
			}

			for _, rPool := range hostGroup.Read {
				rPool.ReInit()
			}
		}

		return true
	})
}
