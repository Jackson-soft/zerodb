package server

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/agent"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/models"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
)

//this is server function

//3个心跳间隔未收到心跳则将proxy的状态置为lost
func (s *Server) proxyLost(cluster, address string) {
	key := fmt.Sprintf(etcdtool.ProxyStatus, cluster, address)
	statusInfo := models.StatusInfo{
		Status: models.StatusLost,
		Time:   time.Now().Unix(),
	}
	data, err := utils.UnLoadJSON(&statusInfo)
	if err != nil {
		glog.GLog.WithFields("cluster", cluster, "address", address, "status", models.StatusLost).Errorln(err)
		return
	}

	if err = s.store.PutData(key, string(data)); err != nil {
		glog.GLog.WithFields("cluster", cluster, "address", address, "status", models.StatusLost).Errorln(err)
	}

	glog.GLog.WithFields("cluster", cluster, "address", address, "status", models.StatusLost).Infoln("proxy lost")
}

//清理已掉线的proxy
func (s *Server) deleteProxy(cluster, address string, force bool) error {
	if len(cluster) == 0 || len(address) == 0 {
		return errors.New("parameter is nil")
	}

	//删除之前判断etcd中是否是up状态，如果是up则不删除
	key := fmt.Sprintf(etcdtool.ProxyStatus, cluster, address)
	if !force {
		data, err := s.store.GetData(key)
		if err != nil {
			return errors.WithStack(err)
		}

		if len(data) == 0 {
			return errors.New("this proxy has already been shut down.")
		}

		statusData := models.StatusInfo{}
		if err = utils.LoadJSON(data, &statusData); err != nil {
			return errors.WithStack(err)
		}

		if statusData.Status == models.StatusUp {
			return nil
		}
	}

	clis, ok := s.proxyClients.Load(cluster)
	if !ok {
		return errors.New("not found proxy client")
	}

	delete(clis.(ProxyCli), address)

	keys := make([]string, 0)

	//删除状态信息
	keys = append(keys, key)

	//删除其他信息
	key = fmt.Sprintf(etcdtool.ProxyInfor, cluster, address)
	keys = append(keys, key)

	if err := s.store.Deletes(keys); err != nil {
		return errors.WithStack(err)
	}

	// 删除集群成员
	index := strings.Index(address, ":")
	host := address[:index]

	return s.store.DelCluster(cluster, host)
}

//TODO 警告信息
func (s *Server) agentLost(address string) error {
	key := fmt.Sprintf(etcdtool.AgentStatus, address)
	statusData := models.StatusInfo{
		Status: models.StatusLost,
		Time:   time.Now().Unix(),
	}
	data, err := utils.UnLoadJSON(&statusData)
	if err != nil {
		glog.GLog.WithFields("agentip", address, "status", models.StatusLost).Errorf("keeper put agent lost failed: %v", err)
		return err
	}
	if err = s.store.PutData(key, string(data)); err != nil {
		glog.GLog.WithFields("agentip", address, "status", models.StatusLost).Errorf("keeper put agent lost failed: %v", err)
		return err
	}
	glog.GLog.WithFields("agentip", address, "status", models.StatusLost).Infoln("keeper put agent lost")

	return nil
}

func (s *Server) agentDown(address string) {
	s.agentClients.Delete(address)

	key := fmt.Sprintf(etcdtool.AgentStatus, address)

	if err := s.store.Delete(key); err != nil {
		glog.GLog.WithFields("agentip", address, "status", models.StatusLost).Errorf("keeper put agent lost failed: %v", err)
		return
	}
	glog.GLog.WithFields("agentip", address, "status", models.StatusLost).Infoln("keeper put agent down")
}

//检测心跳数据
func (s *Server) checkHeartbeat() {
	go s.checkAgent()
	s.checkProxy()
}

