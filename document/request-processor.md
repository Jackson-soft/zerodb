zeroProxy，zeroKeeper，zeroAgent这三个子系统在运行时会产生一些网络交互，而复杂的功能就是建立在这些简单的request-response交互上，本文从角色和功能上定义了所有可能存在的交互。
![](https://ws2.sinaimg.cn/large/006tNc79ly1fnzkz7f50nj316e0cwgod.jpg)
如上图所示，交互分成了4个箭头（Arrow）

```
Request {
    body string
}

Response {
    code int
    data string
    errMsg string
}
```

Response.code
ok: 0
common error:[1,1000)
proxy error:[1000,2000)
keeper error:[2000-3000)
agent error:[3000-4000)

# Arrow 1 Client:zeroProxy Server:zeroKeeper
### zeroProxy心跳

### 获取配置
```
Request.name:"ConfigFetch"
Request.body:"{}"

Response.Processor:ConfigFetchProcessor
Response.data:"yaml content"
```

### 注册
```
Request.name:"ProxyRegister"
Request.body:"{}"

Response.Processor:ProxyRegisterProcessor
Response.data:""
```

### 注销
```
Request.name:"ProxyUnregister"
Request.body:"{
    reason: 'xxx'
}"

Response.Processor:ProxyUnregisterProcessor
Response.data:""
```

### 数据源切换
```
请求数据源切换
Request.name:"ClusterDatasourcesSwitch"
Request.body:"{
    hostGroup:'hostGroup1',
    fromIndex:0,
    toIndex:1,
    reason:"与hostGroup1:0的连接丢失"
}"

Response.Processor:ClusterDatasourcesSwitchProcessor
Response.data:""
Response.errMsg:""
```

# Arrow 2 Client:zeroKeeper Server:zeroProxy 
### 推送配置
```
推送配置
Request.name:"PushConfig"
Request.body:"{
    configId:'1',
    content:'yaml content'
}"

Response.Processor:PushConfigProcessor
Response.data:""
Response.errMsg:""
```

```
回滚配置
Request.name:"RollbackConfig"
Request.body:"{}"

Response.Processor:RollbackConfigProcessor
Response.data:""
Response.errMsg:""
```


### 注销Proxy
```
Request.name:"UnregisterProxy"
Request.body:"{
    reason:'机器磁盘出现问题'
}"

Response.Processor:UnregisterProxyProcessor
Response.data:""
Response.errMsg:""
```

### 激活Proxy
```
Request.name:"ActivateProxy"
Request.body:"{}"

Response.Processor:ActivateProxyProcessor
Response.data:""
Response.errMsg:""
```

### 数据源切换
```
获取proxy切换投票
Request.name:"ProxyVoteFetch"
Request.body:"{
    hostGroup:'hostGroup1',
    toIndex:1
}"

Response.Processor:ProxyVoteFetchProcessor
Response.data:""
Response.errMsg:""
```

```
数据源切换指令
Request.name:"SwitchProxyDatasource"
Request.body:"{
    hostGroup:'hostGroup1',
    toIndex:1
}"

Response.Processor:SwitchProxyDatasourceProcessor
Response.data:""
Response.errMsg:""
```

```
数据源切换回滚
Request.name:"RollbackProxyDatasource"
Request.body:"{
    hostGroup:'hostGroup1'
}"

Response.Processor:RollbackProxyDatasourceProcessor
Response.data:""
Response.errMsg:""
```

```
数据源停止写入
Request.name:"StopWritingProxyDatasource"
Request.body:"{
    hostGroup:'hostGroup1',
    index:0
}"

Response.Processor:RollbackProxyDatasourceProcessor
Response.data:""
Response.errMsg:""
```

```
数据源恢复写入
Request.name:"RecoverWritingProxyDatasource"
Request.body:"{
    hostGroup:'hostGroup1',
    index:0
}"

Response.Processor:RecoverWritingProxyDatasourceProcessor
Response.data:""
Response.errMsg:""
```

``
# Arrow 3 Client:zeroKeeper Server:zeroAgent
### 数据源切换
```
获取负载信息
Request.name:"LoadFetch"
Request.body:"{}"

Response.Processor:LoadFetchProcessor
Response.data:""
Response.errMsg:""
```

```
获取Binlog信息
Request.name:"BinlogFetch"
Request.body:"{
    role:"master" OR "slave"
}"

Response.Processor:BinlogFetchProcessor
Response.data:"{
    file:binlog-1213213.log
    position:12301
}"
Response.errMsg:""
```

# Arrow 4 Client:zeroAgent Server:zeroKeeper
### zeroAgent心跳

