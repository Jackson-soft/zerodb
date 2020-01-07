package server

import (
	"context"
	"fmt"
	"net/http"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	hostgroupcluster = "hostgroupclusters"
)

var (
	ErrHgClusterNameExist    = errors.New("hostgroupcluster name allready exist")
	ErrHgClusterNameNotExist = errors.New("hostgroupcluster name not exist")
)

//增加hostgroup
func (s *Server) addHostGroups(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var hostGroups addHost
	if err := c.BindJSON(&hostGroups); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = err.Error()
		glog.GLog.Errorln(err)
		return
	}

	glog.GLog.WithField("addHost", hostGroups).Infoln("BEGIN: add hostGroups.")

	clis, ok := s.proxyClients.Load(hostGroups.Cluster)
	if !ok {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + hostGroups.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(hostGroups.Cluster)
	if err != nil {
		glog.GLog.WithFields("cluster", hostGroups.Cluster).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithFields("cluster", hostGroups.Cluster).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(hostGroups.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", hostGroups.Cluster).Errorln(err)
		return
	}

	in := new(proxy.AddHostGroupRequest)
	in.Version = nextVersion
	var data []byte
	sMsgs := make([]string, 0)
	bSucceed := true

	for _, v := range hostGroups.Groups {
		data, err = utils.UnLoadYaml(&v)
		if err != nil {
			glog.GLog.WithFields("cluster", hostGroups.Cluster, "api", "addHostGroup").Errorln(err)
			continue
		}
		in.HostGroupCfgData = data
		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)

			glog.GLog.WithField("AddHostGroupRequest", in).Infof("DISPATCH: cli:[%s] is calling AddHostGroup rpc", ip)

			resp, err := cli.AddHostGroup(ctx, in)
			cancel()

			if err != nil {
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
				glog.GLog.WithFields("ip", ip, "cluster", hostGroups.Cluster).Errorln(err)
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithFields("ip", ip, "cluster", hostGroups.Cluster, "errmsg", resp.ErrMsg).Errorln("error")
				}
			}
		}
	}

	if bSucceed {
		if err = s.store.AddHostGroup(hostGroups.Cluster, hostGroups.Groups); err != nil {
			glog.GLog.WithFields("cluster", hostGroups.Cluster).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = inconsistentError + err.Error()
			return
		}
	}

	msg.Data = sMsgs

	glog.GLog.WithField("Result", msg).Infof("END: add hostGroup successfully")
}

// 删除hostgroup
func (s *Server) delHostGroups(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var err error
	var delHosts deleteByNames
	if err = c.BindJSON(&delHosts); err != nil {
		glog.GLog.Errorf("delete hostGroups failed: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	clis, ok := s.proxyClients.Load(delHosts.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + delHosts.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(delHosts.Cluster)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", delHosts.Cluster).Errorln(err)
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithFields("cluster", delHosts.Cluster).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(delHosts.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", delHosts.Cluster).Errorln(err)
		return
	}

	sMsgs := make([]string, 0)
	bSucceed := true
	in := new(proxy.DeleteHostGroupRequest)
	in.Version = nextVersion

	for ip, cli := range clis.(ProxyCli) {
		for _, name := range delHosts.Names {
			in.HostGroupName = name
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			resp, err := cli.DeleteHostGroup(ctx, in)
			cancel()
			if err != nil {
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
				glog.GLog.WithFields("ip", ip, "cluster", delHosts.Cluster).Errorln(err)
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithFields("ip", ip, "cluster", delHosts.Cluster, "code", resp.Code).Errorln(resp.ErrMsg)
				}
			}
		}
	}

	if bSucceed {
		if err = s.store.DeleteGroup(delHosts.Cluster, delHosts.Names); err != nil {
			glog.GLog.WithFields("cluster", delHosts.Cluster).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = inconsistentError + err.Error()
			return
		}
	}

	msg.Data = sMsgs

	glog.GLog.Infof("delete hostGroup successfully")
}

//获取hostGroup web接口
func (s *Server) GetHostGroup(c *gin.Context) {
	msg := NewResult()
	msg.Data = ""
	defer func() {
		c.JSON(http.StatusOK, msg)
	}()
	var (
		snapshot                   = etcdtool.DefaultConfig
		hostgroupName, clusterName string
	)
	clusterName, _ = c.GetQuery("clustername")
	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = clusterNameIsNil
		return
	}
	snapshot, _ = c.GetQuery("snapshot")
	hostgroupName, _ = c.GetQuery("hostgroup")
	hostgroups, err := s.getHostGroups(clusterName, snapshot)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}
	if len(hostgroups) == 0 {
		return
	}
	if len(hostgroupName) == 0 {
		msg.Code = codes.OK
		msg.Data = hostgroups
		return
	}
	for _, hostgroup := range hostgroups {
		if hostgroup.Name == hostgroupName {
			msg.Data = hostgroup
			return
		}
	}
	return
}

