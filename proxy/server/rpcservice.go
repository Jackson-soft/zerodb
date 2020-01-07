//rpc.go 主要包含注册到ProxyServer的rpc方法.

package server

import (
	"context"
	"os"
	"time"

	"git.2dfire.net/zerodb/common/config"
	"gopkg.in/yaml.v3"

	"fmt"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
)

func (s *Server) Restart(ctx context.Context, in *basic.EmptyRequest) (*basic.BasicResponse, error) {
	s.Lock()
	defer s.Unlock()

	glog.Glog.Infof("Prepare for restarting...")

	out := new(basic.BasicResponse)

	s.status = StatusReady

	var err error

	glog.Glog.Infof("Resending [ready] heartbeat...")

	//发送ready心跳
	if err = s.SendHearBeat(); err != nil {
		glog.Glog.Errorln("proxy ready failed: %v", err)
	}

	glog.Glog.Infof("Repulling config from keeper...")
	if err = s.PullConfig(); err != nil {
		glog.Glog.Errorf("proxy pull/parse config failed: %v", err)
	}

	if err = s.proxyEngine.ApplyLatestConfig(*s.shardConfig.GetConfig()); err != nil {
		glog.Glog.Errorf("proxy ApplyLatestConfig failed: %v", err)
		return nil, err
	}
	glog.Glog.Infof("Proxy Restart successfully")

	//注册
	/*
		if err = s.Register(); err != nil {
			glog.ProxyLog.Errorln("proxy register failed: %v", err)
		}
	*/
	if err = s.AddProxyMember(); err != nil {
		glog.Glog.Errorln(err)
	}

	s.status = StatusUp

	out.Code = codes.OK
	return out, nil
}

func (s *Server) GetVote(ctx context.Context, in *proxy.GetVoteRequest) (*proxy.GetVoteResponse, error) {
	glog.Glog.Infof("Receive a hostgroup[%s] vote request from keeper", in.Hostgroup)
	if len(in.Hostgroup) == 0 {
		return nil, errcode.BuildError(errcode.ObjectEmptyErr, "in.Hostgroup")
	}
	hostGroup := s.proxyEngine.GetHostGroupNode(in.Hostgroup)
	if hostGroup == nil {
		return nil, errcode.BuildError(errcode.HostGroupNotExist, in.Hostgroup)
	}
	writeDBPool := hostGroup.Write[in.From]
	if writeDBPool == nil {
		return nil, fmt.Errorf("writeDBPool for index:[%v] doesn't exist", in.From)
	}

	err := writeDBPool.Ping()
	if err != nil {
		glog.Glog.Infof("hostGroup[%s] agrees with switching ", in.Hostgroup)
		return &proxy.GetVoteResponse{
			BasicResp: new(basic.BasicResponse),
			From:      in.From,
			Vote:      true,
		}, nil
	}
	glog.Glog.Infof("hostGroup[%s] disagrees with switching ", in.Hostgroup)
	return &proxy.GetVoteResponse{
		BasicResp: new(basic.BasicResponse),
		From:      in.From,
		Vote:      false,
	}, nil

}

func (s *Server) Version(ctx context.Context, in *basic.EmptyRequest) (*proxy.VersionResponse, error) {
	out := &proxy.VersionResponse{}
	out.BackVersion = s.shardConfig.oldConf.Version
	out.CurVersion = s.shardConfig.currentConf.Version
	return out, nil
}

//RollbackConfig 回滚配置
func (s *Server) RollbackConfig(ctx context.Context, in *basic.EmptyRequest) (*basic.BasicResponse, error) {
	out := new(basic.BasicResponse)
	if err := s.shardConfig.Rollback(); err != nil {
		glog.Glog.Errorf("proxy rollback config failed: %v", err)
		return nil, err
	}
	if err := s.proxyEngine.ApplyBackConfig(s.shardConfig.currentConf); err != nil {
		glog.Glog.Errorf("proxy rollback config failed: %v", err)
		return nil, err
	}
	out.Code = codes.OK
	out.ErrMsg = fmt.Sprintf("rollback success current shardconfig version %s", s.shardConfig.currentConf.Version)
	return out, nil
}

