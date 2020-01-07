package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/keeper/pkg/timer"
	"github.com/pkg/errors"

	"strings"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/models"
	"github.com/gin-gonic/gin"
)

const (
	proxyNotExist             = "no proxy "
	statusError               = "wrong status "
	parameterIsNil            = "nil parameter "
	clusterNameIsNil          = "nil [clusterName]"
	hostGroupClusterNameIsNil = "nil [hostGroupClusterName]"
	inconsistentError         = "inconsistent!!"
)

func rpcMessage(ip, msg string) string {
	return fmt.Sprintf("[ip]:%s,[message]:%s;", ip, msg)
}

func (s *Server) checkHealth(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (s *Server) proxyList(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)
	clusterName := c.Query("cluster_name")
	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [cluster_name]"
		return
	}
	addrs, err := s.store.GetClusters(clusterName)
	if err != nil {
		glog.GLog.Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}
	msg.Data = addrs
}

func (s *Server) login(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)
	userName := c.DefaultPostForm("userName", "")
	password := c.DefaultPostForm("password", "")
	if userName == "" || password == "" {
		msg.Code = codes.ParameterIsNil
		msg.ErrMsg = "userName or password is nil"
		return
	}

	if userName != "admin" || password != "admin" {
		msg.Code = codes.CommonError
		msg.ErrMsg = "userName or password error"
		return
	}

	j := NewJWT()
	token, err := j.GenerateToken(userName, password)
	if err != nil {
		glog.GLog.WithFields("name", userName, "password", password).Errorln(err)
		msg.Code = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}
	date := make(map[string]string)
	date["token"] = token

	msg.Data = date
}

//第一次用来初始化keeper的分库分表配置
func (s *Server) shardConfInit(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	clusterName := c.PostForm("clusterName")
	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}

	//是否强制
	force := c.PostForm("force")

	glog.GLog.WithFields("cluster", clusterName, "force", force).Infoln("query shard config init.")

	file, err := c.FormFile("file")
	if err != nil {
		glog.GLog.WithFields("cluster", clusterName, "force", force).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.DataReadFail
		msg.ErrMsg = err.Error()
		return
	}

	src, err := file.Open()
	if err != nil {
		glog.GLog.WithFields("cluster", clusterName, "force", force, "file", file.Filename).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		glog.GLog.Errorf("config init read file failed: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.DataReadFail
		msg.ErrMsg = err.Error()
		return
	}

	// force 是1为强制更新
	if force == "0" {
		etcdData, err := s.store.GetShardConf(clusterName, "")
		if err != nil {
			glog.GLog.WithField("cluster", clusterName).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			return
		}
		if len(etcdData) > 0 {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = "shard config exist"
			return
		}
	}

	conf := config.Config{}
	if err = utils.LoadYaml(data, &conf); err != nil {
		glog.GLog.WithField("cluster", clusterName).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		return
	}

	s.pushFailed = conf.Basic.PushWhenFail
	bChecked := checkConfig(&conf)

	if !bChecked {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ConfigCheckFailed
		msg.ErrMsg = "shard config check failed."
		return
	}

	if err = s.store.SetShardConf(clusterName, data); err != nil {
		glog.GLog.WithField("cluster", clusterName).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdPutFailed
		msg.ErrMsg = err.Error()
		return
	}
}