func (s *Server) getHostGroups(clusterName, snapshot string) ([]config.HostGroup, error) {
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}
	if len(clusterName) == 0 {
		return nil, errors.New(clusterNameIsNil)
	}
	var hostgroups []config.HostGroup
	key := fmt.Sprintf(etcdtool.ShardConfig, clusterName, snapshot, "hostgroups")
	resp, err := s.store.GetData(key)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	if err = yaml.Unmarshal(resp, &hostgroups); err != nil {
		return nil, err
	}
	return hostgroups, nil
}

func (s *Server) addHostGroupCluster(clusterName string, hgcluster config.HostGroupCluster, do string, version string) error {
	if len(clusterName) == 0 {
		return errors.New(clusterNameIsNil)
	}
	if len(hgcluster.Name) == 0 {
		return errors.New(hostGroupClusterNameIsNil)
	}
	var hgclusters []config.HostGroupCluster
	var hgExist bool
	//get store hostgroupclusters
	hgclusterKey := fmt.Sprintf(etcdtool.ShardConfig, clusterName, etcdtool.DefaultConfig, hostgroupcluster)
	data, err := s.store.GetData(hgclusterKey)
	if err != nil {
		return err
	}
	if data != nil {
		if err := utils.LoadYaml(data, &hgclusters); err != nil {
			return err
		}
		//如果hostgroupcluster存在切do为create返回error
		for i, hgc := range hgclusters {
			if hgc.Name == hgcluster.Name {
				if do == "create" {
					return ErrHgClusterNameExist
				}
				hgExist = true
				hgclusters[i] = hgcluster
			}
		}
	}
	//如果hostgroupcluster 不存在添加数据
	if !hgExist {
		if do == "update" {
			return ErrHgClusterNameNotExist
		}
		hgclusters = append(hgclusters, hgcluster)
	}

	b, err := utils.UnLoadYaml(hgclusters)
	if err != nil {
		return err
	}
	if err := s.store.PutDataWithVersion(clusterName, hgclusterKey, fmt.Sprintf("%s", b), version); err != nil {
		return err
	}

	return nil
}

// AddhostgroupCluster web api
func (s *Server) AddHostGroupCluster(c *gin.Context) {
	msg := NewResult()
	msg.Code = codes.Failed
	msg.ErrorCode = codes.CommonError
	result := make([]string, 0)
	msg.Data = result
	addHostGroupClusterFailLog := "add hostGroupCluster %+v failed %s"
	defer func() {
		c.JSON(http.StatusOK, &msg)
	}()
	var hgclusterReq addhostGroupCluster
	if err := c.BindJSON(&hgclusterReq); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	data, err := yaml.Marshal(hgclusterReq.Hgcluster)
	if err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	clusterName := hgclusterReq.Cluster
	glog.GLog.Infof("webadmin send a add hostgroupcluster [%+v]  request", hgclusterReq)
	defer func() {
		if len(msg.ErrMsg) != 0 {
			glog.GLog.Errorf(addHostGroupClusterFailLog, hgclusterReq, msg.ErrMsg)
		}

	}()
	//request proxy 验证hostgroupcluster合法性
	clis, ok := s.proxyClients.Load(clusterName)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.ErrMsg = proxyNotExist + " [" + clusterName + "]"
		return
	}
	status, err := s.store.GetClusterStatus(clusterName)
	if err != nil {
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	if status != models.StatusUp {
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		return
	}
	version, err := s.store.GetNextVersion(clusterName, false)
	if err != nil {
		msg.ErrMsg = fmt.Sprintf("get nextversion failed %v", err)
		return
	}

	in := new(proxy.AddHostGroupClusterRequest)
	in.HostGroupClusterCfgData = data
	in.Version = version

	//call All Proxy RPC AddHostGroupCluster
	var allProxySucess = true
	msg.Code = codes.Succeed
	for ip, cli := range clis.(ProxyCli) {
		ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
		resp, err := cli.AddHostGroupCluster(ctx, in)
		cancel()
		if err != nil {
			result = append(result, rpcMessage(ip, err.Error()))
			msg.Data = result
			msg.ErrMsg = errors.Wrapf(err, " add [cluster:%s] [proxy:%s] hostgroupcluster error", clusterName, ip).Error()
			allProxySucess = false
			glog.GLog.Errorf("add hostGroupCluster [requset:%+v] to [proxy:%s] failed %v  ", hgclusterReq, ip, err)
		} else {
			result = append(result, rpcMessage(ip, resp.ErrMsg))
		}

	}
	//proxy rpc 调用失败直接返回
	if !allProxySucess {
		msg.Code = codes.Succeed
		msg.ErrorCode = codes.OK
		msg.Data = result
		return
	}
	//add hostgroupcluster to etcd
	if err = s.addHostGroupCluster(clusterName, hgclusterReq.Hgcluster, "create", version); err != nil {
		msg.Code = codes.Failed
		msg.Data = result
		msg.ErrMsg = inconsistentError + err.Error()
		msg.ErrorCode = codes.EtcdPutFailed
		glog.GLog.Errorf("add [cluster:%s] hostGroupCluster failed %+v", err)
		return
	} else {
		msg.Code = codes.Succeed
		msg.ErrorCode = codes.OK
		msg.Data = result
		glog.GLog.Infof("add [cluster:%s] hostGroupCluster successfully %+v", clusterName, hgclusterReq)
		return
	}
}

