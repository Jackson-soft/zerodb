> 在一个分布式集群中，任何一个节点都是有可能宕机的，zeroDB内部的所有子系统（zeroProxy，zeroKeeper，zeroAgent，MySQL，LBS）当然也不例外。那当少量的节点宕机时，整个系统需要继续提供完整或者部分的服务，这种能力或者机制我们一般称为「容错」（fault-tolerant）。那在系统设计初期，周全的分布式系统容错设计将会是一个不可省略性的环节，因为我们要清楚，任何的计算节点都是可能宕机的，计算节点内部的进程也是有可能OOM的，一次网络调用可能失败或者超时。
本文将会讨论zeroDB这个分布式数据库中间件对各种错误情况做出的各种反应。所谓「错误」，其实是一个非常模糊的概念，一次网络请求的超时，在极端严格的情况下，它也能被称之为「错误」，但在宽松的情况下，它只能说是个小的运行异常，达不到去自我恢复的标准。zeroDB在设计时更多得是将「错误」的鉴定和对「错误」的处理的决策权力交给系统的Admin，通过配置的形式。

### proxy

* CASE 1

```
need_vote: true
need_vote.approve_ratio: 50
unregister_when_switch_rejected: true
unregister_switch_rejected_number: 3
unregister_when_switch_rejected.down_host_num: 2
```
![](https://processon.com/chart_image/5a77c423e4b064e9ddb7172b.png?_=100)

proxy向keeper发送数据源切换请求，但由于投票获得同意的票数少于{need_vote.approve_ratio}，所以切换请求被驳回。proxy自己也会定期检查，发现自己的没有{unregister_when_switch_rejected.down_host_num}的dataHost连续请求被驳回的次数超过{unregister_switch_rejected_number}次，所以只做告警。不切换，也不发起注销请求。


* CASE 2

```
need_vote: true
need_vote.approve_ratio: 50
unregister_when_switch_rejected: true
unregister_switch_rejected_number: 3
unregister_when_switch_rejected.down_host_num: 2
```
    
![](http://processon.com/chart_image/5a77c462e4b0615ac04b1b41.png?_=10)

proxy向keeper发送数据源切换请求，但由于投票获得同意的票数少于{need_vote.approve_ratio}，所以切换请求被驳回。proxy自己也会定期检查，发现自己的有超过{unregister_when_switch_rejected.down_host_num}的dataHost连续请求被驳回的次数超过{unregister_switch_rejected_number}次，所以向keeper发起UnregisterRequest，请求让自己从集群中注销掉。

* CASE 3

![](http://processon.com/chart_image/5a77c48be4b064e9ddb718b8.png?_=10)

proxy3失联，这时proxy1,proxy2,proxy4都连不上MySQL_1,向keeper发送数据源切换请求，由于keeper没法获得proxy3的具体状态，这个时候的任何切换都需要拒绝

* CASE 4

```
need_vote: true
need_vote.approve_ratio: 50
unregister_when_switch_rejected: true
unregister_switch_rejected_number: 3
unregister_when_switch_rejected.down_host_num: 2
offline_when_lost_keeper: true
offline_switch_rejected_number: 3
# 至少有两个Host连续被数据源切换失败后，让自己下线
offline_when_lost_keeper.down_host_num: 4
self_back_online: true
```

![](http://processon.com/chart_image/5a77c4c3e4b0615ac04b1ce3.png?_=10)

proxy3发送数据源切换请求，但由于根本连不上keeper，所以没能切换成功。由于自己与keeper已经失联，所以让自己注销已经不可能了。proxy定时检测发现，发现自己的有超过{offline_when_lost_keeper.down_host_num}的Host连续请求失败的次数超过{offline_switch_rejected_number}次，这时proxy3让自己offline，关闭CH。这种offline是不是永久的，取决于{self_back_online}的配置，这里为true，所以检测到与所有的Host都能连接后，再次上线。


### keeper

* CASE 1

![](http://processon.com/chart_image/5a5c973de4b0c090524397b8.png?_=100)

所有proxy「按兵不动」

### mysql or agent

* CASE 1

```
need_vote: true
need_vote.approve_ratio: 50
need_target_load_check: true
need_target_load_check.safe_load: 2
# need_binlog_check: false
need_binlog_check: true
need_binlog_check.safe_binlog_delay: 1000
```

![](http://processon.com/chart_image/5a77ccbde4b064e9ddb73a31.png?_=100)

MySQL_1机器关机，MySQL1和zeroAgent同时死亡。发起数据源切换请求，投票获得同意的票数高于{need_vote.approve_ratio}，由于{need_target_load_check}配置为true，keeper需要向mysql_11的agent请求负载，假设低于{need_target_load_check.safe_load}，接着继续查看配置{need_binlog_check}，为true的话需要查看binlog的延迟情况。发现{need_binlog_check.safe_binlog_delay}有配置，这个表示严格的binlog对比，即需要查看source agent和target agent的binlog position，但是source agent已经挂了，所以这时也不能切换


* CASE 2

```
    need_vote: true
    need_vote.approve_ratio: 50
    need_target_load_check: true
    need_target_load_check.safe_load: 2
    # need_binlog_check: false
    need_binlog_check: true
    need_binlog_check.safe_binlog_delay: 1000
```

![](http://processon.com/chart_image/5a77cde6e4b064e9ddb73f2f.png?_=100)

MySQL_1机器关机，MySQL1和zeroAgent同时死亡。发起数据源切换请求，投票获得同意的票数高于{need_vote.approve_ratio}，由于{need_target_load_check}配置为true，keeper需要向mysql_11的agent请求负载，假设低于{need_target_load_check.safe_load}，接着继续查看配置{need_binlog_check}，为true的话需要查看binlog的延迟情况。如果target agent挂了，那就拒绝请求。