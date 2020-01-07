package frontend

import (
	"testing"

	"git.2dfire.net/zerodb/common/config"
	"gopkg.in/yaml.v3"
)

func TestParseHostGroups(t *testing.T) {
	var s1 = `
host_groups:
-
    name : hostGroup1
    max_conn : 1023
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	var cfg config.Config
	if err := yaml.Unmarshal([]byte(s1), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	testServer, err := NewProxyEngine("127.0.0.1:9696", "utf8mb4", false)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, max_conn is less than 1024")
	}

	var s2 = `
host_groups:
-
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s2), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, name is empty")
	}

	var s3 = `
host_groups:
-
    name : hostGroup1
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s3), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, max_conn is less than 1024")
	}

	var s4 = `
host_groups:
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 1025
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s4), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, init_conn is great than max_conn")
	}

	var s5 = `
host_groups:
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 10
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s5), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, user is empty")
	}

	var s7 = `
host_groups:
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s7), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, write is empty")
	}

	var s8 = `
host_groups:
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306
    read : 10.1.22.1:3306@a,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s8), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, invalid weight of read")
	}

	var s9 = `
host_groups:
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
#    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
`
	if err := yaml.Unmarshal([]byte(s9), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseHostGroupNodes(&cfg, false)
	if err == nil {
		t.Errorf("Should have error here, dup hostGroup")
	}
}

func TestParseSchemas(t *testing.T) {
	testServer, err := NewProxyEngine("127.0.0.1:9696", "utf8mb4", false)
	if err != nil {
		t.Errorf(err.Error())
	}
	var h = `
host_groups:
-
    name : hostGroup1
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
    read : 10.1.22.1:3306@3,10.1.22.1:3306@2
-
    name : hostGroup2
    max_conn : 1024
    init_conn : 10
    user :  zerodb
    password : zerodb@552208
    write : 10.1.22.1:3306,10.1.21.79:3306
    read : 10.1.22.1:3306@3,10.1.22.1:3306@2

`
	var cfg config.Config
	if err := yaml.Unmarshal([]byte(h), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	if err := testServer.ParseHostGroupNodes(&cfg, false); err != nil {
		t.Errorf(err.Error())
	}
	var s1 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 252
    table_sharding: 2
`
	if err := yaml.Unmarshal([]byte(s1), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, schema_sharding is not power of 2")
	}

	var s2 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
    table_sharding: 23
`
	if err := yaml.Unmarshal([]byte(s2), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, table_sharding is not power of 2")
	}

	var s3 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
    table_sharding: 16
`
	if err := yaml.Unmarshal([]byte(s3), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, product is great than 1024")
	}

	var s4 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 1
    table_sharding: 2
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: string
    -
        name: ins_detail
        sharding_key: entity_id
        rule: string
    -
        name: order_sth
        sharding_key: order_id
        rule: int
`
	if err := yaml.Unmarshal([]byte(s4), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, schema_sharding is less than length of sharding_host_groups")
	}

	var s5 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    schema_sharding: 256
    table_sharding: 2
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: string
`
	if err := yaml.Unmarshal([]byte(s5), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, sharding_host_groups empty")
	}

	var s6 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
    table_sharding: 2
    table_configs:
    -
        sharding_key: entity_id
        rule: string

`
	if err := yaml.Unmarshal([]byte(s6), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, table name empty")
	}

	var s7 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
    table_sharding: 2
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
`
	if err = yaml.Unmarshal([]byte(s7), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, rule empty")
	}

	var s8 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
    table_sharding: 2
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: string11
`
	if err := yaml.Unmarshal([]byte(s8), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, rule invalid")
	}

	var s9 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
    table_sharding: 0
    table_configs:
    -
        name: order_ins
        sharding_key: entity_id
        rule: string
    -
        name: order_sth
        sharding_key: order_id
        rule: int
`
	if err := yaml.Unmarshal([]byte(s9), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, table_sharding is 0")
	}

	var s10 = `
schema_configs :
-
    name: hello
    custody: true
`
	if err := yaml.Unmarshal([]byte(s10), &cfg); err != nil {
		t.Errorf(err.Error())
	}
	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, custody schema doesn't have an NonshardingHostGroup")
	}

	var s11 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    table_sharding: 2
`
	if err := yaml.Unmarshal([]byte(s11), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, schema_sharding empty")
	}

	var s12 = `
schema_configs :
-
    name: zerodb
    nonsharding_host_group: hostGroup1
    sharding_host_groups: [hostGroup1,hostGroup2]
    schema_sharding: 256
`
	if err := yaml.Unmarshal([]byte(s12), &cfg); err != nil {
		t.Errorf(err.Error())
	}

	err = testServer.ParseSchemas(&cfg)
	if err == nil {
		t.Errorf("Should have error here, table_sharding empty")
	}

}
