package frontend

import (
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"time"
)

func logSlowSql(mergeTime int64, startTime int64, buildPlanTime int64, gotConnsTime int64, executeTime int64, slowLogTime int, shardingNodes map[int]*backend.ShardingNode, sql string) {
	if ((mergeTime - startTime) / int64(time.Millisecond)) > int64(slowLogTime) {
		monitor.Monitor.IncrSlowLogTotal()

		if len(shardingNodes) == 1 {
			glog.Glog.Warnf("Slow SQL. shardingCount: %d, buildPlan: %v, getShardingBackendConns: %v, executeInMultiNodes: %v, mergeExecResult: %v. shardingNode:%v, SQL: [%s]",
				len(shardingNodes),
				(buildPlanTime-startTime)/int64(time.Millisecond),
				(gotConnsTime-buildPlanTime)/int64(time.Millisecond),
				(executeTime-gotConnsTime)/int64(time.Millisecond),
				(mergeTime-executeTime)/int64(time.Millisecond),
				shardingNodes,
				sql)
		} else {
			glog.Glog.Warnf("Slow SQL. shardingCount: %d, buildPlan: %v, getShardingBackendConns: %v, executeInMultiNodes: %v, mergeExecResult: %v. SQL: [%s]",
				len(shardingNodes),
				(buildPlanTime-startTime)/int64(time.Millisecond),
				(gotConnsTime-buildPlanTime)/int64(time.Millisecond),
				(executeTime-gotConnsTime)/int64(time.Millisecond),
				(mergeTime-executeTime)/int64(time.Millisecond),
				sql)
		}
	}
}
