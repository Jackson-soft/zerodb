# ETCD设计文档

## 读写权限

`etcd`的读写权限都在`keeper`。

## `proxy`的分库分表配置

分库分表的配置文件是采用分段写入`etcd`的方式。其标准格式是`zerodb/shardconf/%s/%s/%s`。

第一个占位符是`proxy`集群的名称；
第二个占位符是配置的快照名称，如果是当前配置则为`default`；
第三个占位符是配置文件的`key`。

配置文件`key`主要有以下几个:

- `basic` 基本配置项。
- `stopservice` 停止服务配置项。
- `switch` 切换配置项。
- `hostgroups` 代理服务后端`mysql`的集群配置项。
- `schemaconfigs` 代理服务的数据库分库分表配置项。

写入`etcd`的时候是由结构体解析成`YAML`文件格式的`[]byte`。解析与组装都是分段完成。

## 心跳信息

心跳设计主要是为了防止在`proxy`切换后端`mysql`主从时不一致导致数据写入不一致。

心跳的数据结构：

```go
// 状态信息
type StatusInfo struct {
	Status string `json:"status"`
	Time   int64  `json:"time"`
}
```

`Time`是写入心跳数据时的`Unix`秒级时间戳。
`Status`主要有以下三种：

```go
const (
	StatusReady = "ready"
	StatusUp    = "up"
	StatusLost  = "lost"
)
```

### `proxy`的心跳信息

`proxy`的`key`格式是`zerodb/proxy/status/%s/%s`。
第一个占位符是`proxy`的集群名称;
第二个占位符是`prxy`的`rpc`地址（ip:port);

`proxy`在创建服务之后会发一个`ready`状态心跳，然后开始解析分库分表配置，一切准备完配后开始发`up`状态心跳。
在`keeper`启动后会开启一个定时任务协程定时检查`etcd`中心跳数据中的写入时间与当前时间之差。
如果这个差值大于3个心跳间隔且状态为`up`则将状态信息修改为`lost`。
如果这个差值大于一分钟则删除该`key`。

### `agent`的心跳信息

`agent`的`key`格式是`zerodb/agent/status/%s`。占位符是`agent`的`rpc`地址（ip:port)。

`agent`的心跳与`proxy`的主要区别是没有`ready`状态，其他逻辑是一致的。

## `mysql`健康信息

`agent`在发送心跳信息的时候会附带上`mysql`的状态信息，其结构如下：
```go
type Mysql struct {
	Status    int32  //mysql当前的运行状态
	Connected int64 //mysql的连接数
	Memory    int64 //mysql当前的内存使用量
}
```