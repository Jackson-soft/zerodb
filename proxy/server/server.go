package server

import (
	"net"
	"time"

	"sync"

	"runtime/debug"

	"net/http"

	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/proxy/metrics"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/pkg/models"
	"git.2dfire.net/zerodb/proxy/proxy/frontend"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "net/http/pprof"

	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

var (
	dstopCh  = make(chan struct{})
	rstopCh  = make(chan struct{})
	hbstopCh = make(chan struct{})
)

type Server struct {
	sync.Mutex

	status string
	conf   *models.Config

	shardConfig *ShardConfig
	keeperCli   keeper.KeeperClient
	ticker      *time.Ticker

	proxyEngine *frontend.ProxyEngine
	// stop service and datasource switch
	hostGroupSwtRejectMap sync.Map
	hostGroupPingFailMap  sync.Map
	getLost               bool // proxy连不上keeper总部
}
type ChServer struct {
	//check healthy server
	ln   net.Listener
	addr string
}

//NewServer 函数创建一个新的Server实例
func NewServer(conf *models.Config) (*Server, error) {
	s := new(Server)
	//在这里对Server初始化

	s.conf = conf

	s.shardConfig = NewShardConf()

	// 在集群发生数据源切换或者更新配置时，严格来说，新的proxy server是不能启动的
	s.status = StatusReady

	//conn, err := grpc.Dial(s.conf.KeeperAddr, grpc.WithInsecure())
	//使用go-grpc-middleware 完成grpc client端重试和连接超时
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeOut)
	defer cancel()
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(MaxRetryNum),
		grpc_retry.WithPerRetryTimeout(time.Second * 10),
	}
	conn, err := grpc.DialContext(ctx, s.conf.KeeperAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	/*
		ipAddr, err := net.ResolveIPAddr("ip", s.conf.KeeperAddr)
		if err != nil {
			return nil, errors.Errorf("can not resolve zeroKeeper's %v ip address", s.conf.KeeperAddr)
		}
		conn, err := grpc.DialContext(ctx, ipAddr.String() + ":5004", grpc.WithInsecure())
	*/
	if err != nil {
		return nil, errors.Wrapf(err, "connect zeroKeeper %v", s.conf.KeeperAddr)
	}

	s.keeperCli = keeper.NewKeeperClient(conn)

	s.proxyEngine, err = frontend.NewProxyEngine(s.conf.ProxyServer, s.conf.Charset, s.conf.LogSQL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.ticker = time.NewTicker(HearbeatTime)

	return s, nil
}

//ServeRPC 方法启动RPC服务
func (s *Server) ServeRPC() error {
	ln, err := net.Listen("tcp", s.conf.RPCServer)
	if err != nil {
		return errors.WithStack(err)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryHandleFunc),
	}
	grpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		grpc_recovery.UnaryServerInterceptor(opts...),
	),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		))
	//grpcServer := grpc.NewServer()
	proxy.RegisterProxyServer(grpcServer, s)
	reflection.Register(grpcServer)

	return grpcServer.Serve(ln)
}

//ServeMetrics 方法启动Metrics服务
func (s *Server) ServeMetrics() {
	metrics.Run(s.conf.MetricsAddr)
}

//Run 主循环
func (s *Server) Run() error {
	var err error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				glog.Glog.Errorf("keeper rpc service panic %s", debug.Stack())
			}
		}()
		if err = s.ServeRPC(); err != nil {
			//log记录rpc服务挂掉的原因
			glog.Glog.Errorln(err)
		}
	}()

	go s.DetectHostGroupAvailability()
	defer close(dstopCh)
	go s.DetectStopService(dstopCh)
	defer close(rstopCh)
	go s.RecoverProxySerivce(rstopCh)
	//go s.proxyEngine.ServeInBlocking(func() {
	//})
	go s.proxyEngine.ServeInBlocking(s.GetProxyCluster)

	if s.conf.DebugMode {
		go s.starDebugServer()
	}

	go s.ServeMetrics()

	s.status = StatusUp
	s.AddProxyMember()

	//启动心跳
	go func() {
		for {
			select {
			case <-hbstopCh:
				return
			case <-s.ticker.C:
				s.SendHearBeat()
			}
		}
	}()
	//for debug print shard config
	/*	go func() {
		for {
			s.shardConfig.Print()
			time.Sleep(10 * time.Second)
		}

	}()*/
	select {}
}

//Stop 方法负责停止Server
func (s *Server) Stop() {
	s.ticker.Stop()
	hbstopCh <- struct{}{}
}

func (s *Server) SetupSignalHandler() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sc
		glog.Glog.Errorf("Got signal [%s] to exit.", sig)
	}()
}

//start checkhealth server
func (s *ChServer) startChServer() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.ln = l
	glog.Glog.Infof("start check healthy server at %s ", s.addr)
	for {
		_, err := s.ln.Accept()
		if err != nil {
			return err
		}
	}
}

//start pprof debug server
func (s *Server) starDebugServer() {
	glog.Glog.Infof("start debug server on %v", s.conf.DebugServer)
	glog.Glog.Errorln(http.ListenAndServe(s.conf.DebugServer, nil))
}

//new checkHelthServer
func NewChser(addr string) (*ChServer, error) {
	if len(addr) == 0 {
		return nil, errors.New("check healthy server config is nil")
	}
	s := &ChServer{}
	s.addr = addr
	return s, nil
}

//close checkhealth server
func (s *ChServer) stopChServer() {
	if s.ln != nil {
		glog.Glog.Infoln("check healthy server stop")
		s.ln.Close()
	}
}

//处理grpc异常
func recoveryHandleFunc(p interface{}) error {
	if p != nil {
		glog.Glog.Errorf("rpc panic err: %v stack:%s", p, debug.Stack())
	}
	return nil
}
