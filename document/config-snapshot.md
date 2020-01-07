> 配置管理与分发

![](http://processon.com/chart_image/5b826f6ae4b06fc64ad58096.png?_=1001)

### 启动keeper，给自己的集群取一个名字：presentation_cluster，初始化一份空白的配置（当然可以初始化一份精心准备的，这里为了演示效果）
```
curl -X POST \
  http://10.1.24.173:5004/api/shardconf_init \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F clusterName=presentation_cluster \
  -F force=ssss \
  -F file=@/Users/eric/Code/go_proj/src/git.2dfire-inc.com/platform/zerodb/zero-proxy/proxy/test-conf/simple.yaml
```

### 启动proxy，指定proxy集群名字为：presentation_cluster
app.yaml
```
clusterName: presentation_cluster
rpcServer: 0.0.0.0:50061
proxyServer: 0.0.0.0:9696
keeperAddr: 10.1.24.173:5003
metricsAddr: 0.0.0.0:50081
checkhealthServer: 0.0.0.0:9001
charset: utf8mb4
debugServer: 0.0.0.0:50012
debugMode: true
logSQL: true
log:
  path: /var/log/zlog1
  level: info
```

### 保存下配置，给它创建一个快照（snapshot）：inital_config_ss
```
curl -X POST \
  http://10.1.24.173:5004/api/snapshot_config \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F snapshotName=inital_config_ss \
  -F clusterName=presentation_cluster
```

### 查看下快照
```
curl -X GET \
  'http://10.1.24.173:5004/api/snapshot_list?clusterName=presentation_cluster'
```

### 给配置添加5个HostGroup
```
curl -X POST \
  http://10.1.24.173:5004/api/add_hostgroups \
  -H 'Content-Type: application/json' \
  -d '{
    "clusterName": "presentation_cluster",
    "groups": [{
        "name": "hostGroup1",
        "max_conn": 2560,
        "init_conn": 10,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.1:3306,10.1.21.79:3306",
        "active_write": 0
    },{
        "name": "hostGroup2",
        "max_conn": 2560,
        "init_conn": 10,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.2:3306,10.1.21.80:3306",
        "active_write": 0
    },{
        "name": "hostGroup3",
        "max_conn": 2560,
        "init_conn": 10,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.3:3306,10.1.21.81:3306",
        "active_write": 0
    },{
        "name": "hostGroup4",
        "max_conn": 2560,
        "init_conn": 10,
        "user": "zerodb",
        "password": "zerodb@552208",
        "write": "10.1.22.4:3306,10.1.21.82:3306",
        "active_write": 0
    },{
        "name": "hostGroupCommon",
        "max_conn": 1024,
        "init_conn": 1,
        "user": "twodfire",
        "password": "123456",
        "write": "common101.my.2dfire-daily.com:3306",
        "read": "common101.my.2dfire-daily.com:3306@2,common101.my.2dfire-daily.com:3306@3",
        "active_write": 0
    }]
}'
```

### 再保存下配置，给它创建一个快照：hostGroup_added_ss
```
curl -X POST \
  http://10.1.24.173:5004/api/snapshot_config \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F snapshotName=hostGroup_added_ss \
  -F clusterName=presentation_cluster
```

### 查看当前的配置和每个snapshot的具体配置信息
```
curl -X GET \
  'http://10.1.24.173:5004/api/config?clusterName=presentation_cluster&snapshotName=' \
  -H 'Cache-Control: no-cache' \
  -H 'Postman-Token: de889de3-4c7d-4674-a7b6-19d7a1326643'
```


### 创建schema：zerodb 和附带的 table
```
curl -X POST \
  http://10.1.24.173:5004/api/add_schema \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: bb9ab927-c4d7-40c2-a847-d98a06814c51' \
  -d '{
    "clusterName": "presentation_cluster",
    "schemas": [
        {
            "Name": "zerodb",
            "sharding_host_groups": [
                "hostGroup1","hostGroup2","hostGroup3","hostGroup4"
            ],
            "rw_split": false,
            "multi_route": true,
    		"init_conn_multi_route": true,
            "schema_sharding": 128,
            "table_sharding": 1,
            "table_configs": [
                {
                    "name": "order_ins",
                    "sharding_key": "entity_id",
                    "rule": "int"
                },{
                    "name": "ins_detail",
                    "sharding_key": "entity_id",
                    "rule": "int"
                },{
                    "name": "order_ins_detail",
                    "sharding_key": "entity_id",
                    "rule": "int"
                }
            ]
        }
    ]
}'
```

### 演示zerodb库的分库分表

### 演示后发现配置没有问题，创建一个快照：schema_zerodb_added_ss
```
curl -X POST \
  http://10.1.24.173:5004/api/snapshot_config \
  -H 'Cache-Control: no-cache' \
  -H 'Postman-Token: ba477cf0-2e75-4b95-b4bb-d60ebb15213d' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F snapshotName=schema_zerodb_added_ss \
  -F clusterName=presentation_cluster
```

### 创建schema：instance 和附带的 table
```
curl -X POST \
  http://10.1.24.173:5004/api/add_schema \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 9efe0263-bf6b-4187-9af7-eaa2ab3575f0' \
  -d '{
    "clusterName": "presentation_cluster",
    "schemas": [
        {
            "Name": "instance",
            "nonsharding_host_group": "hostGroup1",
            "sharding_host_groups": [
                "hostGroup1","hostGroup2","hostGroup3","hostGroup4"
            ],
            "rw_split": false,
            "multi_route": true,
    		"init_conn_multi_route": true,
            "schema_sharding": 128,
            "table_sharding": 1,
            "table_configs": [
                {
                    "name": "order_ins",
                    "sharding_key": "entity_id",
                    "rule": "int"
                },{
                    "name": "ins_detail",
                    "sharding_key": "entity_id",
                    "rule": "int"
                },{
                    "name": "order_ins_detail",
                    "sharding_key": "entity_id",
                    "rule": "int"
                }
            ]
        }
    ]
}'
```

### 演示instance库的分库分表

### 演示后发现配置没有问题，创建一个快照：schema_instance_added_ss
```
curl -X POST \
  http://10.1.24.173:5004/api/snapshot_config \
  -H 'Cache-Control: no-cache' \
  -H 'Postman-Token: 23596756-0c78-46e6-aabf-3ad74ecdbc65' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F snapshotName=schema_instance_added_ss \
  -F clusterName=presentation_cluster
```

### 创建schema：custody 和附带的 table
```
curl -X POST \
  http://10.1.24.173:5004/api/add_schema \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: c149cf65-353f-4a56-9fb6-7346544d07d6' \
  -d '{
    "clusterName": "presentation_cluster",
    "schemas": [
        {
            "Name": "custody",
            "custody": true,
            "nonsharding_host_group": "hostGroupCommon",
            "rw_split": true
        }
    ]
}'
```

### 创建分库+分表 schema：many 和附带的 table
```
curl -X POST \
  http://10.1.24.173:5004/api/add_schema \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 17426685-10e9-4fc2-9579-1665ec6c11fa' \
  -d '{
    "clusterName": "presentation_cluster",
    "schemas": [
        {
            "Name": "many",
            "nonsharding_host_group": "hostGroup1",
            "sharding_host_groups": [
                "hostGroup1","hostGroup2","hostGroup3","hostGroup4"
            ],
            "rw_split": false,
            "multi_route": true,
    		"init_conn_multi_route": true,
            "schema_sharding": 128,
            "table_sharding": 8,
            "table_configs": [
                {
                    "name": "order_ins",
                    "sharding_key": "entity_id",
                    "rule": "int"
                },{
                    "name": "ins_detail",
                    "sharding_key": "entity_id",
                    "rule": "int"
                },{
                    "name": "order_ins_detail",
                    "sharding_key": "entity_id",
                    "rule": "int"
                }
            ]
        }
    ]
}'
```

### 想要回到snapshot为inital_config_ss的版本，并且让所有proxy都生效
```
mysql> show databases;
+-----------+
| DATABASES |
+-----------+
| many      |
| custody   |
| zerodb    |
| instance  |
+-----------+

curl -X PUT \
  'http://10.1.24.173:5004/api/push_config?clusterName=presentation_cluster&snapshotName=schema_zerodb_added_ss' \
  -H 'Cache-Control: no-cache' \
  -H 'Postman-Token: 848a42bc-f63b-4c82-a48e-4d5d8aa8f83d'

mysql> show databases;
+-----------+
| DATABASES |
+-----------+
| zerodb    |
+-----------+
```

### 删除 table：order_ins
```

curl -X POST \
  http://10.1.24.173:5004/api/delete_table \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Postman-Token: 19ef98ce-2605-45d6-bc4e-b5399e756c79' \
  -d 'clusterName=presentation_cluster&tableName=order_ins&schemaName=zerodb'

```

### 删除 schema: zerodb
```
curl -X POST \
  http://10.1.24.173:5004/api/delete_schema \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 543804b1-00ac-4ac9-826d-dd8b68012352' \
  -d '{
    "clusterName": "presentation_cluster",
    "names": [
        "zerodb"
    ]
}'

mysql> show databases;
ERROR 1105 (HY000): no data for show statement
```

### 关闭proxy，增加 schema: zerodb
```
中间有个失败了
```

###




