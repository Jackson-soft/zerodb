package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"sort"
	"strconv"
	"strings"
)

func (fc *FrontendConn) handleRoute(sql string, s *sqlparser.Route) error {
	if len(fc.logicDb) == 0 {
		return errcode.BuildError(errcode.NoDBUsed)
	}

	sql = strings.TrimLeft(strings.ToLower(sql), "route")

	var stmt sqlparser.Statement
	stmt, err := sqlparser.Parse(sql) //解析sql语句,得到的stmt是一个interface
	if err != nil {
		return err
	}

	shardingNodes, err := fc.buildRoutePlan(sql, stmt)
	if err != nil {
		return err
	}

	var rows [][]string
	var names []string
	names = []string{
		"schemaIndex",
		"tableIndex",
		"hostGroup",
		"shardingSQL",
	}
	var Column = len(names)
	var values [][]interface{}

	keys := make([]int, 0, len(shardingNodes))
	for key := range shardingNodes {
		keys = append(keys, key)
	}

	sort.Ints(keys)
	for _, k := range keys {
		sd := shardingNodes[k]
		var hg *backend.HostGroupNode

		result, ok := fc.proxy.HostGroupNodes.Load(sd.HostGroup)
		if ok {
			hg = result.(*backend.HostGroupNode)
		}

		rows = append(rows, []string{
			strconv.Itoa(sd.SchemaIndex),
			strconv.Itoa(sd.TableIndex),
			hg.Write[hg.GetActivedWriteIndex()].GetName(),
			sd.ShardingSQL,
		})
	}

	values = make([][]interface{}, len(rows))
	for i := range rows {
		values[i] = make([]interface{}, Column)
		for j := range rows[i] {
			values[i][j] = rows[i][j]
		}
	}

	if len(rows) == 0 {
		return fmt.Errorf("no data for route statement")
	}

	resultset, err := fc.buildResultset(nil, names, values)
	if err != nil {
		return err
	}
	return fc.writeResultset(fc.status, resultset)
}

func (fc *FrontendConn) buildRoutePlan(sql string, stmt sqlparser.Statement) (map[int]*backend.ShardingNode, error) {
	plan, err := fc.proxy.GetRouter().BuildPlan(fc.logicDb, fc.GetRemoteIP(), stmt, sql, true, true, fc.handlerArena.AllocWithLen(0, 256))
	if err != nil {
		return nil, err
	}
	return plan.ShardingNodes, nil
}
