```yaml
# 配置的版本
config_name: 测试版本

# server user and password
user :  zerodb
password : zerodb

# only log the query that take more than slow_log_time ms
slow_log_time : 1000

# 如果推送proxy配置失败了，是否需要继续推送后续的proxy
push_when_fail: false

stop_service:
    offline_on_lost_keeper: true
    offline_swh_rejected_num: 3
    offline_down_host_num: 4
    offline_recover: true

# 「数据源切换」配置，也有Host级别的
switch:
#    # 默认值为true
#    need_vote: true
#    # 只有当获得了 >= n% 的投票后才能数据源切换
#    vote_approve_ratio: 50
#
#    # 和target agent有关
#    # 默认值为true
#    need_load_check: true
#    safe_load: 2
#
#    # 与source agent和target agent都有关
#    # need_binlog_check: false
#    need_binlog_check: true
#    #    mysql> show master status;
#    #    +------------------+-----------+--------------+------------------+-------------------+
#    #    | File             | Position  | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
#    #    +------------------+-----------+--------------+------------------+-------------------+
#    #    | mysql-bin.000056 | 176224573 |              |                  |                   |
#    #    +------------------+-----------+--------------+------------------+-------------------+
#    safe_binlog_delay: 1000
    need_vote: true
    vote_approve_ratio: 50
    need_load_check: true
    safe_load: 8
    need_binlog_check: false
    safe_binlog_delay: 1000

host_groups:
-
    name : hostGroup1
    max_conn : 2560
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
    active_write: 0
-
    name : hostGroup2
    max_conn : 2560
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.2:3306,10.1.21.80:3306
    active_write: 0
-
    name : hostGroup3
    max_conn : 2560
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.3:3306,10.1.21.81:3306
    active_write: 0
-
    name : hostGroup4
    max_conn : 2560
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.4:3306,10.1.21.82:3306
    active_write: 0
-
    name : hostGroupCommon
    max_conn : 1024
    init_conn : 1
    user :  twodfire
    password : 123456
    active_write: 0
    write : common101.my.2dfire-daily.com:3306
    read : common101.my.2dfire-daily.com:3306@2,common101.my.2dfire-daily.com:3306@3

############################# 日常配置 #####################################
-
    name : hostGroup11
    max_conn : 2560
    init_conn : 10
    user : order
    password : order@552208
    write : 10.1.6.101:3306
    active_write: 0

-
    name : hostGroup12
    max_conn : 2560
    init_conn : 10
    user : order
    password : order@552208
    write : 10.1.6.102:3306
    active_write: 0
-
    name : hostGroup13
    # default max conns for mysql server
    max_conn : 2560
    init_conn : 10
    user : order
    password : order@552208
    write : 10.1.6.103:3306
    active_write: 0
-
    name : hostGroup14
    # default max conns for mysql server
    max_conn : 2560
    init_conn : 10
    user : order
    password : order@552208
    write : 10.1.6.104:3306
    active_write: 0

schema_configs :
-
    name: zerodb
    sharding_host_groups: [hostGroup1,hostGroup2,hostGroup3,hostGroup4]
    rw_split: false
    multi_route: true
    init_conn_multi_route: true
    ################## 分库分表配置范围 #########################
    # 1 <= table_sharding <= 1024 / (schema_sharding)
    # 1 <= table_sharding <= 8
    schema_sharding: 128
    table_sharding: 1
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: int
    -
        name: ins_detail
        sharding_key: entity_id
        rule: int
    -
        name: order_ins_detail
        sharding_key: entity_id
        rule: int
-
    name: instance
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2,hostGroup3,hostGroup4]
    rw_split: false
    multi_route: true
    init_conn_multi_route: true
    ################## 分库分表配置范围 #########################
    # 1 <= table_sharding <= 1024 / (schema_sharding)
    # 1 <= table_sharding <= 8
    schema_sharding: 128
    table_sharding: 1
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: int
    -
        name: ins_detail
        sharding_key: entity_id
        rule: int
    -
        name: order_ins_detail
        sharding_key: entity_id
        rule: int
-
    # 单独用于托管某个数据库
    name: custody
    rw_split: true
    custody: true
    multi_route: true
    init_conn_multi_route: true
    nonsharding_host_group: hostGroupCommon
-
    name: many
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2,hostGroup3,hostGroup4]
    rw_split: false
    multi_route: true
    init_conn_multi_route: true
    schema_sharding: 128
    table_sharding: 8
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: int
    -
        name: ins_detail
        sharding_key: entity_id
        rule: int
    -
        name: order_ins_detail
        sharding_key: entity_id
        rule: int


################################ 与Cobar保持一致的分库分表配置 ######################################
-
    name: order
    nonsharding_host_group: hostGroup11
    sharding_host_groups: [hostGroup11,hostGroup12,hostGroup13,hostGroup14]
    rw_split: false
    multi_route: true
    init_conn_multi_route: true
    schema_sharding: 128
    table_sharding: 1
    table_configs:
    -
        name: paydetail
        sharding_key: entity_id
        rule: string
    -
        name: waitinginstanceinfo
        sharding_key: entity_id
        rule: string
    -
        name: totalpayinfo
        sharding_key: entity_id
        rule: string
    -
        name: payinfo
        sharding_key: entity_id
        rule: string
    -
        name: servicebillinfo
        sharding_key: entity_id
        rule: string
    -
        name: specialfee
        sharding_key: entity_id
        rule: string
    -
        name: orderdetail
        sharding_key: entity_id
        rule: string
    -
        name: takeout_order_extra
        sharding_key: entity_id
        rule: string
    -
        name: takeout_instance_extra
        sharding_key: waitinginstance_id
        rule: string
    -
        name: globalcodeorder
        sharding_key: global_code
        rule: string
    -
        name: simplecodeorder
        sharding_key: simple_code
        rule: string
    -
        name: instancedetail
        sharding_key: entity_id
        rule: int
    -
        name: customer_order_relation
        sharding_key: customerregister_id
        rule: string
    -
        name: waitingorderdetail
        sharding_key: entity_id
        rule: int
    -
        name: waitingordercrid
        sharding_key: customerregister_id
        rule: string
    -
        name: waiting_pay
        sharding_key: entity_id
        rule: int
    -
        name: order_snapshot
        sharding_key: entity_id
        rule: int
    -
        name: order_bill
        sharding_key: entity_id
        rule: int
    -
        name: instancebill
        sharding_key: entity_id
        rule: int
    -
        name: discount_detail
        sharding_key: entity_id
        rule: int
    -
        name: presell_order_extra
        sharding_key: entity_id
        rule: int
    -
        name: order_refund
        sharding_key: entity_id
        rule: int
    -
        name: refund_pay_item
        sharding_key: entity_id
        rule: int
```