func (s *Server) checkAgent() {
	glog.GLog.Infoln("check agent heartbeat start ...")
	agents, err := s.store.GetDatas(etcdtool.AgentPrefix)
	if err != nil {
		glog.GLog.Errorln(err)
		return
	}
	if len(agents) > 0 {
		statusOps := make([]clientv3.Op, 0)
		var interval int64
		var index int
		var address string
		infor := models.StatusInfo{}
		for key, value := range agents {
			infor = models.StatusInfo{}
			if err = utils.LoadJSON(value, &infor); err != nil {
				continue
			}
			index = strings.LastIndex(key, "/")
			if index > 0 {
				address = key[index+1:]
			} else {
				continue
			}
			interval = time.Now().Unix() - infor.Time
			if interval >= HBDown {
				// 超过丢失心跳间隔，down
				s.agentClients.Delete(address)
				// 删除 rpc client 、删除 etcd 状态及基本信息
				tkey := fmt.Sprintf(etcdtool.AgentInfor, address, MysqlKey)
				statusOps = append(statusOps, clientv3.OpDelete(tkey))

				tkey = fmt.Sprintf(etcdtool.AgentInfor, address, SystemKey)
				statusOps = append(statusOps, clientv3.OpDelete(tkey))

				statusOps = append(statusOps, clientv3.OpDelete(key))
			} else if infor.Status == models.StatusUp {
				if interval < HBLost {
					if _, ok := s.agentClients.Load(address); !ok {
						conn, err := grpc.Dial(address, grpc.WithInsecure())
						if err != nil {
							glog.GLog.WithField("address", address).Errorf("create agent rpc client failed: %v", err)
						} else {
							s.agentClients.Store(address, agent.NewAgentClient(conn))
							glog.GLog.WithField("address", address).Infoln("agent register success")
						}
					}
				} else if interval >= HBLost {
					// 超过三个心跳间隔，lost
					lost := models.StatusInfo{
						Status: models.StatusLost,
						Time:   time.Now().Unix(),
					}
					if lostData, err := utils.UnLoadJSON(&lost); err == nil {
						statusOps = append(statusOps, clientv3.OpPut(key, string(lostData)))
					} else {
						glog.GLog.Errorln("proxy lost unload json failed: %v", err)
					}
				}
			}
		}

		if len(statusOps) > 0 {
			ctx, cancel := context.WithTimeout(context.TODO(), TimeOut)
			_, err = s.store.GetClient().Txn(ctx).Then(statusOps...).Commit()
			cancel()
			if err != nil {
				glog.GLog.Errorln("check heartbeat write etcd data failed: %v", err)
				return
			}
		}
	}

	glog.GLog.Infoln("check agent heartbeat sucess ...")
}