func (s *Server) SwitchDataSource(ctx context.Context, in *proxy.SwitchDatasourceRequest) (*proxy.SwitchDatasourceResponse, error) {
	out := new(basic.BasicResponse)

	if len(in.Hostgroup) == 0 {
		return nil, errcode.BuildError(errcode.ObjectEmptyErr, "in.Hostgroup")
	}
	hostGroup := s.proxyEngine.GetHostGroupNode(in.Hostgroup)
	if hostGroup == nil {
		return nil, errcode.BuildError(errcode.HostGroupNotExist, hostGroup)
	}

	var msg string
	var err error
	msg, err = hostGroup.SwitchDatasource(int(in.To))
	if err != nil {
		glog.Glog.Errorf("SwitchDataSource Failed %v", err)
		return nil, err
	}
	// 重置ping计数器
	s.hostGroupPingFailMap.Store(in.Hostgroup, 0)

	out.Code = codes.OK
	out.ErrMsg = msg
	return &proxy.SwitchDatasourceResponse{
		BasicResp: out,
	}, nil
}

//停止hostgroup的数据库写操作
func (s *Server) StopWritingAbility(ctx context.Context, in *proxy.StopWriteRequest) (*basic.BasicResponse, error) {
	if err := s.proxyEngine.StopWritingAbility(in.Hostgroup); err != nil {
		return nil, err
	}

	return new(basic.BasicResponse), nil
}

func (s *Server) RecoverWritingAbility(ctx context.Context, in *proxy.RecoverWriteRequest) (*basic.BasicResponse, error) {
	if err := s.proxyEngine.RecoverWritingAbility(in.Hostgroup); err != nil {
		return nil, err
	}

	return new(basic.BasicResponse), nil
}