// 查看proxy集群的状态
func (s *Server) clusterStatus(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	clusterName := c.Query("clusterName")
	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}

	glog.GLog.WithField("cluster", clusterName).Infoln("cluster status.")

	key := fmt.Sprintf(etcdtool.StatusPrefix, clusterName)
	status, err := s.store.GetDatas(key)
	if err != nil {
		glog.GLog.WithField("cluster", clusterName).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	statusData := make(map[string]string)
	var index int
	if len(status) > 0 {
		for key, value := range status {
			index = strings.LastIndex(key, "/")
			if index > 0 {
				key = key[index+1:]
				infor := models.StatusInfo{}
				if err = utils.LoadJSON(value, &infor); err != nil {
					glog.GLog.WithField("cluster", clusterName).Errorln(err)
					continue
				}
				statusData[key] = infor.Status
			}
		}
	} else {
		msg.Data = make([]interface{}, 0)
		return
	}

	key = fmt.Sprintf(etcdtool.ProxyInforPrefix, clusterName)
	data, err := s.store.GetDatas(key)
	if err != nil {
		glog.GLog.WithField("cluster", clusterName).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	if len(data) > 0 {
		retData := make([]systemInfo, 0)
		for key, value := range data {
			info := systemInfo{}
			index = strings.LastIndex(key, "/")
			if index > 0 {
				key = key[index+1:]
				if err = utils.LoadJSON(value, &info); err != nil {
					glog.GLog.WithField("cluster", clusterName).Errorln(err)
					continue
				}
				for k, v := range statusData {
					if key == k {
						info.Status = v
					}
				}
				info.IP = key
				retData = append(retData, info)
			}
		}

		key = fmt.Sprintf(etcdtool.ProxyVersionPrefix, clusterName)
		data, err = s.store.GetDatas(key)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			glog.GLog.WithField("cluster", clusterName).Errorln(err)
			return
		}

		if len(data) > 0 {
			for key, v := range data {
				index = strings.LastIndex(key, "/")
				if index > 0 {
					key = key[index+1:]
					for i := range retData {
						if retData[i].IP == key {
							retData[i].ConfVersion = string(v)
						}
					}
				}
			}
		}
		msg.Data = retData
	}
}

//proxy集群数据源切换
func (s *Server) switchDB(c *gin.Context) {
	msg := NewResult()
	msg.Code = codes.Failed
	msg.ErrorCode = codes.CommonError
	defer func() {
		c.JSON(http.StatusOK, &msg)
	}()
	var switchDataSource SwitchDataSource
	if err := c.BindJSON(&switchDataSource); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	clusterName := switchDataSource.Cluster
	from := int32(*switchDataSource.From)
	to := int32(*switchDataSource.To)
	hostGroup := switchDataSource.HostGroup
	if from == to {
		msg.ErrMsg = fmt.Sprintf("from index %d equal to 'to' index %d ignore it", from, to)
		return
	}

	w, err := s.store.GetWriteIPs(clusterName, hostGroup)
	if err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	max := len(strings.Split(strings.TrimSpace(w), ","))
	if err := s.checkSwitchIndex(from, to, max); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	glog.GLog.Infof("web admin switchDB hostGroup: %s, from: %d, toIndex: %d, reason: %s ", hostGroup, from, to, switchDataSource.Reason)

	//检查集群的状态
	clis, ok := s.proxyClients.Load(clusterName)
	proxyCount := len(clis.(ProxyCli))
	if !ok || proxyCount == 0 {
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
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithField("cluster", clusterName).Warnln(statusError)
		return
	}
	//向集群中的所有proxy发起数据源切换,
	proxyRespMsg := make([]string, 0)
	in := new(proxy.SwitchDatasourceRequest)
	in.Hostgroup = hostGroup
	in.From = from
	in.To = to
	defer func() {
		msg.Data = proxyRespMsg
	}()
	msg.Code = codes.Succeed
	msg.ErrorCode = codes.OK

	//call proxy rpc
	var allProxySucess = true
	for ip, cli := range clis.(ProxyCli) {
		ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
		resp, err := cli.SwitchDataSource(ctx, in)
		cancel()
		if err != nil {
			proxyRespMsg = append(proxyRespMsg, rpcMessage(ip, err.Error()))
			glog.GLog.Errorf("[proxy:ip] failed switch data source form %d to %d  %v ", ip, from, to)
			allProxySucess = false
		} else {
			proxyRespMsg = append(proxyRespMsg, rpcMessage(ip, resp.BasicResp.ErrMsg))
		}

	}
	if !allProxySucess {
		return
	}
	//set active write etcd

	if err := s.store.SetHostGroupActiveWrite(clusterName, hostGroup, int(to)); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdPutFailed
		e := errors.Errorf("set hostGroup %s activeWrite  failed %v", hostGroup, err)
		glog.GLog.Errorln(e)
		msg.ErrMsg = inconsistentError + e.Error()
		glog.GLog.Errorf("failed set %s activeWrite %d to etcd", hostGroup, to)
		return
	}
	glog.GLog.Infof("sucess switch data source from %d to %d  proxy response %+v", from, to, proxyRespMsg)
	return

}

