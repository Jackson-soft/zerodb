> ZeroDB是一个支持MySQL协议的分布式分库分表中间件。

### zeroProxy
是基于MySQL协议的分库分表引擎，有以下诸多优点：支持兼容Cobar的分库，分表，读写分离，Proxy->MySQL水位控制，Route执行计划

### zeroKeeper
zeroDB的集群控制中心，负责收集zeroProxy、zeroAgent心跳，协调zeroProxy数据源切换，管理zeroProxy集群状态，管理zeroProxy配置，Web API中心

### zeroAgent
MySQL的agent，负责收集MySQL的binlog信息，收集机器负载


![](http://processon.com/chart_image/5a4391bfe4b0daa64fecd228.png?_=111)