在keeper管理proxy集群的过程中，keeper需要记录并且管理proxy的状态。

![](http:3000->//processon.com/chart_image/5a72886ce4b0615ac0451d81.png_=10)

举一个完整的例子来说明proxy状态转移，有四台机器(或者容器)需要启动proxy,它们的端口都是3000，分别是
10.25.0.11:3000
10.25.0.12:3000
10.25.0.13:3000
10.25.0.14:3000


* Keeper中有对于proxy的记录，由于proxy还未启动，所以keeper中proxy中是空的。
```
```

* 将前面三台proxy先后启动，proxy立刻把自己的状态设置为「ready」，然后通过与keeper建立连接并且发送「ready」心跳，这是一个比较短暂的中间过程。
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->ready
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->ready
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->ready
```

* 只有当顺利建立心跳了，然后向keeper拉取「启动配置」后，开始初始化数据库连接池。初始化完后，将自己的状态设置为「up」，发送「up」心跳。
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
```

* 运行一段时间后，keeper和「10.25.0.13」之间的网络中断了，proxy还在发送无效的「up」心跳，keeper的定时器检测到后，keeper将它的状态修改成「lost」，
* 这时任何「数据源切换」或者「推送配置」的操作都将无法进行。
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->lost
```

* 当keeper和「10.25.0.13」的网络恢复后，「10.25.0.13」的状态又直接恢复成「up」
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
```

* 「10.25.0.14」启动，keeper收到了「10.25.0.14」的「ready」心跳。这时任何「数据源切换」或者「推送配置」的操作都是无法进行的。
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.14:3000->ready
```

*「10.25.0.14」不幸没有顺利启动，进程死亡，keeper检测到「10.25.0.14」太久都是「ready」，并且没有新的心跳，直接将「10.25.0.14」删除
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
```

*「10.25.0.14」由于明确的情况启动失败，那这种就和上面描述的情况有所不同，当明确的启动失败后，proxy向keeper发送一个「down」心跳，keeper如果收到，
直接将处于「ready」状态的「10.25.0.14」删除掉，这时就不用依赖keeper自身的定时任务来清理长久的「ready」状态的proxy，提早将proxy集群清理完毕 // 第一期没有计划实现
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
```


* 管理员重新将「10.25.0.14」启动，并且成功了
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.14:3000->up
```

* 「10.25.0.14」与MySQL连接出现问题，数量超过阈值，触发offline，向keeper发送「注销proxy」的命令，keeper将对应的proxy从列表中删除，
* 删除后不影响所有剩余「up」状态的proxy的「数据源切换」或者「推送配置」功能
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.13:3000->up
```

* proxy「10.25.0.13」万一是被恶意关闭的，那么keeper的定时器检测到后，也当它是「lost」，此时也是不能进行数据源切换的。如果真的想要把「10.25.0.13」
* 关闭，只能通过API「remove proxy」，将该proxy彻底从test_cluster中删除。
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.14:3000->up
```

* proxy「10.25.0.14」与keeper失联，keeper检测到后，将「10.25.0.14」设置为「lost」
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.14:3000->lost
```

* proxy「10.25.0.14」的网络迟迟不恢复，直到keeper将「lost」的「10.25.0.14」删除掉，但「10.25.0.14」复活了，这时「10.25.0.14」发送的「up」心跳被keeper
重新感知到，但「10.25.0.14」已经被集群移除，所以keeper给proxy发送一个「重启」命令，让「10.25.0.14」重新获取配置
```
/zerodb/proxy/test_cluster/status/10.25.0.11:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.12:3000->up
/zerodb/proxy/test_cluster/status/10.25.0.14:3000->up
```