//注销proxy
func (s *Server) unregister(c *gin.Context) {
	msg := NewResult()
	msg.Code = codes.Failed
	msg.ErrorCode = codes.CommonError
	unregFailLog := "unreg cluster[%s] proxy[%s] failed %s"
	rpcResult := make([]string, 1)
	defer func() {
		msg.Data = rpcResult
		c.JSON(http.StatusOK, &msg)
	}()
	var unregProxy unregProxy
	if err := c.BindJSON(&unregProxy); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	cluster := unregProxy.Cluster
	address := unregProxy.Address
	reason := unregProxy.Reason
	glog.GLog.Infof("webadmin send a unreg cluster [%s] proxy [%s] reason [%s] request", cluster, address, reason)
	defer func() {
		if len(msg.ErrMsg) != 0 {
			glog.GLog.Errorf(unregFailLog, cluster, address, msg.ErrMsg)
		}

	}()
	var err error
	if _, _, err = net.SplitHostPort(address); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	//call proxy remove rpc to stop proxy
	clis, ok := s.proxyClients.Load(cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.ErrMsg = proxyNotExist + " [" + cluster + "]"
		return
	}
	cli := clis.(ProxyCli)[address]
	ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
	defer cancel()
	resp, err := cli.Remove(ctx, new(basic.EmptyRequest))
	if err != nil {
		msg.Code = codes.Succeed
		msg.ErrMsg = err.Error()
		rpcResult = append(rpcResult, rpcMessage(address, err.Error()))
		return
	}
	rpcResult = append(rpcResult, rpcMessage(address, resp.ErrMsg))
	//delete proxyClients cache and  proxy data in etcd
	if err = s.deleteProxy(cluster, address, true); err != nil {
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}
	glog.GLog.Infof("unreg cluster[%s] proxy[%s] reason [% s] success %+v", cluster, address, reason, rpcResult)
	msg.ErrorCode = codes.OK
	msg.Code = codes.Succeed
}

//添加分库分表配置快照
func (s *Server) snapshotConfig(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	clusterName := c.PostForm("clusterName")

	name := c.PostForm("snapshotName")

	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}

	if len(name) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [name]"
		return
	}

	now := time.Now()
	name = fmt.Sprintf("%s-%04d-%02d-%02d-%02d-%02d-%02d", name, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	glog.GLog.WithFields("cluster", clusterName, "snapshot", name).Infoln("query snapshot config.")

	if err := s.store.CreateSnapshot(clusterName, name); err != nil {
		glog.GLog.WithFields("cluster", clusterName, "snapshot", name).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdPutFailed
		msg.ErrMsg = err.Error()
		return
	}
}

//获取所有的配置快照
func (s *Server) snapshotList(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	clusterName := c.Query("clusterName")

	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	glog.GLog.WithField("cluster", clusterName).Infoln("query snapshot list.")

	retData, err := s.store.GetSnapshotList(clusterName)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", clusterName).Errorln(err)
		return
	}

	msg.Data = retData
}

//获取所有集群名称
func (s *Server) clusterList(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)
	list, err := s.store.GetClusterList()
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.Errorln(err)
		return
	}

	msg.Data = list
}

