package frontend

import (
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"strings"
)

func (fc *FrontendConn) handleExplain(sql string, s *sqlparser.Explain) error {
	if len(fc.logicDb) == 0 {
		return errcode.BuildError(errcode.NoDBUsed)
	}

	sql = strings.TrimPrefix(strings.ToLower(sql), "explain ")

	var stmt sqlparser.Statement
	stmt, err := sqlparser.Parse(sql) //解析sql语句,得到的stmt是一个interface
	if err != nil {
		return err
	}

	plan, err := fc.buildExplainPlan(sql, stmt)
	if err != nil {
		return err
	}

	conns, err := fc.getShardingBackendConns(false, true, plan)
	if err != nil {
		return err
	}

	var rs []*mysql.Result
	rs, err = fc.executeInMultiNodes(conns, plan.ShardingNodes, nil, fc.proxy.LogSQL)
	monitor.Monitor.IncrClientQPS(fc.logicDb)

	defer fc.closeShardConns(conns, false)

	if err != nil {
		return err
	}

	r, err := fc.mergeResults(rs, nil)

	if err != nil {
		return err
	}

	return fc.writeResultset(r.Status, r.Resultset)
}

func (fc *FrontendConn) buildExplainPlan(sql string, stmt sqlparser.Statement) (*Plan, error) {
	plan, err := fc.proxy.GetRouter().BuildPlan(fc.logicDb, fc.GetRemoteIP(), stmt, sql, true, true, fc.handlerArena.AllocWithLen(0, 256))
	if err != nil {
		return nil, err
	}

	if plan != nil && !plan.hasShardingKey {
		return nil, errcode.BuildError(errcode.NoShardingKeyExplain)
	}

	// rewrite shardingNodes
	for _, value := range plan.ShardingNodes {
		value.ShardingSQL = "explain " + value.ShardingSQL
	}

	return plan, nil
}
