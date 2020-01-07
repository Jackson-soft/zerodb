# ZeroDB演示

## 集群的状态演变

### Proxy

整个Proxy的状态变迁过程是 空 -> `ready` -> `up` -> `lost` -> 空 。

#### 正常流程

空 -> `ready` -> `up` -> `lost` -> 空

#### 空 -> `up`

这种情况就是网络抖动等问题导致`proxy`在正常发`up`心跳过程中出现心跳中断，然后恢复。如果这个过程时间超过心跳丢失时间间隔`keeper`会把`proxy`的心跳状态从`etcd`中删除。
这种情况会导致`proxy`重新拉取配置并关闭代理服务重新加载配置信息。

```shell
ps -ef | grep zeroProxy
```

```shell
export ETCDCTL_API=3
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/proxy/status/cluster-func-test/
```

### Agent

相对于`proxy`的状态变迁，`agent`相对简单： 空 -> `up` -> 空。

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/agent/status/
```

## web-api的演示

### 分库分表配置的初始化

+ 强制初始化

```shell
curl -X POST \
  http://10.1.24.173:5004/api/shardconf_init \
  -H "content-type: multipart/form-data" \
  -F "clusterName=cluster-func-test" \
  -F "force=tre" \
  -F "file=@zero-proxy/proxy/test-conf/proxy_conf_test.yaml" | jq
```

+ 非强制初始化

```shell
curl -X POST \
  http://10.1.24.173:5004/api/shardconf_init \
  -H "content-type: multipart/form-data" \
  -F "clusterName=fusu-test" \
  -F "file=@zero-proxy/proxy/test-conf/proxy_conf_test.yaml" | jq
```

```shell
export ETCDCTL_API=3
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/
```

### 获取配置信息

```shell
curl -X GET http://10.1.24.173:5004/api/config?clusterName=cluster-func-test&snapshotName= | jq
```

### 创建配置快照

```shell
curl -X POST \
  http://10.1.24.173:5004/api/snapshot_config \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'clusterName=cluster-func-test&snapshotName=test_xb' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/test_xb/
```

### 快照列表

```shell
curl -X GET http://10.1.24.173:5004/api/snapshot_list?clusterName=cluster-func-test | jq
```

### 获取proxy集群的状态信息

```shell
curl -X GET \
  'http://10.1.24.173:5004/api/proxy_cluster/status?clusterName=cluster-func-test' | jq
```

### 下线某一台proxy

**注意这里的address是proxy的RPC地址（ip:port)**

```shell
curl -X PUT \
  'http://10.1.24.173:5004/api/proxy/unregister?clusterName=cluster-func-test&address=10.1.21.242:5006&reason=xx' | jq
```

非法的地址

```shell
curl -X PUT \
  'http://10.1.24.173:5004/api/proxy/unregister?clusterName=cluster-func-test&address=10.1.21.242&reason=xx' | jq
```

### 推送配置

+ 推送默认配置

```shell
curl -X PUT 'http://10.1.24.173:5004/api/push_config?clusterName=cluster-func-test' | jq
```

+ 推送某个配置快照

```shell
curl -X PUT 'http://10.1.24.173:5004/api/push_config?clusterName=cluster-func-test&snapshotName=test_xx' | jq
```

### 手动切换数据库

```shell
curl -X PUT \
  'http://10.1.24.173:5004/api/proxy_cluster/switch?hostGroup=hostGroup1&from=0&to=1&clusterName=cluster-func-test&reason=this%20is%20a%20reason.' | jq
```

### 增加hostGroup

```shell
curl -X POST \
  http://10.1.24.173:5004/api/add_hostgroups \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "cluster-func-test",
    "groups": [{
        "name": "hostGroup100",
        "max_conn": 1025,
        "init_conn": 100,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.1:3306,10.1.21.79:3306"
    }]
}' | jq
```

最大连接数太小的情况

```shell
curl -X POST \
  http://10.1.24.173:5004/api/add_hostgroups \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "cluster-func-test",
    "groups": [{
        "name": "hostGroup100",
        "max_conn": 100,
        "init_conn": 100,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.1:3306,10.1.21.79:3306"
    }]
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/hostgroups
```

### 删除HostGroup

```shell
curl -X POST \
  http://10.1.24.173:5004/api/delete_hostgroups \
  -H 'Content-Type: application/json' \
  -d '{
	"clusterName": "cluster-func-test",
	"names":["hostGroup100"]
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/hostgroups
```

### 更新StopService

```shell
curl -X POST \
  http://10.1.24.173:5004/api/update_stopservice \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "cluster-func-test",
    "service": {
        "offline_on_lost_keeper": true,
        "offline_swh_rejected_num": 4,
        "offline_down_host_num": 3,
        "offline_recover": false
    }
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/stopservice
```

### 新增Schema

```shell
curl -X POST \
  http://10.1.24.173:5004/api/add_schema \
  -H 'Content-Type: application/json' \
  -d '{
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
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/schemaconfigs
```

### 删除schema

```shell
curl -X POST \
  http://10.1.24.173:5004/api/delete_schema \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "cluster-func-test", 
    "names": [
        "xxx_test", 
        "member"
    ]
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/schemaconfigs
```

### 新增Table

```shell
curl -X POST \
  http://10.1.24.173:5004/api/add_table \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "cluster-func-test", 
    "schemaName": "empty", 
    "tables": [
        {
            "name": "xxx_text", 
            "sharding_key": "shard_key", 
            "rule": "hash"
        }
    ]
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/schemaconfigs
```

### 删除table

```shell
curl -X POST \
  http://10.1.24.173:5004/api/delete_table \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'clusterName=cluster-func-test&tableName=xxx_text&schemaName=empty' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/schemaconfigs
```

### 更新读写分离

```shell
curl -X POST \
  http://10.1.24.173:5004/api/update_shardrw \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "cluster-func-test", 
    "shardrw": {
        "zerodb": true
    }
}' | jq
```

```shell
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/schemaconfigs
```