func (s *Server) checkProxy() {
	glog.GLog.Infoln("check proxy heartbeat start ...")
	proxys, err := s.store.GetDatas(etcdtool.ProxyPrefix)
	if err != nil {
		glog.GLog.Errorln(err)
		return
	}

	if len(proxys) > 0 {
		statusOps := make([]clientv3.Op, 0)
		var interval int64
		var index int
		var want, address, cluster string
		var infor models.StatusInfo

		wantCli := sync.Map{}

		for key, value := range proxys {
			infor = models.StatusInfo{}
			if err = utils.LoadJSON(value, &infor); err != nil {
				glog.GLog.Errorln("proxy lost load json failed: %v", err)
				continue
			}
			want = strings.TrimPrefix(key, etcdtool.ProxyPrefix)
			index = strings.LastIndex(want, "/")
			if index > 0 {
				cluster = want[:index]
				address = want[index+1:]
			} else {
				continue
			}

			// 时间间隔
			interval = time.Now().Unix() - infor.Time
			if interval >= HBDown {
				// 如果状态不是「models.StatusUp」超过丢失心跳间隔，down
				// 删除 rpc client 、删除 etcd 状态及基本信息
				clis, ok := s.proxyClients.Load(cluster)
				if ok {
					pClis := clis.(ProxyCli)
					if len(pClis) == 0 {
						s.proxyClients.Delete(cluster)
					} else {
						delete(pClis, address)
					}
				}

				tkey := fmt.Sprintf(etcdtool.ProxyInfor, cluster, address)
				statusOps = append(statusOps, clientv3.OpDelete(tkey))
				glog.GLog.Infof("proxy down to add  opdelete: %s", tkey)

				statusOps = append(statusOps, clientv3.OpDelete(key))
				glog.GLog.Infof("proxy down to add  opdelete: %s", key)
			} else if infor.Status == models.StatusUp {
				if interval < HBLost {
					// 未超过Lost间隔，创建rpc client
					cli, ok := s.proxyClients.Load(cluster)
					glog.GLog.WithFields("cluster", cluster, "address", address).Infof("start create proxy rpc client: %s, ok: %v, cli: %v", key, ok, cli)
					if !ok {
						cli = make(ProxyCli)
					}
					if _, ok = cli.(ProxyCli)[address]; !ok {
						conn, err := newRPCConn(address)
						if err != nil {
							glog.GLog.Errorf("create rpc client failed : %v", err)
						} else {
							cli.(ProxyCli)[address] = proxy.NewProxyClient(conn)
							s.proxyClients.Store(cluster, cli)
							glog.GLog.WithFields("cluster", cluster, "address", address).Infof("create proxy rpc client success: %s", key)
						}
					}

				} else if interval >= HBLost {
					// 超过三个心跳间隔，lost
					lost := models.StatusInfo{
						Status: models.StatusLost,
						Time:   time.Now().Unix(),
					}
					if lostData, err := utils.UnLoadJSON(&lost); err == nil {
						statusOps = append(statusOps, clientv3.OpPut(key, string(lostData)))
					} else {
						glog.GLog.Errorln("proxy lost unload json failed: %v", err)
					}
				}
			}

			//同步keeper内存中的client
			cli, ok := s.proxyClients.Load(cluster)
			if ok {
				wCli, ok := wantCli.Load(cluster)
				if !ok {
					wCli = make(ProxyCli)
				}
				if pCli, ok := cli.(ProxyCli)[address]; ok {
					wCli.(ProxyCli)[address] = pCli
				}
				wantCli.Store(cluster, wCli)
			}
		}
		s.proxyClients = wantCli

		if len(statusOps) > 0 {
			ctx, cancel := context.WithTimeout(context.TODO(), TimeOut)
			_, err = s.store.GetClient().Txn(ctx).Then(statusOps...).Commit()
			cancel()
			if err != nil {
				glog.GLog.Errorln("check heartbeat write etcd data failed: %v", err)
			}
		}
	}

	glog.GLog.Infoln("check proxy heartbeat sucess ...")
}

//通过etcd重建rpc的客户端
func (s *Server) newRPCCliByDB() error {
	glog.GLog.Infoln("create rpc client start ....")
	agents, err := s.store.GetDatas(etcdtool.AgentPrefix)
	if err != nil {
		glog.GLog.Errorf("create rpc client failed: %v", err)
		return err
	}

	var index int
	var conn *grpc.ClientConn

	mToDelete := make([]string, 0)
	statusInfo := models.StatusInfo{}

	if len(agents) > 0 {
		for key, value := range agents {
			statusInfo = models.StatusInfo{}
			if err = utils.LoadJSON(value, &statusInfo); err != nil {
				continue
			}
			if statusInfo.Status == models.StatusUp {
				index = strings.LastIndex(key, "/")
				if index > 0 {
					addr := key[index+1:]
					//grpc.WithBlock() 主要是用来测试rpc是否能连通
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
					conn, err = grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
					cancel()
					if err == nil {
						s.agentClients.Store(addr, agent.NewAgentClient(conn))
						glog.GLog.Infof("success dial agent %s RPC client", addr)
					} else {
						//失败则加入删除队列
						mToDelete = append(mToDelete, key)
					}
				}
			} else if statusInfo.Status == models.StatusLost {
				mToDelete = append(mToDelete, key)
			}
		}
	}

	proxys, err := s.store.GetDatas(etcdtool.ProxyPrefix)
	if err != nil {
		glog.GLog.Infof("failed get proxy client: %v", err)
		return err
	}

	if len(proxys) > 0 {
		for key, value := range proxys {
			statusInfo = models.StatusInfo{}
			if err = utils.LoadJSON(value, &statusInfo); err != nil {
				continue
			}
			if statusInfo.Status == models.StatusUp {
				key = strings.TrimPrefix(key, etcdtool.ProxyPrefix)
				//index = strings.LastIndex(key, "/")
				index = strings.LastIndexByte(key, '/')
				if index > 0 {
					cluster := key[:index]
					address := key[index+1:]
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
					conn, err = grpc.DialContext(ctx, address, grpc.WithBlock(), grpc.WithInsecure())
					cancel()
					if err == nil {
						cli, ok := s.proxyClients.Load(cluster)
						if !ok {
							cli = make(ProxyCli)
						}
						if _, ok = cli.(ProxyCli)[address]; !ok {
							cli.(ProxyCli)[address] = proxy.NewProxyClient(conn)
							s.proxyClients.Store(cluster, cli)
							glog.GLog.Infof("success dial proxy %s:%s RPC client", cluster, address)
						}
					} else {
						mToDelete = append(mToDelete, key)
					}
				}
			} else if statusInfo.Status == models.StatusLost {
				mToDelete = append(mToDelete, key)
			}
		}
	}

	if len(mToDelete) > 0 {
		if err = s.store.Deletes(mToDelete); err != nil {
			glog.GLog.Warnf("new rpc client from etcd delete key failed: %v", err)
		}
	}

	glog.GLog.Infoln("create rpc client end ....")
	return nil
}

