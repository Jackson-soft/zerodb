> zeroDB的一大特色是拥有一套可预测的，智能的数据源切换功能，使得对数据源的访问有高可用性。

### 负载

* 1.发起切换请求

![](http://processon.com/chart_image/5a781154e4b0615ac04c14a0.png?_=1001)

* 2.向target_agent获取负载


![](http://processon.com/chart_image/5a781316e4b064e9ddb81af7.png?_=1001)

* 3.根据负载情况决定是否切换

![](http://processon.com/chart_image/5a78138be4b0874437bbe4df.png?_=1001)


### 投票

* 1.发起切换请求

![](http://processon.com/chart_image/5a7813bbe4b064e9ddb81db7.png?_=1001)

* 2.发起切换请求

![](http://processon.com/chart_image/5a78154de4b059c41ab5d65b.png?_=1001)

* 3.唱票，并且同意

![](http://processon.com/chart_image/5a78159ae4b0615ac04c276b.png?_=1001)

* 4.发起切换请求指令

![](http://processon.com/chart_image/5a7816eee4b059c41ab5dd7e.png?_=1001)

### Binlog同步

* 1.发起切换请求

![](http://processon.com/chart_image/5a781c7ae4b064e9ddb842b2.png?_=1001)

* 2.向agent询问同步情况

![](http://processon.com/chart_image/5a781ca3e4b059c41ab5f688.png?_=1001)

* 3.proxy停止写入

![](http://processon.com/chart_image/5a781cf9e4b0874437bc0cf0.png?_=1001)

* 3.1.恢复写入

![](http://processon.com/chart_image/5a781d21e4b0615ac04c47c0.png?_=1001)

* 4.slave追上master

![](http://processon.com/chart_image/5a781dc7e4b0874437bc1027.png?_=1001)