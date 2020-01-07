package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/concurrency"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/models"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

/*
//RegisterProxy proxy register into keeper
func (s *Server) RegisterProxy(ctx context.Context, in *keeper.RegisterRequest) (*basic.BasicResponse, error) {
	if len(in.ClusterName) == 0 || len(in.Address) == 0 {
		return nil, errors.New("proxy register clusterName or Address is nil")
	}

	isLoop, _, _ := utils.IsLoopback(in.Address)
	if isLoop {
		glog.GLog.Warnf("keeper register proxy address: %s is loopback", in.Address)
		return nil, errors.New("proxy address is loopback")
	}
	//增加grp client端重试和初始化连接超时
	conn, err := newRPCConn(in.Address)
	if err != nil {
		glog.GLog.Errorf("proxy registery to keeper failed : %v", err)
		return nil, err
	}

	cli, ok := s.proxyClients.Load(in.ClusterName)
	if !ok {
		cli = make(ProxyCli)
	}
	cli.(ProxyCli)[in.Address] = proxy.NewProxyClient(conn)
	s.proxyClients.Store(in.ClusterName, cli)

	metrics.ProxyGauge.Inc()

	glog.GLog.WithFields( "clusterName", in.ClusterName, "Address", in.Address ).Infoln("proxy register success")
	out := new(basic.BasicResponse)
	out.Code = codes.OK
	return out, nil
}

//RegisterAgent agent register into keeper
func (s *Server) RegisterAgent(ctx context.Context, in *keeper.RegisterRequest) (*basic.BasicResponse, error) {
	isLoop, _, _ := utils.IsLoopback(in.Address)
	if isLoop {
		return nil, errors.New("keeper register proxy address is loopback")
	}

	conn, err := grpc.Dial(in.Address, grpc.WithInsecure())
	if err != nil {
		glog.GLog.Errorf("keeper register agent failed: %v", err)
		return nil, err
	}

	s.agentClients.Store(in.Address, agent.NewAgentClient(conn))

	glog.GLog.WithFields( "Address", in.Address ).Infoln("agent register success")

	metrics.AgentGauge.Inc()

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	return out, nil
}
*/

func (s *Server) AgentHeartbeat(ctx context.Context, in *keeper.AgentHeartRequest) (*basic.BasicResponse, error) {
	glog.GLog.WithFields("address", in.Address).Infoln("agent heartbeat.")

	datas := make(map[string]string)

	key := fmt.Sprintf(etcdtool.AgentStatus, in.Address)
	status := models.StatusInfo{
		Status: models.StatusUp,
		Time:   time.Now().Unix(),
	}

	data, err := utils.UnLoadJSON(&status)
	if err != nil {
		glog.GLog.Errorf("agent heartbeat unload status json failed: %v", err)
	} else {
		datas[key] = string(data)
	}

	key = fmt.Sprintf(etcdtool.AgentInfor, in.Address, MysqlKey)

	data, err = utils.UnLoadJSON(in.Mysql)
	if err != nil {
		glog.GLog.Errorf("agent heartbeat unload mysql json failed: %v", err)
	} else {
		datas[key] = string(data)
	}

	key = fmt.Sprintf(etcdtool.AgentInfor, in.Address, SystemKey)

	data, err = utils.UnLoadJSON(in.SystemInfo)
	if err != nil {
		glog.GLog.Errorf("agent heartbeat unload system json failed: %v", err)
	} else {
		datas[key] = string(data)
	}

	if err = s.store.Puts(datas); err != nil {
		glog.GLog.WithFields("address", in.Address).Errorln(err)
	}

	out := new(basic.BasicResponse)
	out.Code = codes.OK
	return out, nil
}