//isSwitchFeq 检测切换请求是否在时间间隔内,如果在间隔内true否则返回false
func (s *Server) isSwitchFeq(ctx context.Context, clusterName string, hostgroup string, ip string) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	if resp, err := s.store.GetClient().Get(c, etcdtool.SwitchDBIntervalPrefix+clusterName+"/"+hostgroup); err != nil {
		return false, err
	} else if resp.Count > 0 {
		return true, nil
	}
	return false, nil

}

//isLoadHigh 判断mysql负载是否过高,如果负载大于safeload返回false,否则返回true
func (s *Server) isAgentLoadHigh(ctx context.Context, clusterName string, hostgroup string, to int32, safeload float64) error {
	var ips []string
	str, err := s.store.GetWriteIPs(clusterName, hostgroup)
	if err != nil {
		return err
	}
	if len(str) == 0 {
		return errors.New("agent is nil")
	}
	for _, ip := range strings.Split(str, ",") {
		str := strings.SplitAfter(ip, ",")
		ips = append(ips, str[0]+s.conf.AgentPort)
	}
	/*	if !utils.ValidIP(ips) {
		return errors.New("agent ip valid")
	}*/
	//check agent index
	if len(ips) > int(to) {
		c, cancel := context.WithTimeout(ctx, TimeOut)
		cli, ok := s.agentClients.Load(ips[to])
		if !ok {
			cli, err = newAgentClient(ips[to])
			if err != nil {
				return err
			}
			s.agentClients.Store(ips[to], cli)
		}

		resp, err := cli.(agent.AgentClient).GetLoad(c, &agent.LoadRequest{LoadLimit: safeload})
		defer cancel()
		if err != nil {
			return err
		}
		if resp.Code == codes.OK {
			return nil
		}
		return errors.Wrapf(errors.New(resp.ErrMsg), "agent %s", ips[to])

	} else {
		return errors.Errorf("switch to index %v doesn't exist", to)
	}
}