//updateStopService 更新
func (s *Server) updateStopService(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var parameter stopService
	var err error

	if err = c.BindJSON(&parameter); err != nil {
		glog.GLog.Errorf("update stopService error: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + parameter.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(parameter.Cluster)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithFields("cluster", parameter.Cluster).Warnln(statusError)
		return
	}

	var data []byte
	data, err = utils.UnLoadYaml(parameter.Service)
	if err != nil {
		glog.GLog.Errorf("add schema unloadyaml error: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	nextVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", parameter.Cluster).Errorln(err)
		return
	}

	in := new(proxy.UpdateStopServiceRequest)
	in.StopServiceData = data
	in.Version = nextVersion

	sMsgs := make([]string, 0)
	bSucceed := true

	for ip, cli := range clis.(ProxyCli) {
		ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
		resp, err := cli.UpdateStopService(ctx, in)
		cancel()
		if err != nil {
			sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
			bSucceed = false
			glog.GLog.WithFields("ip", ip, "cluster", parameter.Cluster).Errorln(err)
		} else {
			sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
			if resp.Code != codes.OK {
				bSucceed = false
				glog.GLog.WithFields("ip", ip, "cluster", parameter.Cluster, "code", resp.Code).Errorln(resp.ErrMsg)
			}
		}
	}

	if bSucceed {
		if err = s.store.UpdateStopSve(parameter.Cluster, parameter.Service); err != nil {
			glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = err.Error()
		}
	}

	msg.Data = sMsgs
}

//推送配置
func (s *Server) pushConfig(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	clusterName := c.PostForm("clusterName")
	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}

	snapshotName := c.PostForm("snapshotName")
	if len(snapshotName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [snapshotName]"
		return
	}

	todoTime := c.PostForm("doTime")

	if len(todoTime) > 0 {
		//定时任务，检测时间格式正确性
		tTime, err := utils.ParseTime(todoTime)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = err.Error()
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
			return
		}
		//如果当前时间在给定时间之后则是错误
		bAfter := time.Now().After(tTime)
		if bAfter {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = "time error."
			return
		}
		args := map[string]string{
			"cluster":  clusterName,
			"snapshot": snapshotName,
		}
		_, err = s.tTimer.Register(timer.Single, tTime.Sub(time.Now()), s.todoPush, args)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = "timer register failed"
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
			return
		}
	} else {
		//立即推送配置
		glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName).Infoln("push config.")
		clis, ok := s.proxyClients.Load(clusterName)
		if !ok || len(clis.(ProxyCli)) == 0 {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = proxyNotExist + " [" + clusterName + "]"
			return
		}

		status, err := s.store.GetClusterStatus(clusterName)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
			return
		}

		if status != models.StatusUp {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
			return
		}

		data, err := s.store.GetShardConf(clusterName, snapshotName)
		if err != nil {
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			return
		}

		sVersion, err := s.store.GetNextVersion(clusterName, true)
		if err != nil {
			glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			return
		}

		in := &proxy.PushConfigRequest{
			Data:    data,
			Version: sVersion,
		}

		sMsgs := make([]string, 0)
		bSucceed := true

		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			resp, err := cli.PushConfig(ctx, in)
			cancel()
			if err != nil {
				glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(resp.ErrMsg)
				}
			}
		}

		if bSucceed {
			if err = s.store.SetShardConf(clusterName, data); err != nil {
				glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Errorln(err)
				msg.Code = codes.Failed
				msg.ErrorCode = codes.EtcdPutFailed
				msg.ErrMsg = inconsistentError + err.Error()
				return
			}
		}

		if len(sMsgs) > 0 {
			msg.Data = sMsgs
		}
	}

	glog.GLog.WithFields("cluster", clusterName, "snapshot", snapshotName, "todoTime", todoTime).Infoln("push config successful.")
}

