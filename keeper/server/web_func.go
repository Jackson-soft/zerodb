package server

import (
	"context"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/models"
)

// web common function

// 定时推送配置任务
func (s *Server) todoPush(args interface{}) {
	v, ok := args.(map[string]string)
	if !ok {
		return
	}
	clusterName := v["cluster"]
	snapshotName := v["snapshot"]
	clis, ok := s.proxyClients.Load(clusterName)
	if !ok {
		return
	}

	status, err := s.store.GetClusterStatus(clusterName)
	if err != nil {
		glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName).Errorln(err)
		return
	}

	if status != models.StatusUp {
		glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName).Errorln("proxy status is not up.")
		return
	}

	data, err := s.store.GetShardConf(clusterName, snapshotName)
	if err != nil {
		glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName).Errorln(err)
		return
	}

	sVersion, err := s.store.GetNextVersion(clusterName, true)
	if err != nil {
		glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName).Errorln(err)
		return
	}

	in := &proxy.PushConfigRequest{
		Data:    data,
		Version: sVersion,
	}

	bSucceed := true

	for ip, cli := range clis.(ProxyCli) {
		ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
		resp, err := cli.PushConfig(ctx, in)
		cancel()
		if err != nil {
			bSucceed = false
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "ip", ip).Errorln(err)
		} else {
			if resp.Code != codes.OK {
				bSucceed = false
				glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "ip", ip).Errorln(resp.ErrMsg)
			}
		}
	}
	if bSucceed {
		if err = s.store.SetShardConf(clusterName, data); err != nil {
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName).Errorln(err)
		}
	}
}

func checkConfig(conf *config.Config) bool {
	// basic
	if len(conf.Basic.ConfigName) == 0 || len(conf.Basic.User) == 0 || len(conf.Basic.Password) == 0 || conf.Basic.SlowLogTime <= 0 {
		return false
	}

	//stop_service
	if conf.StopService.OfflineSwhRejectedNum <= 0 || conf.StopService.OfflineDownHostNum < 0 {
		return false
	}

	// switch
	if conf.Switch.SafeLoad <= 0 || conf.Switch.VoteApproveRatio <= 0 || conf.Switch.VoteApproveRatio > 100 || conf.Switch.SafeBinlogDelay <= 0 || conf.Switch.Binlogwaittime <= 0 || conf.Switch.Frequency <= 0 || conf.Switch.BackendPingInterval < 0 {
		return false
	}

	return true
}
