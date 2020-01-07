package frontend

import (
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"time"
)

func (fc *FrontendConn) handleWriteSQL(sql string, stmt sqlparser.Statement, args []interface{}) error {
	if len(fc.logicDb) == 0 {
		return errcode.BuildError(errcode.NoDBUsed)
	}
	sn := fc.proxy.GetSchemaNode(fc.logicDb)
	if sn == nil {
		return errcode.BuildError(errcode.DBNotExist, fc.logicDb)
	}

	startTime := time.Now().UnixNano()

	plan, err := fc.proxy.GetRouter().BuildPlan(fc.logicDb, fc.GetRemoteIP(), stmt, sql, fc.permitMultiRoute(fc.status), sn.MultiRoutePermitted, fc.handlerArena.AllocWithLen(0, 128))
	if err != nil {
		return err
	}

	buildPlanTime := time.Now().UnixNano()
	conns, err := fc.getShardingBackendConns(false, false, plan)
	needRollback := err != nil

	// 所有的连接一起归还给连接池
	defer fc.closeShardConns(conns, needRollback)

	if err != nil {
		return err
	}
	if conns == nil {
		return fc.writeOK(nil)
	}
	gotConnsTime := time.Now().UnixNano()

	var rs []*mysql.Result

	rs, err = fc.executeInMultiNodes(conns, plan.ShardingNodes, args, fc.proxy.LogSQL)

	monitor.Monitor.IncrClientTPS(fc.logicDb)

	executeTime := time.Now().UnixNano()
	if err == nil {
		err = fc.mergeExecResult(rs)
	}
	mergeTime := time.Now().UnixNano()

	logSlowSql(mergeTime, startTime, buildPlanTime, gotConnsTime, executeTime, fc.proxy.slowLogTime, plan.ShardingNodes, sql)

	return err
}