//回滚配置
func (s *Server) rollbackConfig(c *gin.Context) {
	msg := NewResult()
	msg.Code = codes.Failed
	msg.ErrorCode = codes.CommonError
	data := make([]string, 0)
	rollbackfaillog := "rollback cluster[%s] config failed %s"
	defer func() {
		c.JSON(http.StatusOK, &msg)
	}()
	var cluster cluster
	var err error
	if err = c.BindJSON(&cluster); err != nil {
		msg.ErrMsg = err.Error()
		return
	}
	clusterName := cluster.Cluster
	glog.GLog.Infof("webadmin send a rollback cluster [%s] config request", clusterName)
	defer func() {
		if len(msg.ErrMsg) != 0 {
			glog.GLog.Errorf(rollbackfaillog, clusterName, msg.ErrMsg)
		}

	}()
	var wg sync.WaitGroup
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

	// call all proxy rollbackConfig rpc
	for key, value := range clis.(ProxyCli) {
		wg.Add(1)
		go func(host string, client proxy.ProxyClient) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			resp, err := client.RollbackConfig(ctx, &basic.EmptyRequest{})
			cancel()
			if err != nil {
				glog.GLog.Errorf("rollback cluster[%s] config  failed, proxy[%s] return error %s", clusterName, host, err.Error())
				data = append(data, rpcMessage(host, err.Error()))
				return
			}
			data = append(data, rpcMessage(host, resp.ErrMsg))
		}(key, value)
	}
	//因为不用写etcd所以回滚是永远成功的.
	wg.Wait()
	//回滚etcd
	if err := s.store.RollBackLastModify(clusterName); err != nil {
		msg.Code = codes.Failed
		msg.ErrMsg = "rollback etcd failed"
		msg.Data = data
		msg.ErrorCode = codes.EtcdPutFailed
		return
	}
	glog.GLog.Infof("rollback etcd sucess")

	msg.Code = codes.Succeed
	msg.Data = data
	msg.ErrorCode = codes.OK
}

// updateBasic
func (s *Server) updateBasic(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)
	parameter := updateBasic{}
	if err := c.BindJSON(&parameter); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.Errorln(err)
		return
	}

	data, err := utils.UnLoadJSON(&parameter.Basic)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	sVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdPutFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	sMsgs := make([]string, 0)
	bSucceed := true

	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if ok {
		status, err := s.store.GetClusterStatus(parameter.Cluster)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			return
		}

		if status != models.StatusUp {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
			return
		}

		in := new(proxy.UpdateBasicRequest)
		in.BasicData = data
		in.Version = sVersion
		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.TODO(), WebTimeOut)
			resp, err := cli.UpdateBasic(ctx, in)
			cancel()
			if err != nil {
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
				glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithField("cluster", parameter.Cluster).Errorln(resp.ErrMsg)
				}
			}
		}
	}

	if bSucceed {
		key := fmt.Sprintf(etcdtool.ShardConfig, parameter.Cluster, etcdtool.DefaultConfig, etcdtool.Basic)

		ops := make(map[string]string)
		ops[key] = string(data)

		key = fmt.Sprintf(etcdtool.ConfigVersion, parameter.Cluster)

		ops[key] = sVersion

		if err = s.store.Puts(ops); err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = inconsistentError + err.Error()
			glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			return
		}
	}
	if len(sMsgs) == 0 {
		msg.Data = nil
	} else {
		msg.Data = sMsgs
	}
}

// updateSwitch
func (s *Server) updateSwitch(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)
	parameter := updateSwitch{}
	if err := c.BindJSON(&parameter); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.Errorln(err)
		return
	}
	data, err := utils.UnLoadJSON(&parameter.Switch)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	sVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdPutFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	sMsgs := make([]string, 0)
	bSucceed := true
	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if ok {
		status, err := s.store.GetClusterStatus(parameter.Cluster)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdGetFailed
			msg.ErrMsg = err.Error()
			glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			return
		}

		if status != models.StatusUp {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
			return
		}

		in := new(proxy.UpdateSwitchRequest)
		in.SwitchData = data
		in.Version = sVersion
		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.TODO(), WebTimeOut)
			resp, err := cli.UpdateSwitch(ctx, in)
			cancel()
			if err != nil {
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
				glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithField("cluster", parameter.Cluster).Errorln(resp.ErrMsg)
				}
			}
		}
	}
	if bSucceed {
		key := fmt.Sprintf(etcdtool.ShardConfig, parameter.Cluster, etcdtool.DefaultConfig, etcdtool.Switchdb)

		ops := make(map[string]string)
		ops[key] = string(data)

		key = fmt.Sprintf(etcdtool.ConfigVersion, parameter.Cluster)
		ops[key] = sVersion

		if err = s.store.Puts(ops); err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = inconsistentError + err.Error()
			glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			return
		}
	}
	if len(sMsgs) == 0 {
		msg.Data = nil
	} else {
		msg.Data = sMsgs
	}
}
