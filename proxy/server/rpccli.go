package server

import (
	"context"
	"net"
	"strconv"

	"runtime/debug"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/system"
	"git.2dfire.net/zerodb/common/system/load"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"github.com/pkg/errors"
)

//rpc client interface

//PullConfig pull proxy config from keeper and send status to keeper
func (s *Server) PullConfig() error {
	in := new(keeper.PullConfigRequest)
	in.ClusterName = s.conf.ClusterName

	ctx, cancel := context.WithTimeout(context.TODO(), TimeOut)
	r, err := s.keeperCli.PullConfig(ctx, in)
	cancel()
	if err != nil {
		return errors.WithStack(err)
	}

	if err = s.shardConfig.LoadConfig(r.Data); err != nil {
		return errors.WithStack(err)
	}
	//第一次启动备份shardconfig
	s.shardConfig.SetVersion(r.Version)
	//s.shardConfig.oldConf.Version = r.Version

	if err = s.proxyEngine.ParserConfig(s.shardConfig.GetConfig()); err != nil {
		glog.Glog.Errorf("proxy parser shard config failed: %v", err)
		return errors.WithStack(err)
	}

	glog.Glog.Infof("Config has been fetched from keeper[%s].", s.conf.KeeperAddr)
	return nil
}

// Register 注册到keeper
/*
func (s *Server) Register() error {
	in := new(keeper.RegisterRequest)
	in.Address = s.conf.RPCServer

	loop, _, port := utils.IsLoopback(in.Address)
	if loop {
		ip := utils.GetLocalIP()
		in.Address = net.JoinHostPort(ip, port)
	}

	//TODO
	in.ClusterName = s.conf.ClusterName
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	resp, err := s.keeperCli.RegisterProxy(ctx, in)
	cancel()
	if err != nil {
		glog.ProxyLog.Errorf("proxy register failed: %v", err)
		return err
	} else if resp.Code != codes.OK {
		glog.ProxyLog.Errorf("proxy register failed: %v, code: %d", resp.ErrMsg, resp.Code)
		return errors.New(resp.ErrMsg)
	}
	return nil
}
*/

// add proxy member
func (s *Server) AddProxyMember() error {
	loop, host, port := utils.IsLoopback(s.conf.ProxyServer)
	if loop {
		host = utils.GetLocalIP()
	}
	nPort, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		glog.Glog.Errorln(err)
		return err
	}
	in := &keeper.ProxyMemberRequest{
		ClusterName: s.conf.ClusterName,
		Host:        host,
		Port:        nPort,
		Weight:      1,
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	resp, err := s.keeperCli.AddProxyMember(ctx, in)
	cancel()
	if err != nil {
		glog.Glog.Errorf("proxy add failed: %s", err.Error())
		return err
	}
	if resp.Code != codes.OK {
		glog.Glog.Errorf("proxy add member failed: %v, code: %d", resp.Code, resp.ErrMsg)
	}

	return nil
}

//SendHearBeat send hearbeat to keeper
func (s *Server) SendHearBeat() error {
	//recover sendheartbeat panic
	defer func() {
		if r := recover(); r != nil {
			glog.Glog.Errorln(debug.Stack())
		}
	}()
	in := new(keeper.ProxyHeartRequest)

	in.Address = s.conf.RPCServer

	loop, _, port := utils.IsLoopback(in.Address)
	if loop {
		ip := utils.GetLocalIP()
		in.Address = net.JoinHostPort(ip, port)
	}

	in.Status = s.status
	in.ClusterName = s.conf.ClusterName
	in.ConfVersion = s.shardConfig.GetVersion()

	in.SystemInfo = &keeper.SystemInfo{
		MemLoad: 0,
		CpuLoad: 0,
		LoadAvg: &keeper.LoadAvg{},
	}

	var err error

	//TODO
	in.SystemInfo.CpuLoad, err = system.CPULoad()
	if err != nil {
		glog.Glog.Warnf("get cpuload failed: %v", err)
	}

	in.SystemInfo.MemLoad, err = system.MemLoad()
	if err != nil {
		glog.Glog.Warnf("get memlogad failed: %v", err)
	}

	loadAvg, err := load.LoadAvg()
	if err != nil {
		glog.Glog.Warnf("get loadavg failed: %v", err)
	} else {
		in.SystemInfo.LoadAvg.Avg1Min = loadAvg.Avg1min
		in.SystemInfo.LoadAvg.Avg5Min = loadAvg.Avg5min
		in.SystemInfo.LoadAvg.Avg15Min = loadAvg.Avg15min
	}
	glog.Glog.WithFields("status", in.Status, "address", in.Address, "cluster", in.ClusterName).Infof("Sending heartbeat to keeper[%s]...", s.conf.KeeperAddr)

	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	_, err = s.keeperCli.ProxyHeartbeat(ctx, in)
	cancel()
	if err != nil {
		s.getLost = true
		glog.Glog.Errorf("proxy senthearbeat failed: %v", err)
		return err
	} else {
		s.getLost = false
	}

	return nil
}

//added by huajia,for cobar
func (s *Server) GetProxyCluster() {
	in := new(keeper.ClustersRequest)
	in.ClusterName = s.conf.ClusterName


	ctx, cancel := context.WithTimeout(context.TODO(), TimeOut)
	r, err := s.keeperCli.ProxyClusters(ctx, in)
	cancel()
	if err != nil {
		return
	}
	if r.Infos == nil {
		return
	}

	s.proxyEngine.SetProxyCluster(r.Infos)

	return
}
//end added
