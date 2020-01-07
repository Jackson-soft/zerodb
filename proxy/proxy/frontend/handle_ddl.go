package frontend

import (
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
)

func (fc *FrontendConn) handleDdl(stmt sqlparser.Statement, sql string) error {
	if len(fc.logicDb) == 0 {
		return errcode.BuildError(errcode.NoDBUsed)
	}
	sn := fc.proxy.GetSchemaNode(fc.logicDb)
	if sn == nil {
		return errcode.BuildError(errcode.DBNotExist, fc.logicDb)
	}

	plan, err := fc.proxy.GetRouter().BuildPlan(fc.logicDb, fc.GetRemoteIP(), stmt, sql, fc.permitMultiRoute(fc.status), sn.MultiRoutePermitted, fc.handlerArena.AllocWithLen(0, 128))
	if err != nil {
		return err
	}
	conns, err := fc.getShardingBackendConns(false, false, plan)
	defer fc.closeShardConns(conns, err != nil)
	if err != nil {
		return err
	}
	if conns == nil {
		return fc.writeOK(nil)
	}

	var rs []*mysql.Result

	rs, err = fc.executeInMultiNodes(conns, plan.ShardingNodes, nil, fc.proxy.LogSQL)
	if err == nil {
		err = fc.mergeExecResult(rs)
	}

	return err
}