func (s *Server) AddHostGroup(ctx context.Context, in *proxy.AddHostGroupRequest) (*basic.BasicResponse, error) {
	hostGroupCfg := config.HostGroup{}
	if err := yaml.Unmarshal(in.HostGroupCfgData, &hostGroupCfg); err != nil {
		return nil, err
	}
	var msg string
	var err error
	if msg, err = s.proxyEngine.AddHostGroup(hostGroupCfg); err != nil {
		return nil, err
	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.AddHostGroup(hostGroupCfg)
	s.shardConfig.SetVersion(in.Version)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) AddHostGroupCluster(ctx context.Context, in *proxy.AddHostGroupClusterRequest) (*basic.BasicResponse, error) {
	hostGroupClusterCfg := config.HostGroupCluster{}
	if err := yaml.Unmarshal(in.HostGroupClusterCfgData, &hostGroupClusterCfg); err != nil {
		return nil, err
	}
	var msg string
	var err error
	if msg, err = s.proxyEngine.AddHostGroupCluster(hostGroupClusterCfg); err != nil {
		return nil, err
	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}

	s.shardConfig.AddHostGroupCluster(hostGroupClusterCfg)
	s.shardConfig.SetVersion(in.Version)
	glog.Glog.Infof("success add hostgroupcluster %+v", in)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) UpdateHostGroupCluster(ctx context.Context, in *proxy.UpdateHostGroupClusterRequest) (*basic.BasicResponse, error) {
	hostGroupClusterCfg := config.HostGroupCluster{}
	if err := yaml.Unmarshal(in.HostGroupClusterCfgData, &hostGroupClusterCfg); err != nil {
		return nil, err
	}
	var msg string
	var err error
	if msg, err = s.proxyEngine.UpdateHostGroupCluster(hostGroupClusterCfg); err != nil {
		return nil, err
	}

	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}

	s.shardConfig.DelHostGroupCluster(hostGroupClusterCfg.Name)
	s.shardConfig.AddHostGroupCluster(hostGroupClusterCfg)
	s.shardConfig.SetVersion(in.Version)
	glog.Glog.Infof("success update hostgroupcluster %+v", in)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

//PushConfig 提供向proxy推送配置的rpc调用
func (s *Server) PushConfig(ctx context.Context, in *proxy.PushConfigRequest) (*basic.BasicResponse, error) {
	out := new(basic.BasicResponse)

	var err error
	if err = s.shardConfig.LoadConfig(in.Data); err != nil {
		glog.Glog.Errorf("proxy pushconfig failed: %v", err)
		return nil, err
	}
	//set version
	s.shardConfig.currentConf.Version = in.Version
	if err = s.proxyEngine.ApplyLatestConfig(*s.shardConfig.GetConfig()); err != nil {
		glog.Glog.Errorf("proxy ApplyLatestConfig failed: %v", err)
		return nil, err
	}
	out.Code = codes.OK
	out.ErrMsg = "PushConfig success, " + "current version: " + in.Version
	return out, nil
}

func (s *Server) DeleteHostGroupCluster(ctx context.Context, in *proxy.DeleteHostGroupClusterRequest) (*basic.BasicResponse, error) {
	var msg string
	var err error
	if msg, err = s.proxyEngine.DeleteHostGroupCluster(in.HostGroupClusterName); err != nil {
		return nil, err
	}

	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.DelHostGroupCluster(in.HostGroupClusterName)
	s.shardConfig.SetVersion(in.Version)
	glog.Glog.Infof("success delete hostgroupcluster %+v", in)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) DeleteHostGroup(ctx context.Context, in *proxy.DeleteHostGroupRequest) (*basic.BasicResponse, error) {
	var msg string
	var err error

	if msg, err = s.proxyEngine.DeleteHostGroup(in.HostGroupName); err != nil {
		return nil, err
	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}

	s.shardConfig.DelHostGroup(in.HostGroupName)
	s.shardConfig.SetVersion(in.Version)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) AddSchema(ctx context.Context, in *proxy.AddSchemaRequest) (*basic.BasicResponse, error) {
	schemaCfg := config.SchemaConfig{}
	if err := yaml.Unmarshal(in.SchemaCfgData, &schemaCfg); err != nil {
		return nil, err
	}

	var msg string
	var err error
	if msg, err = s.proxyEngine.AddSchema(schemaCfg); err != nil {
		glog.Glog.Errorf("add schema[%s] failed: %v", schemaCfg.Name, err)
		return nil, err
	}
	//back current config
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.AddSchema(schemaCfg)
	s.shardConfig.SetVersion(in.Version)

	glog.Glog.Infof("success add schema %+v", in)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) DeleteSchema(ctx context.Context, in *proxy.DeleteSchemaRequest) (*basic.BasicResponse, error) {
	var msg string
	var err error

	if msg, err = s.proxyEngine.DeleteSchema(in.SchemaName); err != nil {
		glog.Glog.Errorf("delete schema[%s] failed: %v", in.SchemaName, err)
		return nil, err
	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}

	s.shardConfig.DelSchema(in.SchemaName)
	s.shardConfig.SetVersion(in.Version)
	glog.Glog.Infof("success del schema %+v", in)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) AddTable(ctx context.Context, in *proxy.AddTableRequest) (*basic.BasicResponse, error) {
	tableCfg := config.TableConfig{}
	if err := yaml.Unmarshal(in.TableCfgData, &tableCfg); err != nil {
		return nil, err
	}

	var msg string
	var err error
	if msg, err = s.proxyEngine.AddTable(tableCfg, in.SchemaName); err != nil {
		glog.Glog.Errorf("add table(s) of schema[%s] failed: %v", in.SchemaName, err)
		return nil, err
	}

	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.AddTable(in.SchemaName, tableCfg)
	s.shardConfig.SetVersion(in.Version)

	glog.Glog.Infof("success add table %+v", in)
	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) DeleteTable(ctx context.Context, in *proxy.DeleteTableRequest) (*basic.BasicResponse, error) {
	var msg string
	var err error
	if msg, err = s.proxyEngine.DeleteTable(in.TableName, in.SchemaName); err != nil {
		glog.Glog.Errorf("delete table[%s] of schema[%s] failed: %v", in.TableName, in.SchemaName, err)
		return nil, err
	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.DelTable(in.SchemaName, in.TableName)
	s.shardConfig.SetVersion(in.Version)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) UpdateSchemaRWSplit(ctx context.Context, in *proxy.UpdateSchemaRWSplitRequest) (*basic.BasicResponse, error) {
	var msg string
	var err error

	if msg, err = s.proxyEngine.UpdateSchemaRWSplit(in.SchemaName, in.RwSplit); err != nil {
		return nil, err
	}

	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.UpdateRwSplit(in.SchemaName, in.RwSplit)
	s.shardConfig.SetVersion(in.Version)

	glog.Glog.Infof("success add table %+v", in)

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = msg + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) UpdateStopService(ctx context.Context, in *proxy.UpdateStopServiceRequest) (*basic.BasicResponse, error) {
	stopServiceCfg := config.StopService{}
	if err := yaml.Unmarshal(in.StopServiceData, &stopServiceCfg); err != nil {
		return nil, err
	}
	// TODO need more specific validations
	if stopServiceCfg.OfflineDownHostNum <= 0 {

	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.currentConf.StopService = stopServiceCfg
	s.shardConfig.SetVersion(in.Version)

	glog.Glog.Infof("success update stopService %+v", in)
	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = "UpdateStopService success" + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) UpdateSwitch(ctx context.Context, in *proxy.UpdateSwitchRequest) (*basic.BasicResponse, error) {
	switchCfg := config.SwitchDB{}
	if err := yaml.Unmarshal(in.SwitchData, &switchCfg); err != nil {
		return nil, err
	}
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.currentConf.Switch = switchCfg
	s.shardConfig.SetVersion(in.Version)
	glog.Glog.Infof("success update switch %+v", in)
	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = "UpdateSwitch success" + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

func (s *Server) UpdateBasic(ctx context.Context, in *proxy.UpdateBasicRequest) (*basic.BasicResponse, error) {
	basicCfg := config.Basic{}
	if err := yaml.Unmarshal(in.BasicData, &basicCfg); err != nil {
		return nil, err
	}

	s.proxyEngine.ReconfigBasic(&basicCfg)
	// TODO need more specific validations
	if err := s.shardConfig.BackupCurConfig(); err != nil {
		return nil, err
	}
	s.shardConfig.currentConf.Basic = basicCfg
	s.shardConfig.SetVersion(in.Version)
	glog.Glog.Infof("success update basic %+v", in)
	out := new(basic.BasicResponse)
	out.Code = codes.OK
	out.ErrMsg = "UpdateBasic success" + ", " + s.shardConfig.oldConf.Version + "->" + s.shardConfig.currentConf.Version
	return out, nil
}

//unregi RPC方法负责注销proxy
func (s *Server) Remove(ctx context.Context, in *basic.EmptyRequest) (*basic.BasicResponse, error) {
	hbstopCh <- struct{}{}
	glog.Glog.Infoln("stop heartbeat service ")
	//stop Recover Services
	rstopCh <- struct{}{}
	glog.Glog.Infoln("stop recover service")
	dstopCh <- struct{}{}
	glog.Glog.Infoln("stop detect  stop service")
	//stop  Mysql Engine
	s.proxyEngine.CloseProxyEngine()
	glog.Glog.Infoln("stop Mysql Engine")
	go s.DelayExit(10 * time.Second)
	out := new(basic.BasicResponse)
	out.ErrMsg = "Success OK"
	return out, nil

}

func (s *Server) DelayExit(duration time.Duration) {
	time.Sleep(duration)
	os.Exit(0)
}
