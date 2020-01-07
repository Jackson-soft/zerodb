package server

import (
	"net"
	"net/http"
	"sync"
	"time"

	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	"git.2dfire.net/zerodb/keeper/pkg/models"

	_ "net/http/pprof"

	"runtime/debug"

	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/keeper/metrics"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/timer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	conf *models.Config

	store        *etcdtool.Store
	proxyClients sync.Map //proxy rpc client map
	agentClients sync.Map //agent rpc client map
	//timer        *time.Ticker //全局定时器
	tTimer *timer.Timer
	//done         chan bool
	pushFailed bool //推送配置如果有失败是否继续
}

func NewServer(conf *models.Config) (*Server, error) {
	s := new(Server)
	s.conf = conf
	var err error
	s.store, err = etcdtool.NewStore(s.conf.Persistence)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//s.done = make(chan bool)
	//s.timer = time.NewTicker(Heartbeat)
	s.tTimer = timer.NewTimer(time.Second, 5)

	if err := s.newRPCCliByDB(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Run() error {
	if _, err := s.tTimer.Register(timer.Repetition, HeartbeatCheckInterval, func(arg interface{}) { s.checkHeartbeat() }, nil); err != nil {
		return err
	}

	go s.startWeb()

	go s.startMetrics()

	go s.startRPC()

	if s.conf.DebugMode {
		go s.starDebugServer()
	}

	select {}
}

func (s *Server) startRPC() {
	ln, err := net.Listen("tcp", s.conf.RPCServer)
	if err != nil {
		glog.GLog.Fatalln(err)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryHandleFunc),
	}
	rpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		grpc_recovery.UnaryServerInterceptor(opts...),
	),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		))
	keeper.RegisterKeeperServer(rpcServer, s)
	reflection.Register(rpcServer)
	glog.GLog.Infof("keeper RPC server running on %s", s.conf.RPCServer)
	if err := rpcServer.Serve(ln); err != nil {
		glog.GLog.Fatalln(err)
	}
}

func (s *Server) startWeb() {
	engine := s.creatHander()
	glog.GLog.Infof("start web server on %v", s.conf.WebServer)
	if err := http.ListenAndServe(s.conf.WebServer, engine); err != nil {
		glog.GLog.Fatalln(err)
	}
}

//startMetrics 方法启动Metrics服务
func (s *Server) startMetrics() {
	if err := metrics.Run(s.conf.MetricsAddr); err != nil {
		glog.GLog.Fatalln(err)
	}
}

func (s *Server) Close() {
	s.tTimer.Stop()
}

//start pprof debug server
func (s *Server) starDebugServer() {
	glog.GLog.Infof("start debug server on %v", s.conf.DebugServer)
	if err := http.ListenAndServe(s.conf.DebugServer, nil); err != nil {
		glog.GLog.Fatalln(err)
	}
}

func recoveryHandleFunc(p interface{}) error {
	if p != nil {
		glog.GLog.Errorf("rpc panic %s", debug.Stack())
	}
	return nil
}
