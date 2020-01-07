# zeroproxy 恢复和停止服务

### zeroproxy 停止服务相关的配置

* OfflineOnLostKeeper 	  //当proxy被拒绝多次后是否向keeper注销自己
* OfflineSwhRejectedNum.  //proxy被拒绝几次后加入注销队列
* OfflineDownHostNum       //跟多少台mysql host断开后向注销自己
* OfflineRecover.                  //跟多台mysql host断开后是否下线自己默认必须为true

### zeroproxy停止服务

> **zeroproxy 启动后会在后台跑一个协程,这个协程会一直监控跟mysql的连接情况,如果无法正常连接的mysql达到一定数量,且切换多次被拒绝。zeroproxy将会从keeper注销自己,关闭mysql代理引擎,关闭checkhealth 服务.**

1. 当无法连通的mysql达到一定数量时,zeroproxy会关闭check health server 以切断lbs进来的mysql连接流量.
2. 关闭mysql引擎. 
3. 向keeper注销自己.(即keeper从etcd中删除zeroproxy)

### ZeroProxy 恢复服务

> **当proxy 和所有的mysql连接状态正常时,zeroproxy 会试图重新启动服务向zerokeeper注册自己,以提供服务.**

1. 试图恢复时发现zerokeeper中仍然保留着自己(etcd 中zeroproxy 没有被删除).开启mysql引擎，开启check health 服务, 发送up心跳。
2. 试图恢复时发现zerokeeper中zeroproxy已经被删除,重新拉取配置,重新启动.

 

