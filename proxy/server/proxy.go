package server

import (
	"context"
	"time"

	"net"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"github.com/pkg/errors"
)

//proxy.go包含跟mysql代理服务相关的注册到server的方法
func (s *Server) GetTPS() int64 {
	return monitor.Monitor.GetTPS()
}

func (s *Server) GetQPS() int64 {
	return monitor.Monitor.GetQPS()
}

/*func (s *Server) DeleteRead(node, readAddr string) error {
    n := s.proxyEngine.GetHostGroupNode(node)
    if n == nil {
        return fmt.Errorf("invalid node [%s].", node)
    }
    return n.DeleteRead(readAddr)
}

func (s *Server) AddRead(node, readAddr string) error {
    n := s.proxyEngine.GetHostGroupNode(node)
    if n == nil {
        return fmt.Errorf("invalid node [%s].", node)
    }
    return n.DeleteRead(readAddr)
}*/

func (s *Server) DetectHostGroupAvailability() {
	round := int64(1)
	for {
		s.proxyEngine.HostGroupNodes.Range(func(key, value interface{}) bool {
			hostGroupNode, ok := value.(*backend.HostGroupNode)

			// 当心跳失败的程度达到一定值，向keeper发起数据源切换请求
			start := time.Now().UnixNano()
			err := hostGroupNode.Heartbeat()
			end := time.Now().UnixNano()

			var state string

			if err != nil {
				state = "FAILED"
				result, ok := s.hostGroupPingFailMap.Load(hostGroupNode.Cfg.Name)
				if ok {
					s.hostGroupPingFailMap.Store(hostGroupNode.Cfg.Name, result.(int)+1)
				} else {
					s.hostGroupPingFailMap.Store(hostGroupNode.Cfg.Name, 1)
				}
				//glog.Glog.Infof("hostGroup %s status %s! ", hostGroupNode.Cfg.Name, state)
				glog.Glog.Errorf("hostGroup %s status %s! ", hostGroupNode.Cfg.Name, state)
			} else {
				state = "OK"
				//added by huajia
				result, ok := s.hostGroupPingFailMap.Load(hostGroupNode.Cfg.Name)
				if ok && result.(int) > 0 {
					hostGroupNode.InitWrite(hostGroupNode.Cfg.Write, true)
					hostGroupNode.InitRead(hostGroupNode.Cfg.Read, true)
				}
				//end added
				s.hostGroupPingFailMap.Store(hostGroupNode.Cfg.Name, 0)
			}

			// 暂停打印
			if 1 > 2 {
				glog.Glog.Infof("HostGroup: [%s]. ActiveDataSource: [%v], Ping time: [%v]ms, Round:[%v], State: [%s]",
					hostGroupNode.Cfg.Name,
					hostGroupNode.GetActivedWriteIndex(),
					(end-start)/int64(time.Millisecond),
					round,
					state)
			}

			result, ok := s.hostGroupPingFailMap.Load(hostGroupNode.Cfg.Name)
			//modified by huajia
			if ok && result.(int) >= 3 && hostGroupNode.Cfg.EnableSwitch {
				//if ok && result.(int) >= 3 {
				//end modified
				/*				if s.status == StatusUnreg {
								return true
							}*/
				//glog.DSSwitchLog.Warnf("HostGroup: [%s]. The current dataSource can't be reached for 3 times. Prepare for switch request.", hostGroupNode.Cfg.Name)
				// 发起数据源切换请求
				fromIndex := hostGroupNode.GetActivedWriteIndex()
				toIndex := hostGroupNode.GetNextWriteIndex()

				glog.Glog.Warnf("HostGroup: [%s]. The current datasource can't be reached for 3 times. Prepare for switch request from: %d to %d.",
					hostGroupNode.Cfg.Name,
					fromIndex,
					toIndex)
				req := new(keeper.SwitchDBRequest)
				req.ClusterName = s.conf.ClusterName
				req.From = int32(fromIndex)
				req.To = int32(toIndex)
				req.Hostgroup = hostGroupNode.Cfg.Name
				req.ProxyIP = utils.GetLocalIP()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
				resp, err := s.keeperCli.SwitchDB(ctx, req)
				cancel()
				if err != nil || resp.BasicResp == nil {
					result, ok := s.hostGroupSwtRejectMap.Load(hostGroupNode.Cfg.Name)
					if ok {
						s.hostGroupSwtRejectMap.Store(hostGroupNode.Cfg.Name, result.(int)+1)
					} else {
						s.hostGroupSwtRejectMap.Store(hostGroupNode.Cfg.Name, 1)
					}
					// 请求失败
					glog.Glog.Errorf("request error %v", err)
				} else if resp.BasicResp.Code != codes.OK {
					// 由于各种原因被拒绝
					glog.Glog.Warnf("HostGroup: %s . Request reject: errCode: %d, errMsg: %s", hostGroupNode.Cfg.Name, resp.BasicResp.Code, resp.BasicResp.ErrMsg)

					result, ok := s.hostGroupSwtRejectMap.Load(hostGroupNode.Cfg.Name)
					if ok {
						s.hostGroupSwtRejectMap.Store(hostGroupNode.Cfg.Name, result.(int)+1)
					} else {
						s.hostGroupSwtRejectMap.Store(hostGroupNode.Cfg.Name, 1)
					}
				} else {
					// 成功请求，并且切换
					glog.Glog.Warnf("HostGroup: %s . datasource switched: from: %d to %d", hostGroupNode.Cfg.Name, fromIndex, toIndex)
					//成功切换后刷新hostGroupSwtRejectMap
					if _, ok := s.hostGroupSwtRejectMap.Load(hostGroupNode.Cfg.Name); ok {
						s.hostGroupSwtRejectMap.Delete(hostGroupNode.Cfg.Name)
					}
				}
			} else {
				if _, ok := s.hostGroupSwtRejectMap.Load(hostGroupNode.Cfg.Name); ok {
					s.hostGroupSwtRejectMap.Delete(hostGroupNode.Cfg.Name)
				}
			}
			return true
		})
		round++
		time.Sleep(2 * time.Second)
	}
}