func (s *Server) UpdateHostGroupCluster(c *gin.Context) {
	msg := NewResult()
	msg.Code = codes.Failed
	msg.ErrorCode = codes.CommonError
	result := make([]string, 0)
	msg.Data = result
	updateHostGroupClusterFailLog := "update hostGroupCluster %+v failed %s"
	defer func() {
		c.JSON(http.StatusOK, &msg)
	}()
	var hgclusterReq addhostGroupCluster
	glog.GLog.Infof("webadmin send a update hostgroupcluster [%+v]  request", hgclusterReq)
	defer func() {
		if len(msg.ErrMsg) != 0 {
			glog.GLog.Errorf(updateHostGroupClusterFailLog, hgclusterReq, msg.ErrMsg)
		}

	}()
	if err := c.BindJSON(&hgclusterReq); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	data, err := yaml.Marshal(hgclusterReq.Hgcluster)
	if err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	clusterName := hgclusterReq.Cluster

	//request proxy 验证hostgroupcluster合法性
	clis, ok := s.proxyClients.Load(clusterName)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.ErrMsg = proxyNotExist + " [" + clusterName + "]"
		return
	}

	status, err := s.store.GetClusterStatus(clusterName)
	if err != nil {
		glog.GLog.WithFields("cluster", clusterName).Errorln(err)
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithFields("cluster", clusterName).Warnln(statusError)
		return
	}

	version, err := s.store.GetNextVersion(clusterName, false)
	if err != nil {
		msg.ErrMsg = fmt.Sprintf("get nextversion failed %v", err)
		return
	}
	var allProxySucess = true
	for ip, cli := range clis.(ProxyCli) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		in := new(proxy.UpdateHostGroupClusterRequest)
		in.HostGroupClusterCfgData = data
		in.Version = version
		resp, err := cli.UpdateHostGroupCluster(ctx, in)
		cancel()
		if err != nil {
			msg.Code = codes.Succeed
			result = append(result, rpcMessage(ip, err.Error()))
			msg.Data = result
			msg.ErrMsg = errors.Wrapf(err, " add [cluster:%s] [proxy:%s] hostgroupcluster error", clusterName, ip).Error()
			allProxySucess = false
			glog.GLog.Errorf(" add [cluster:%s] [proxy:%s] hostgroupcluster %v", clusterName, ip, err)
		} else {
			result = append(result, rpcMessage(ip, resp.ErrMsg))
		}

	}
	if !allProxySucess {
		return
	}

	//add hostgroupcluster to etcd
	if err := s.addHostGroupCluster(clusterName, hgclusterReq.Hgcluster, "update", version); err != nil {
		msg.Code = codes.Failed
		msg.Data = result
		msg.ErrMsg = inconsistentError + err.Error()
		msg.ErrorCode = codes.EtcdPutFailed
		return
	}
	msg.Code = codes.Succeed
	msg.ErrorCode = codes.OK
	msg.Data = result
	glog.GLog.Infof("update hostGroupCluster successfully %+v", hgclusterReq)
}

