package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"testing"
	"time"
)

func TestLogSlowSqlSingleNode(*testing.T) {
	startTime := 0
	buildTime := 1
	gotConnsTime := 2
	executeTime := 3
	mergeTime := 10000 * time.Millisecond
	slowLogTime := 1
	sql := "select * from order_ins where entity_id = '33'"
	shardingNodes := make(map[int]*backend.ShardingNode)
	shardingNodes[1] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from zerodb153.order_ins304 where entity_id = '33'",
		HostGroup:    "hostGroup3",
	}

	logSlowSql(int64(mergeTime), int64(startTime), int64(buildTime), int64(gotConnsTime), int64(executeTime), slowLogTime, shardingNodes, sql)
}

func TestLogSlowSqlMultiNodes(*testing.T) {
	startTime := 0
	buildTime := 1
	gotConnsTime := 2
	executeTime := 3
	mergeTime := 10000 * time.Millisecond
	slowLogTime := 1
	sql := "select * from order_ins where entity_id = '33'"
	shardingNodes := make(map[int]*backend.ShardingNode)
	shardingNodes[1] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from zerodb153.order_ins304 where entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	shardingNodes[2] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from zerodb153.order_ins304 where entity_id = '33'",
		HostGroup:    "hostGroup3",
	}

	logSlowSql(int64(mergeTime), int64(startTime), int64(buildTime), int64(gotConnsTime), int64(executeTime), slowLogTime, shardingNodes, sql)
}