//判断binlog的位置是否在规定的差距之内,如果大于规定的差距,返回error,否则等待完全同步,如果在规定的时间内完全同步返回nil,否则返回error
func (s *Server) isBinlogFullySynced(ctx context.Context, clusterName string, hostgroup string, from int32, to int32, config config.SwitchDB) error {

	safeBinlogDalay := config.SafeBinlogDelay
	binlogwaittime := config.Binlogwaittime
	var ips []string
	//获取指定hostgroup的所有writeip.
	str, err := s.store.GetWriteIPs(clusterName, hostgroup)
	if err != nil {
		return err
	}
	if len(str) == 0 {
		return errors.New("agent is nil")
	}
	for _, ip := range strings.Split(str, ",") {
		str := strings.SplitAfter(ip, ",")
		ips = append(ips, str[0]+s.conf.AgentPort)
	}
	//检查切换索引值,如果索引值大于write的数量返回错误.
	if len(ips) > int(to) && len(ips) > int(from) {

		cliTo, ok := s.agentClients.Load(ips[to])
		if !ok {
			cliTo, err = newAgentClient(ips[to])
			if err != nil {
				return err
			}
		}
		cliFrom, ok := s.agentClients.Load(ips[from])
		if !ok {
			cliFrom, err = newAgentClient(ips[from])
			if err != nil {
				return err
			}
		}
		c, cancel := context.WithTimeout(ctx, time.Second*time.Duration(binlogwaittime))
		defer cancel()
		for {
			select {
			case <-c.Done():
				return errors.New("binlog sync timeout  ")
			default:
			}
			r1, err := cliFrom.(agent.AgentClient).GetBinLog(c, &agent.BinLogRequest{Role: "master"})
			if err != nil {
				return errors.Wrapf(err, "get %s binlog failed", ips[from])
			}
			r2, err := cliTo.(agent.AgentClient).GetBinLog(ctx, &agent.BinLogRequest{Role: "slave"})
			if err != nil {
				return errors.Wrapf(err, "get %s binlog failed", ips[to])
			}
			//glog.GLog.WithFields("from_binlog", ips[from], "to_binlog", r2).Infoln("success get  binlog")
			glog.GLog.Infof("success get binlog master binlog:%s,%v,slave binlog %s,%v", ips[from], r1, ips[to], r2)
			if r1.File != r2.File {
				return errors.Errorf("binlog filename %s,%s doesn't match ", r1.File, r2.File)
			}
			dst := r1.Position - r2.Position
			if dst == 0 {
				return nil
			}
			if float64(dst) > safeBinlogDalay {
				return errors.Errorf("binlog delay more than safebinlogDelay,position1:%d,position1:%d", r1.Position, r2.Position)
			}
			time.Sleep(time.Millisecond * 100)
		}

	} else {
		return errors.Errorf("switch to index %v doesn't exist", to)
	}

}

//获取切换的配置
func (s *Server) getSwitchConfig(clusterName string) (config.SwitchDB, error) {
	var switchConf config.SwitchDB
	key := fmt.Sprintf(etcdtool.ShardConfig, clusterName, etcdtool.DefaultConfig, "switch")
	kvs, err := s.store.GetData(key)
	if err != nil {
		return switchConf, err
	}
	if kvs == nil {
		return switchConf, errors.New("switch config is nil")
	}
	err = utils.LoadYaml(kvs, &switchConf)
	if err != nil {
		return switchConf, err
	}
	return switchConf, err
}

//初始化clusterName下所有proxy的proxyClient
func (s *Server) initProxyClients(clusterName string) error {
	addrs, err := s.store.GetProxyAddrs(clusterName)
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if err := s.initProxyClient(clusterName, addr); err != nil {
			return err
		}

	}
	return nil
}

//初始化proxyClient
func (s *Server) initProxyClient(clusterName string, addr string) error {
	var cli ProxyCli
	if c, ok := s.proxyClients.Load(clusterName); !ok {
		cli = make(ProxyCli)
		s.proxyClients.Store(clusterName, cli)
	} else {
		cli = c.(ProxyCli)
	}
	conn, err := newRPCConn(addr)
	if err != nil {
		return err
	}
	cli[addr] = proxy.NewProxyClient(conn)
	glog.GLog.Infof("success init proxy %s:%s RPC client\n", clusterName, addr)
	return nil
}