//DetectStopService 方法负责在满足条件时下线proxy,向keeper注销自己
func (s *Server) DetectStopService(stopDssCh chan struct{}) {
	conf := s.shardConfig.currentConf.StopService
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-stopDssCh:
			return
		case <-t.C:
			{
				unConnectHostNum := 0
				s.hostGroupSwtRejectMap.Range(func(key interface{}, value interface{}) bool {
					if value.(int) >= conf.OfflineSwhRejectedNum {
						unConnectHostNum += 1
						glog.Glog.Infof("%s switch rejetct %d times\n", key.(string), value.(int))
					}
					return true
				})
				if unConnectHostNum >= conf.OfflineDownHostNum && unConnectHostNum > 0 {
					if s.proxyEngine.GetState() != StatusOffline {
						glog.Glog.Warnf("proxy can not connect %d hostgroup,prepare offline self \n", unConnectHostNum)
						if err := s.offline(); err != nil {
							glog.Glog.Errorf("proxy offline failed %v ", err)
						} else {
							glog.Glog.Infoln("proxy success offline ")
						}
					}
					if s.status != StatusUnreg {
						if err := s.unregister(); err != nil {
							glog.Glog.Errorf("proxy unreg failed %v", err)
						} else {
							glog.Glog.Infoln("proxy success unreg ")
						}
					}
				}
			}
		}
	}
}

//恢复proxy service
func (s *Server) RecoverProxySerivce(rpsCh chan struct{}) {
	allowRecover := s.shardConfig.currentConf.StopService.OfflineRecover
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-rpsCh:
			return
		case <-t.C:
			//default:
			{
				if !allowRecover {
					continue
				}
				allHostConnected := true
				s.hostGroupPingFailMap.Range(func(key, value interface{}) bool {
					if value.(int) > 0 {
						//glog.Glog.Infof("can't connect %s ping failed %d times", key, value)
						glog.Glog.Errorf("can't connect %s ping failed %d times", key, value)
						allHostConnected = false
						return allHostConnected
					}
					return true
				})
				if !allHostConnected {
					continue
				}
				if s.proxyEngine.GetState() == StatusOffline {
					glog.Glog.Infoln("all hostgroup connected try recover self")
					//恢复proxyEngine
					if err := s.online(); err != nil {
						glog.Glog.Errorf("recover proxy service failed %v", err)
					} else {
						s.status = StatusUp
						glog.Glog.Infoln("recover proxy service success")
					}
				}
			}
		}
	}
}

//向keeper注销自己
func (s *Server) unregister() error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	req := CreeatProxyReq(s.conf.ClusterName, s.conf.RPCServer)
	_, err := s.keeperCli.Unregister(ctx, req)
	if err != nil {
		return err
	}
	s.status = StatusUnreg
	return nil
}

//通过关闭checkhealth来切断lbs代理进来的mysql流量
func (s *Server) offline() error {
	conf := s.shardConfig.currentConf.StopService
	if !conf.OfflineOnLostKeeper {
		return errors.New("proxy configuration skip offline ")
	}
	/*	if s.proxyEngine.GetState() == StatusOffline {
		return nil
	}*/
	//关闭mysql引擎
	s.proxyEngine.CloseProxyEngine()
	s.proxyEngine.SetState(StatusOffline)
	return nil
}

//通过打开checkhealth使proxy online
func (s *Server) online() error {
	//开启mysql引擎
	go s.proxyEngine.ServeInBlocking(func() {
	})

	s.proxyEngine.SetState(StatusOnline)
	glog.Glog.Infoln("mysql online")
	return nil
}

func CreeatProxyReq(clusterName string, address string) *keeper.Proxy {
	req := new(keeper.Proxy)
	req.ClusterName = clusterName
	req.Address = address
	loop, _, port := utils.IsLoopback(req.Address)
	if loop {
		ip := utils.GetLocalIP()
		req.Address = net.JoinHostPort(ip, port)
	}
	return req
}