func (s *Server) ProxyHeartbeat(ctx context.Context, in *keeper.ProxyHeartRequest) (*basic.BasicResponse, error) {
	out := new(basic.BasicResponse)
	out.Code = codes.OK
	//如果proxy status为unreg 只记录日志不操作etcd
	if in.Status == models.StatusUnreg {
		glog.GLog.WithFields("clusterName", in.ClusterName, "address", in.Address, "Status", in.Status).Infoln("proxy heartbeat")
		return out, nil
	}

	//判断proxy status如果当前状态为空不允许从空变为up,记录日志不操作etcd
	statusData, err := s.store.GetData(fmt.Sprintf(etcdtool.ProxyStatus, in.ClusterName, in.Address))
	if err != nil {
		glog.GLog.WithFields("clusterName", in.ClusterName, "address", in.Address, "Status", in.Status).Errorln(err)
		return out, nil
	}

	if len(statusData) == 0 && in.Status == models.StatusUp {
		go func(clusterName, address string) {
			glog.GLog.WithFields("clusterName", clusterName, "address", address).Infoln("restart proxy.")

			conn, err := newRPCConn(address)
			if err != nil {
				glog.GLog.Errorf("keeper create proxy rpc client failed : %v", err)
				return
			}

			cli := proxy.NewProxyClient(conn)
			// 先获取心跳锁
			key := fmt.Sprintf(etcdtool.HeartbeatLock, clusterName, address)
			hearMutex, err := s.store.NewMutex(key, 10)
			if err != nil {
				glog.GLog.Errorf("keeper new hearbeat failed: %s", err.Error())
				return
			}
			// 上锁
			err = hearMutex.Lock()
			if err != nil {
				glog.GLog.Errorf("keeper hearbeat lock failed: %s", err.Error())
				return
			}
			//需要proxy 重新加载配置以保证一致
			in := new(basic.EmptyRequest)
			resp, err := cli.Restart(context.Background(), in)
			if err != nil {
				glog.GLog.WithFields("clusterName", clusterName, "address", address).Errorln(err)
			} else {
				if resp.Code != codes.OK {
					glog.GLog.WithFields("clusterName", clusterName, "address", address).Errorf("proxy restart failed, code: %d, msg: %s", resp.Code, resp.ErrMsg)
				}
			}
		}(in.ClusterName, in.Address)

		out := new(basic.BasicResponse)
		out.Code = codes.OK

		return out, nil
	}

	datas := make(map[string]string)

	key := fmt.Sprintf(etcdtool.ProxyStatus, in.ClusterName, in.Address)
	status := models.StatusInfo{
		Status: in.Status,
		Time:   time.Now().Unix(),
	}

	// 状态
	data, err := utils.UnLoadJSON(&status)
	if err != nil {
		glog.GLog.Errorf("proxy heartbeat status unload json failed: %v", err)
	} else {
		datas[key] = string(data)
	}

	glog.GLog.WithFields("clusterName", in.ClusterName, "address", in.Address, "Status", in.Status).Infoln("proxy heartbeat")

	// 集群信息
	key = fmt.Sprintf(etcdtool.ProxyInfor, in.ClusterName, in.Address)
	data, err = utils.UnLoadJSON(in.SystemInfo)
	if err != nil {
		glog.GLog.Errorf("proxy heartbeat system information unloadjson failed: %v", err)
	} else {
		datas[key] = string(data)
	}

	// 集群配置版本号
	key = fmt.Sprintf(etcdtool.ProxyVersion, in.ClusterName, in.Address)
	datas[key] = in.ConfVersion

	if err = s.store.Puts(datas); err != nil {
		glog.GLog.WithFields("clusterName", in.ClusterName, "address", in.Address, "Status", in.Status).Errorln(err)
	}

	return out, nil
}

func (s *Server) PullConfig(ctx context.Context, in *keeper.PullConfigRequest) (*keeper.PullConfigResponse, error) {
	data, err := s.store.GetShardConf(in.ClusterName, "")
	if err != nil {
		glog.GLog.Errorf("keeper pullconfig, etcd get shard config failed: %v", err)
		return nil, err
	}

	key := fmt.Sprintf(etcdtool.ConfigVersion, in.ClusterName)
	bVersion, err := s.store.GetData(key)
	if err != nil {
		glog.GLog.Errorf("keeper pullconfig, etcd get config version failed: %v", err)
		return nil, err
	} else if len(bVersion) == 0 {
		return nil, errors.New("config version is nil.")
	}

	glog.GLog.Infof("A proxy from cluster[%s] pulls config successfully.", in.ClusterName)
	out := new(keeper.PullConfigResponse)
	out.Data = data
	out.Version = string(bVersion)
	return out, nil
}