//initagent agentClients为空时初始化agentClients
func (s *Server) initAgentClient(clusterName string) error {
	var clientLen int
	s.agentClients.Range(func(key interface{}, value interface{}) bool {
		clientLen++
		return true
	})
	if clientLen != 0 {
		return errors.New("AgentClients not nil")
	}

	agents, err := s.store.GetAllWriteIP(clusterName)
	if err != nil {
		return errors.Wrapf(err, "failed get %s cluster agent ip", clusterName)
	}
	if len(agents) == 0 {
		return errors.New("agents is nil")
	}
	/*if !utils.ValidIP(ips) {
		return errors.Errorf("valid agent ip %v ", ips)
	}*/
	for _, ips := range agents {
		for _, a := range strings.Split(ips, ",") {
			ip := strings.SplitAfter(a, ",")[0] + s.conf.AgentPort
			agentClient, err := newAgentClient(ip)
			if err != nil {
				return err
			}
			s.agentClients.Store(ip, agentClient)
		}

	}

	return nil
}

func newAgentClient(target string) (agent.AgentClient, error) {
	conn, err := newRPCConn(target)
	if err != nil {
		return nil, err
	}
	return agent.NewAgentClient(conn), nil
}

func createSwitchDBResp(code codes.Code) *keeper.SwitchDBResponse {
	return &keeper.SwitchDBResponse{
		BasicResp: &basic.BasicResponse{
			Code:   code,
			ErrMsg: codes.CodeText(code),
		},
	}

}
func newRPCConn(addr string) (*grpc.ClientConn, error) {
	//ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, DialTimeOut)
	//defer cancel()
	/*	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(MaxRetryNum),
		grpc_retry.WithPerRetryTimeout(time.Second * 5),
	}*/
	//conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	//conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "create %s rpc connection failed", addr)
	}
	return conn, nil
}

//recover proxy writeing ability
func RecoverProxyWriting(clusterName string, hostgroup string, addr string, cli proxy.ProxyClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DefaultSwitchDBInterval))
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			break
		default:
		}
		_, err := cli.RecoverWritingAbility(context.Background(), &proxy.RecoverWriteRequest{Hostgroup: hostgroup})
		if err != nil {
			glog.GLog.WithFields("proxyIP", addr, "hostgroup", hostgroup).Errorf("recover hostgroup DB write failed %v", err)
		} else {
			glog.GLog.WithFields("clusterName", clusterName, "proxyIP", addr).Infoln("Recover proxy WriteAblility success")
			break
		}
	}
}

//delete server proxy client
func (s *Server) deleteProxyCli(clustername string, addr string) error {
	clis, ok := s.proxyClients.Load(clustername)
	if !ok {
		return errors.Errorf("cluster %s doesn't exist", clustername)
	}
	c := clis.(ProxyCli)
	_, ok = c[addr]
	if ok {
		delete(c, addr)
	}
	return nil
}

//以etcd中的proxy为准,验证etcd中的Proxy和内存中的ProxyClient是否一致,如果不存在就创建,如果有多余就删除
func (s *Server) CheckProxyClients(clusterName string) error {
	//get proxy addrs
	addrs, err := s.store.GetProxyAddrs(clusterName)
	if err != nil {
		return err
	}
	clis, ok := s.proxyClients.Load(clusterName)
	if !ok && len(addrs) == 0 {
		return nil
	}
	//遍历etcd中的proxy ip
	for _, addr := range addrs {
		_, ok := clis.(ProxyCli)[addr]
		if !ok {
			err := s.initProxyClient(clusterName, addr)
			if err != nil {
				return err
			}
		}
	}
	//遍历内存中的proxycliens
	for k, _ := range clis.(ProxyCli) {
		if !utils.Contains(addrs, k) {
			if err := s.deleteProxyCli(clusterName, k); err != nil {
				return err
			}
		}
	}
	return nil
}

//checkSwitchIndex 检测数据源切换索引值.
func (s *Server) checkSwitchIndex(from int32, to int32, max int) error {
	if to < 0 || from < 0 {
		return errors.Errorf("dataSwitch index must  greater than zero [from:%d] [to %d]", from, to)
	}
	if to >= int32(max) || from >= int32(max) {
		return errors.Errorf("dataSwitch to index  from index must less than max [from:%d] [to:%d],[max:%d]", from, to, max)
	}
	return nil
}
