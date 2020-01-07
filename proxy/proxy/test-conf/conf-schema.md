basic:
    #必填
    config_name: 测试版本
    #必填
    user :  zerodb
    #必填，*号
    password : zerodb
    #必填，> 0
    slow_log_time : 1000
    #必填
    push_when_fail: false

stop_service:
    #必填
    offline_on_lost_keeper: true
    #必填, > 0
    offline_swh_rejected_num: 3
    #必填, >= 0
    offline_down_host_num: 4
    #必填
    offline_recover: true

switch:
    #必填
    need_vote: true
    #必填, 0 < vote_approve_ratio <= 100
    vote_approve_ratio: 50
    #必填
    need_load_check: true
    #必填, > 0
    safe_load: 2

    #必填
    need_binlog_check: true
    #必填, > 0
    safe_binlog_delay: 1000
    #必填, > 0
    binlog_wait_time: 5 # sec

    #必填, > 0
    frequency: 10 # sec
    #必填, >= 0
    backend_ping_interval: 5 # sec

- HostGroup:
	#必填
    name : hostGroup1
    #必填，max_conn >= 1024
    max_conn : 2560
    #必填，0 < init_conn <= max_conn
    init_conn : 10
    #必填
    user : order
    #必填，以*号显示
    password : order@552208
    #必填，ip:port
    write : 10.1.6.101:3306,10.1.6.101:3306
	read : 10.1.22.1:3306@3,10.1.22.1:3306@2
	#必填，0 <= active_write <= len(write) - 1
    active_write: 0
    #default true
    enable_switch: true


- host_group_clusters :
    #必填
    name: cluster_order
    nonsharding_host_group: hostGroup1
    #不重复，len(sharding_host_groups) = 2 ^ n，有顺序
    sharding_host_groups: [hostGroup1,hostGroup2,hostGroup3,hostGroup4]

```
func IsPowerOfTwo(x int) bool {
    return (x > 0) && ((x & (x - 1)) == 0)
}
```

- schema :
    #必填
    name: boss
    #必选
    host_group_cluster: cluster_boss
    #必填
    rw_split: false
    #必填
    multi_route: true
    #必填
    init_conn_multi_route: true
    #必填，2 ^ n
    schema_sharding: 128
    #必填，2 ^ n
    table_sharding: 1

 -
    #必填
    name: ir_menu_ext
    #必填
    sharding_key: entity_id
    #必填，int || string
    rule: int