//注销某个proxy
func (s *Server) Unregister(ctx context.Context, in *keeper.Proxy) (*basic.BasicResponse, error) {
	//通过分布式锁操作防止数据并发
	out := new(basic.BasicResponse)
	session, err := concurrency.NewSession(s.store.GetClient())
	if err != nil {
		e := errors.Wrapf(err, "%v", "create cluster lock session failed")
		glog.GLog.WithFields("clusterName", in.ClusterName).Errorln(e)
		return nil, e
	}
	defer session.Close()
	lock := concurrency.NewLocker(session, etcdtool.HostLockPrefix+in.ClusterName)
	err = lock.Lock()
	defer lock.Unlock()
	if err != nil {
		if err == concurrency.ErrLockAllreadyHold {
			glog.GLog.WithFields("clusterName", in.ClusterName).Warnf("lock allready hold by other proxy")
			out.Code = codes.LockAllReadyHold
			return out, nil
		}
		e := errors.Wrap(err, "get lock failed")
		glog.GLog.WithFields("clusterName", in.ClusterName).Errorln(e)
		return nil, e
	} else {
		glog.GLog.WithFields("clusterName", in.ClusterName).Infoln("proxy get lock")
	}
	var proxyCount int
	var statusinfo models.StatusInfo
	cli := s.store.GetClient()
	ctx1, cancel := context.WithTimeout(context.Background(), TimeOut)
	resp, err := cli.Get(ctx1, fmt.Sprintf(etcdtool.StatusPrefix, in.ClusterName), clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "read data from etcd failed")
	}
	if resp.Count == 0 {
		return nil, errors.New("proxy status is nil")
	}
	for _, kv := range resp.Kvs {
		err := utils.LoadJSON(kv.Value, &statusinfo)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal proxy status data %s failed", kv.Value)
		}
		if statusinfo.Status == models.StatusUp {
			proxyCount += 1
		}
	}
	if proxyCount <= UnregNum {
		err := errors.Errorf("could't unregister proxy current proxy count %v", proxyCount)
		glog.GLog.WithFields("clusterName", in.ClusterName, "proxy", in.Address).Errorln(err)
		return nil, err
	}
	ctx2, cancel := context.WithTimeout(ctx, time.Second*3)
	_, err = cli.Delete(ctx2, fmt.Sprintf(etcdtool.ProxyStatus, in.ClusterName, in.Address))
	cancel()
	if err != nil {
		glog.GLog.Errorf("unregister proxy %s,%s failed %v ", in.ClusterName, in.Address, err)
		return nil, err
	} else {
		glog.GLog.Infof("proxy %s %s unregister success!", in.ClusterName, in.Address)
	}
	if err := s.deleteProxyCli(in.ClusterName, in.Address); err != nil {
		return nil, errors.Wrap(err, "delete Proxy Client failed")
	}
	return out, nil
}
func (s *Server) GetProxyStatus(ctx context.Context, in *keeper.Proxy) (*keeper.ProxyStatusResponse, error) {
	var statusinfo models.StatusInfo
	out := new(keeper.ProxyStatusResponse)
	cli := s.store.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	resp, err := cli.Get(ctx, fmt.Sprintf(etcdtool.ProxyStatus, in.ClusterName, in.Address))
	if err != nil {
		return nil, errors.Wrap(err, "read data from etcd failed")
	}
	if resp.Count == 0 {
		return nil, errors.New("proxy status is nil")
	}
	value := resp.Kvs[0].Value

	if err = utils.LoadJSON(value, &statusinfo); err != nil {
		return nil, errors.Wrapf(err, "unmarshal proxy status data %s failed", value)
	}
	if statusinfo.Status == "" {
		out.Status = models.StatusUnreg
		return out, nil
	}
	out.Status = statusinfo.Status
	return out, nil
}

