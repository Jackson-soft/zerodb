# ZeroDB 管理端API文档

## API参数说明

- 采用`restful`风格`API`设计。
- 假设`ZeroDB`的`Web Server IP`和端口是:`127.0.0.1:10087`，`web`用户名:`admin`，密码:`admin`
- 统一的返回格式：`code`参考「request-processor」文档
- 返回数据格式为*json*，其结构如下：

``` go
{
    code: int // 状态码 1成功 0失败
    errorCode: int // 错误码
    data: interface{} //范型数据
    message: string // 信息
}
```

## API接口导航

- [登录](#登录)
- [查看proxy集群的状态](#查看proxy集群的状态)
- [注销proxy](#注销proxy)
- [激活proxy](#激活proxy)
- [初始化分库分表配置](#初始化分库分表配置)
- [查看分库分表配置](#查看分库分表配置)
- [查看集群列表](#查看集群列表)
- [导出配置文件](#导出配置文件)
- [推送配置](#推送配置)
- [回滚配置](#回滚配置)
- [创建分库分表配置快照](#创建分库分表配置快照)
- [推送分库分表配置](#推送分库分表配置)
- [查看分库分表配置快照列表](#查看分库分表配置快照列表)
- [获取HostGroups](#获取HostGroups)
- [proxy集群数据源切换](#proxy集群数据源切换)
- [增加HostGroups](#增加HostGroups)
- [删除HostGroups](#删除HostGroups)  //删除HostGroup时先删除引用
- [更新StopService](#更新StopService)
- [增加Schema](#增加Schema)
- [删除Schema](#删除Schema)  //删除Schema时先删除所有挂在该Schema下的所有Table
- [增加Table](#增加Table)
- [删除Table](#删除Table)  //需要锁住proxy的读写
- [更新读写分离](#更新读写分离)

### 登录

Action: POST  
URL: /login

参数：

```text
userName: admin
password: admin
```

返回：

```json
{
    "code": 0,
    "message": "",
    "data": {"token": "xx"}
}
```

### 查看proxy集群的状态

Action: GET  
URL: /api/proxy_cluster/status

参数：

```go
clusterName: string
```

返回：

```text
{
    code: 0
    message: ""
    data: {
        ip: {
            cpuLoad: float64
            memLoad: float64
            loadAvg: {
                avg1min: float64
                avg5min: float64
                avg15min: float64
            }
        }
    }
}
```

### 初始化分库分表配置

Action: POST  
Content-Type: multipart/form-data  
URL: /api/shardconf_init

参数：

```go
clusterName: string //集群名称
force: string //是否强制初始化,为空时有默认配置时不做初始化动作
file: string  //指定配置文件,上传的yaml格式的配置文件
```

返回：

```json
{
    "code": 0,
    "message": "",
    "data": ""
}
```

### 查看分库分表配置

Action: GET
URL: /api/config

参数:

```go
clusterName: string //集群名称
snapshotName: string // 快照名称
```

返回：

```json
{
    "code": 0,
    "message": "",
    "data": "" //object
}
```

### 查看集群列表

Action: GET  
URL: /api/cluster_list

返回:

```json
{
    "code":0,
    "message":"",
    "data": "" //[]string
}
```

### 导出配置文件

Action: GET  
URL: /api/cluster_list

参数:

```go
clusterName: string //集群名称
snapshotName: string // 快照名称
```

返回:

```json
{
    "code":0,
    "message":"",
    "data": "" //string  文件名，yaml格式
}
```

### 移除proxy

Action: PUT  
URL: /api/proxy/unregister

参数:

```go
address: string //proxy的address=>ip:port
clusterName: string //集群名称
reason: string  //简要的理由
```

返回结果：

```json
{
    "code":0,
    "message":"",
    "data": ""
}
```

### 推送配置

Action: PUT  
URL: /api/push_config
Timeout: 60sec

参数：

```go
snapshotName: string //配置快照，为空则推送默认配置
clusterName: string //集群名称
```

返回结果：

```json
{
    "code":0,
    "message":"",
    "data": ""
}
```

### proxy集群数据源切换

Action: PUT  
URL: /api/proxy_cluster/switch
Timeout: 60sec

参数:

```go
hostGroup: string //hostGroup名字
from: int //切出的host的index
to: int //切到hostGroup对应host的index
clusterName: string //集群名称
reason: string //简要的理由，必须参数
```

返回:

```json
{
    "code":0,
    "message":"",
    "data": ""
}
```

### 创建分库分表配置快照

Action: POST  
URL: /api/snapshot_config

参数：

```go
clusterName: string
snapshotName: string
```

返回：

```json
{
    "code": 0,
    "message": "",
    "data": ""
}
```

### 查看分库分表配置快照列表

Action: GET  
URL: /api/snapshot_list

参数：

```go
clusterName: string
```

返回：

```json
{
    "code": 0,
    "message": "",
    "data": ["xxxx", "xxxx1"]
}
```

### 获取HostGroups

Action: GET
URL: /api/hostgroups

参数：

```go
clusterName: string
snapshotName: string
```

返回：

```json
{
    "code": 0,
    "message": "",
    "data": ["xxxx", "xxxx1"]
}
```

### 增加HostGroups

Action: POST  
URL: /api/add_hostgroups  
Content-Type: application/json
Timeout: 60sec

参数：

```text
{
    clusterName: string
    groups: [{
        name: string
        max_conn: in
        init_conn: int
        user: string
        password: string
        write: string
        read: string
        down_after_noalive: bool
        switch: {
            need_vote: bool
            need_load_check: bool
            need_binlog_check: bool
        }]
    }
}
```

返回结果：

```json
{
    "code":0,
    "message":"",
    "data": map[string]string //如果有失败的，则返回key:ip,value:ErrorMsg 的map
}
```

示例入参：

```json
{
    "clusterName": "cluster_test",
    "groups": [{
        "name": "hostGroup1000",
        "max_conn": 10000,
        "init_conn": 100,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.1:3306,10.1.21.79:3306"
    }]
}
```

### 删除HostGroups

Action: POST  
URL: /api/delete_hostgroups  
Content-Type: application/json
Timeout: 60sec

参数：

```go
{
    clusterName: string
    names: []string
}
```

返回结果：

```json
{
    "code":0,
    "message":"",
    "data": map[string]string //如果有失败的，则返回key:ip,value:ErrorMsg 的map
}
```

示例入参：

```json
{
    "clusterName": "cluster_test",
    "names":["hostGroup1000", "hostGroup6"]
}
```

### 更新StopService

Action: POST  
URL: /api/update_stopservice  
Content-Type: application/json

参数：

```text
{
    clusterName: string
    service: {
        unreg_on_swh_rejected: bool
        unreg_rejected_num: int
        unreg_down_host_num: int
        offline_on_lost_keeper: bool
        offline_swh_rejected_num: int
        offline_down_host_num: int
        offline_recover: bool
    }
}
```

返回结果：

```json
{
    "code":0,
    "data": "",
    "message":""
}
```

参数示例：

```json
{
    "clusterName": "cluster_test",
    "service": {
        "unreg_on_swh_rejected": true,
        "unreg_rejected_num": 4,
        "unreg_down_host_num": 3,
        "offline_on_lost_keeper": true,
        "offline_swh_rejected_num": 4,
        "offline_down_host_num": 3,
        "offline_recover": false
    }
}
```

### 新增Schema

Action:POST  
URL:/api/add_schema  
Content-Type: application/json
Timeout: 60sec

参数：

```text
{
    clusterName: string
    schemas:[{
        name: string
        custody: bool
        nonsharding_host_group: string
        sharding_host_groups: []string
        schema_sharding: int
        table_sharding: int
        rw_split: bool
        table_configs: [{
            name: string
            sharding_key: string
            rule: string
        }]
    }]
}
```

返回结果：

```json
{
    "code":0,
    "data":"",
    "message":""
}
```

示例入参：

```json
{
    "clusterName": "cluster_test",
    "schemas": [
        {
            "Name": "xxx_test",
            "custody": true,
            "nonsharding_host_group": "hostGroup1",
            "sharding_host_groups": [
                "hostGroup1",
                "hostGroup2",
                "hostGroup3",
                "hostGroup4"
            ],
            "schema_sharding": 64,
            "table_sharding": 2,
            "table_configs": [
                {
                    "name": "ship",
                    "sharding_key": "entity_id",
                    "rule": "string"
                }
            ]
        }
    ]
}
```

### 删除schema

Action:POST  
URL:/api/delete_schema  
Content-Type: application/json
Timeout: 60sec

参数：

```go
{
    clusterName: string
    names: []string
}
```

返回结果：

```json
{
    "code":0,
    "message":"",
    "data": map[string]string //如果有失败的，则返回key:ip,value:ErrorMsg 的map
}
```

示例入参：

```json
{
    "clusterName": "cluster_test",
    "names": [
        "xxx_test",
        "member"
    ]
}
```

### 新增Table

Action: POST  
URL: /api/add_table  
Content-Type: application/json
Timeout: 60sec

参数：

```go
{
    clusterName: string
    schemaName: string
    tables: [{
        name: string
        sharding_key: string
        rule: string
    }]
}
```

返回结果：

```json
{
    "code":0,
    "data": "",
    "message":""
}
```

示例入参：

```json
{
    "clusterName": "cluster_test", 
    "schemaName": "empty", 
    "tables": [
        {
            "name": "xxx_text", 
            "sharding_key": "shard_key", 
            "rule": "hash"
        }
    ]
}
```

### 删除table

Action: POST  
URL: /api/delete_table
Timeout: 60sec

参数：

```go
clusterName: string
tableName：string //table的唯一标识
schemaName： string //schema的唯一标识
```

返回结果：

```json
{
    "code":0,
    "data":"map[string]string",
    "message":""
}
```

### 更新读写分离

Action: POST  
URL: /api/update_shardrw  
Content-Type: application/json

参数：

```go
{
    clusterName: string
    shardrw: map[string]bool // key:schemaName value: rw
}
```

返回结果：

```json
{
    "code":0,
    "data": "",
    "message":""
}
```

示例入参：

```json
{
    "clusterName": "cluster_test",
    "shardrw": {
        "hostGroup1": true, 
        "hostGroup2": false
    }
}
```
