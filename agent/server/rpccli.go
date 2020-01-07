package server

import (
	"context"
	"net"

	"git.2dfire.net/zerodb/agent/glog"
	"git.2dfire.net/zerodb/common/system"
	"git.2dfire.net/zerodb/common/system/load"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
)

//rpc client interface

//Register
/*
func (s *Server) Register() error {
	in := new(keeper.RegisterRequest)
	in.Address = s.conf.RPCServer

	loop, _, port := utils.IsLoopback(in.Address)
	if loop {
		ip := utils.GetLocalIP()
		in.Address = net.JoinHostPort(ip, port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	resp, err := s.keeperCli.RegisterAgent(ctx, in)
	cancel()
	if err != nil {
		glog.AgentLog.Errorf("agent register failed: %v", err)
		return errors.WithStack(err)
	} else if resp.Code != codes.OK {
		glog.AgentLog.Errorf("agent register failed with code: %d, errmsg: %s", resp.Code, resp.ErrMsg)
		return errors.New(fmt.Sprintf("agent register failed with code: %d, errmsg: %s", resp.Code, resp.ErrMsg))
	}
	return nil
}
*/

//heartbeat
func (s *Server) heartbeat() {
	in := new(keeper.AgentHeartRequest)

	in.Mysql = new(keeper.Mysql)

	mysql, err := s.mysqlInfo()
	if err != nil {
		glog.Glog.Errorf("agent heartbeat get mysql information failed: %v", err)
		in.Mysql.Status = mySQLDie
	} else {
		in.Mysql = mysql
	}

	in.Address = s.conf.RPCServer
	// TODO
	loop, _, port := utils.IsLoopback(in.Address)
	if loop {
		ip := utils.GetLocalIP()
		in.Address = net.JoinHostPort(ip, port)
	}

	in.SystemInfo = new(keeper.SystemInfo)
	in.SystemInfo.CpuLoad, err = system.CPULoad()
	if err != nil {
		glog.Glog.Errorf("agent AgentHeartbeat cpu failed: %v", err)
	}

	in.SystemInfo.MemLoad, err = system.MemLoad()
	if err != nil {
		glog.Glog.Errorf("agent AgentHeartbeat memory failed: %v", err)
	}

	in.SystemInfo.LoadAvg = new(keeper.LoadAvg)
	loadAvg, err := load.LoadAvg()
	if err != nil {
		glog.Glog.Errorf("agent AgentHeartbeat avg failed: %v", err)
	} else {
		in.SystemInfo.LoadAvg.Avg1Min = loadAvg.Avg1min
		in.SystemInfo.LoadAvg.Avg5Min = loadAvg.Avg5min
		in.SystemInfo.LoadAvg.Avg15Min = loadAvg.Avg15min
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	resp, err := s.keeperCli.AgentHeartbeat(ctx, in)
	cancel()
	if err != nil {
		glog.Glog.Errorf("agent AgentHeartbeat failed: %v", err)
	} else if resp.Code != codes.OK {
		glog.Glog.Errorf("agent AgentHeartbeat failed with code: %d, errmsg: %s", resp.Code, resp.ErrMsg)
	}

	glog.Glog.Infof("agent heartbeat address: %s", in.Address)
}