// delete hostgroupCluster web api
func (s *Server) DelHostGroupCluster(c *gin.Context) {
	msg := NewResult()
	msg.Code = codes.Failed
	msg.ErrorCode = codes.CommonError
	result := make([]string, 0)
	msg.Data = result
	delHostGroupClusterFailLog := "delete hostGroupCluster %+v failed %s"
	defer func() {
		c.JSON(http.StatusOK, &msg)
	}()
	var delHgclusterReq delHostGroupCluster
	if err := c.BindJSON(&delHgclusterReq); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	glog.GLog.Infof("webadmin send a delete hostgroupcluster [%+v]  request", delHgclusterReq)
	defer func() {
		if len(msg.ErrMsg) != 0 {
			glog.GLog.Errorf(delHostGroupClusterFailLog, delHgclusterReq, msg.ErrMsg)
		}

	}()

	//send rpc invoke to proxy
	clusterName := delHgclusterReq.Cluster
	hgclusterNames := delHgclusterReq.HostGroupClusterNames
	clis, ok := s.proxyClients.Load(clusterName)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.ErrMsg = proxyNotExist + " [" + clusterName + "]"
		return
	}

	status, err := s.store.GetClusterStatus(clusterName)
	if err != nil {
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	if status != models.StatusUp {
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		return
	}
	version, err := s.store.GetNextVersion(clusterName, false)
	if err != nil {
		msg.ErrMsg = fmt.Sprintf("get nextversion failed %v", err)
		return
	}
	var allProxySucess = true
	for ip, cli := range clis.(ProxyCli) {
		for _, hgclusterName := range hgclusterNames {
			ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
			in := new(proxy.DeleteHostGroupClusterRequest)
			in.HostGroupClusterName = hgclusterName
			in.Version = version
			resp, err := cli.DeleteHostGroupCluster(ctx, in)
			cancel()
			if err != nil {
				result = append(result, rpcMessage(ip, err.Error()))
				msg.Data = result
				msg.ErrMsg = errors.Wrapf(err, " del [cluster:%s] [proxy:%s] [hostgroupcluster:%s]  error", clusterName, ip, hgclusterName).Error()
				glog.GLog.Errorln(err)
				msg.Code = codes.Succeed
				allProxySucess = false
				glog.GLog.Errorf(" del [cluster:%s] [proxy:%s] [hostgroupcluster:%s]  %v", clusterName, ip, hgclusterName, err)
			} else {
				result = append(result, rpcMessage(ip, resp.ErrMsg))
			}

		}

	}
	if !allProxySucess {
		return
	}

	err = s.delhostGroupCluster(delHgclusterReq, version)
	if err != nil {
		msg.Code = codes.Failed
		msg.Data = result
		msg.ErrorCode = codes.EtcdPutFailed
		msg.ErrMsg = inconsistentError + err.Error()
		return
	}
	msg.Code = codes.Succeed
	msg.ErrorCode = codes.OK
	msg.Data = result

	glog.GLog.Infof("delete hostGroupCluster successfully %v", delHgclusterReq)
}

func (s *Server) delhostGroupCluster(req delHostGroupCluster, version string) error {
	//clusterName, snapshot string, hgcluster string
	clusterName := req.Cluster
	//snapshot := req.SnapsHot
	snapshot := etcdtool.DefaultConfig
	hgclusterNames := req.HostGroupClusterNames
	if len(clusterName) == 0 {
		return errors.New(clusterNameIsNil)
	}
	/*	if len(snapshot) == 0 {
		snapshot = "default"
	}*/
	//hgClusterNotExist := errors.New("hostgroupcluster doesn't exist")
	var hgclusters []config.HostGroupCluster
	hgclusterKey := fmt.Sprintf(etcdtool.ShardConfig, clusterName, snapshot, hostgroupcluster)
	s.store.GetData(hgclusterKey)
	data, err := s.store.GetData(hgclusterKey)
	if err != nil {
		return err
	}
	if data == nil {
		return errors.Errorf("check input could't get data  use  %s", hgclusterKey)
	}
	err = utils.LoadYaml(data, &hgclusters)
	if err != nil {
		return err
	}
	var notExistHgclusters []string
	for _, hgclusterName := range hgclusterNames {
		name := hgclusterName
		for i, hgc := range hgclusters {
			if name == hgc.Name {
				//index = i + 1
				hgclusters = append(hgclusters[:i], hgclusters[i+1:]...)
				break
			}
			if i == len(hgclusters)-1 {
				notExistHgclusters = append(notExistHgclusters, name)
			}

		}

	}
	b, err := utils.UnLoadYaml(hgclusters)
	if err != nil {
		return err
	}
	//put data with version
	if err := s.store.PutDataWithVersion(clusterName, hgclusterKey, fmt.Sprintf("%s", b), version); err != nil {
		return err
	}
	/*	if err := s.store.PutData(hgclusterKey, fmt.Sprintf("%s", b)); err != nil {
		return err
	}*/
	/*	if len(notExistHgclusters) != 0 {
		return errors.WithMessage(hgClusterNotExist, fmt.Sprintf("%s", notExistHgclusters))
	}*/
	return nil
}