func (s *Server) SwitchDB(ctx context.Context, in *keeper.SwitchDBRequest) (*keeper.SwitchDBResponse, error) {
	//初始化数据源切换日志
	//TODO glog.GLog.SetLogFile(SwitchDBLog, 0, 0)
	hostgroup := in.Hostgroup
	to := in.To
	from := in.From
	proxyIP := in.ProxyIP
	out := new(keeper.SwitchDBResponse)
	glog.GLog.WithFields("clusterName", in.ClusterName, "proxyIP", in.ProxyIP, "hostgroup", in.Hostgroup, "from", in.From, "to", in.To).Infoln("proxy send switch request")
	if from == to {
		e := fmt.Errorf("from index %d equal to 'to' index %d ignore it", from, to)
		glog.GLog.Errorln(e)
		out.BasicResp.Code = codes.OK
		out.BasicResp.ErrMsg = e.Error()
	}
	w, err := s.store.GetWriteIPs(in.ClusterName, hostgroup)
	if err != nil {
		return nil, err
	}
	max := len(strings.Split(strings.TrimSpace(w), ","))
	if err := s.checkSwitchIndex(from, to, max); err != nil {
		return nil, err
	}
	//检查proxy集群的状态
	status, err := s.store.GetAllStatus(in.ClusterName)
	if err != nil {
		e := errors.Wrap(err, "get cluster status failed")
		glog.GLog.WithField("clusterName", in.ClusterName).Errorln(e)
		return nil, e
	}
	if status == nil {
		e := errors.Errorf("cluster:%s doesn't exist", in.ClusterName)
		glog.GLog.Errorln(e)
		return nil, e
	}
	for _, v := range status {
		if v != models.StatusUp {
			e := errors.Errorf("cluster %s doesn't up %s", in.ClusterName, fmt.Sprintf("%s", status))
			glog.GLog.Errorln(e)
			return nil, e
		}
	}

	if err := s.CheckProxyClients(in.ClusterName); err != nil {
		return nil, err
	}
	ok, err := s.isSwitchFeq(ctx, in.ClusterName, hostgroup, proxyIP)
	if err != nil {
		e := errors.Wrap(err, "check switch frequency  failed")
		glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", in.Hostgroup).Errorln(e)
		return nil, e
	}
	if ok {
		return createSwitchDBResp(codes.SwitchDBFrequently), nil
	}
	//获取切换配置数据
	switchDBConfig, err := s.getSwitchConfig(in.ClusterName)
	if err != nil {
		e := errors.Wrap(err, "get switchDB config failed")
		glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", in.Hostgroup).Errorln(e)
		return nil, e
	}
	//判断mysql负载,如果要切换到的mysql的负载大于配置值不允许切换
	if switchDBConfig.NeedLoadCheck {
		err := s.isAgentLoadHigh(ctx, in.ClusterName, hostgroup, to, switchDBConfig.SafeLoad)
		if err != nil {
			e := errors.Wrap(err, "check agent load")
			glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", in.Hostgroup).Errorln(e)
			return nil, e
		} else {
			glog.GLog.Infoln("agent load under safelaod")
		}
	} else {
		glog.GLog.WithField("NeedLoadCheck", switchDBConfig.NeedLoadCheck).Warnln("skip agent loadcheck")
	}
	//获取host的锁
	session, err := concurrency.NewSession(s.store.GetClient())
	if err != nil {
		e := errors.Wrapf(err, "%v", "create lock session failed")
		glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", in.Hostgroup).Errorln(e)
		return nil, e
	}
	defer session.Close()
	lock := concurrency.NewLocker(session, etcdtool.HostLockPrefix+in.ClusterName+"/"+hostgroup)
	err = lock.Lock()
	defer lock.Unlock()
	if err != nil {
		if err == concurrency.ErrLockAllreadyHold {
			glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", in.Hostgroup).Warnf("lock allready hold by other proxy")
			return createSwitchDBResp(codes.LockAllReadyHold), nil
		}
		e := errors.Wrap(err, "get lock failed")
		glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", in.Hostgroup).Errorln(e)
		return nil, e
	} else {
		glog.GLog.WithFields("proxyip", proxyIP).Infoln("proxy get lock")
	}
	//获取投票
	var voteResult float64
	clis, ok := s.proxyClients.Load(in.ClusterName)
	if !ok {
		glog.GLog.WithField("clusterName", in.ClusterName).Warnln("proxyclients is nil start init proxyclients")
		err := s.CheckProxyClients(in.ClusterName)
		if err != nil {
			e := errors.Wrap(err, "init ProxyClients faield ")
			glog.GLog.WithField("clusterName", in.ClusterName).Errorln(e)
			return nil, e
		}
		clis, _ = s.proxyClients.Load(in.ClusterName)

	} else {
		for k, _ := range clis.(ProxyCli) {
			glog.GLog.Infof("proxy cluster %s client %s ready ", in.ClusterName, k)
		}

	}

	//判断binlog同步是否在
	if switchDBConfig.NeedBinlogCheck {
		//stop proxy write
		for ip, cli := range clis.(ProxyCli) {
			c, cancel := context.WithTimeout(ctx, TimeOut)
			_, err := cli.StopWritingAbility(c, &proxy.StopWriteRequest{Hostgroup: hostgroup})
			cancel()
			if err != nil {
				e := errors.Wrap(err, "stop hostgroup  write failed")
				glog.GLog.WithFields("clusterName", in.ClusterName, "proxyIP", ip, "hostgroup", hostgroup).Errorln(e)
				return nil, e
			} else {
				glog.GLog.WithField("proxyIP", ip).Infoln("stop proxy WritingAbility")
			}
		}
		//协程用于恢复此集群所有proxy的写操作直到所有的proxy的写操作被恢复.
		defer func() {
			for ip, cli := range clis.(ProxyCli) {
				go RecoverProxyWriting(in.ClusterName, hostgroup, ip, cli)
			}
		}()

		err := s.isBinlogFullySynced(ctx, in.ClusterName, hostgroup, from, to, switchDBConfig)
		if err != nil {
			e := errors.Wrap(err, "check binlog sync failed")
			glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", hostgroup).Errorln(e)
			return nil, e
		} else {
			glog.GLog.Infoln("binlog has been synchronized ")
		}
	} else {
		glog.GLog.WithField("NeedBinlogCheck", switchDBConfig.NeedBinlogCheck).Warnln("skip binlog check")
	}
	if switchDBConfig.NeedVote {
		proxyCount := len(clis.(ProxyCli))
		for ip, cli := range clis.(ProxyCli) {
			c, cancel := context.WithTimeout(ctx, time.Second*8)
			r, err := cli.GetVote(c, &proxy.GetVoteRequest{Hostgroup: hostgroup, From: from})
			cancel()
			if err != nil {
				e := errors.Wrapf(err, "get proxy %s vote failed", ip)
				glog.GLog.WithFields("clusterName", in.ClusterName, "Hostgroup", hostgroup, "To", to).Errorln(e)
				return nil, e
			}
			if r.Vote {
				glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", hostgroup, "from", from, "to", to, "proxyIP", ip, "vote", r.Vote).Infoln("proxy argree switch")
				voteResult++
			} else {
				glog.GLog.WithFields("clusterName", in.ClusterName, "hostgroup", hostgroup, "from", from, "to", to, "proxyIP", ip, "vote", r.Vote).Infoln("proxy disagree switch")
			}
		}

		//判断票数是否过半
		//TODO 投票比例定义在配置中存储在etcd中
		if voteResult <= float64(proxyCount)*switchDBConfig.VoteApproveRatio/float64(100) {
			glog.GLog.WithFields("voteResult", voteResult, "proxyCount", proxyCount).Errorln(" less than need vote")
			return createSwitchDBResp(codes.LessHalfVote), nil
		} else {
			glog.GLog.WithFields("clusterName", in.ClusterName, "voteResult", voteResult, "proxyCount", proxyCount).Infoln("proxy cluster agree switch")
		}
	} else {
		glog.GLog.WithField("NeedVote", switchDBConfig.NeedVote).Warnln("skip proxy vote")
	}
	//关闭proxy写操作等待binlog完全日志

	for ip, cli := range clis.(ProxyCli) {
		c, cancel := context.WithTimeout(ctx, TimeOut)
		_, err := cli.SwitchDataSource(c, &proxy.SwitchDatasourceRequest{Hostgroup: hostgroup, From: from, To: to})
		cancel()

		//TODO 如果设置数据库返回一个error需要报警给DBA介入处理
		if err != nil {
			e := errors.Wrap(err, "set  proxy dbsource failed")
			glog.GLog.WithFields("clusterName", in.ClusterName, "Hostgroup", hostgroup, "From", from, "To", to, "proxyIP", ip).Errorln(e)
			return nil, e
		} else {
			glog.GLog.WithField("proxyIP", ip).Infoln("switch proxy datasource success")
		}

	}
	glog.GLog.WithFields("clusterName", in.ClusterName, "Hostgroup", hostgroup, "From", from, "To", to, "proxyIP", proxyIP).Infoln("switch all proxy dbsource success")

	//TODO 记录切换事件

	//切换成功增加切换频率控制

	lease, err := s.store.GetClient().Grant(ctx, int64(switchDBConfig.Frequency))
	if err != nil {
		e := errors.Wrap(err, "add hostgroup switch frequency failed")
		glog.GLog.Errorln(e)
		return nil, e
	}
	_, err = s.store.GetClient().Put(ctx, etcdtool.SwitchDBIntervalPrefix+in.ClusterName+"/"+hostgroup, "", clientv3.WithLease(lease.ID))
	if err != nil {
		e := errors.Wrap(err, "add hostgroup switch frequency failed")
		glog.GLog.Errorln(e)
		return nil, e
	}
	//返回切换成功

	//set activeWrite
	err = s.store.SetHostGroupActiveWrite(in.ClusterName, in.Hostgroup, int(in.To))
	if err != nil {
		e := errors.Wrap(err, "set hostgroup activeWrite  failed")
		glog.GLog.Errorln(e)
		return nil, err
	}

	glog.GLog.Infoln("Cluster[%s], HostGroup[%s] switchDB successfully", in.ClusterName, in.Hostgroup)
	return createSwitchDBResp(codes.OK), nil
}

// 增加集群成员
func (s *Server) AddProxyMember(ctx context.Context, in *keeper.ProxyMemberRequest) (*basic.BasicResponse, error) {
	if err := s.store.AddCluster(in.ClusterName, in.Host, in.Port, in.Weight); err != nil {
		glog.GLog.Errorln(err)
		return nil, err
	}
	return &basic.BasicResponse{
		Code:   codes.OK,
		ErrMsg: "",
	}, nil
}

// 获取集群列表
func (s *Server) ProxyClusters(ctx context.Context, in *keeper.ClustersRequest) (*keeper.ClustersResponse, error) {
	if len(in.ClusterName) == 0 {
		return nil, errors.New("clusterName is nil.")
	}
	datas, err := s.store.GetClusters(in.ClusterName)
	if err != nil {
		return nil, err
	}
	infos := make([]*keeper.ClusterInfo, 0)
	for i := range datas {
		infos = append(infos, &keeper.ClusterInfo{
			Host:   datas[i].Host,
			Port:   datas[i].Port,
			Weight: datas[i].Weight,
		})
	}

	return &keeper.ClustersResponse{
		Infos: infos,
	}, nil
}
