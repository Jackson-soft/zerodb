package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
)

// 不进行路由
func (fc *FrontendConn) handleSimpleSelect(sql string, stmt *sqlparser.SimpleSelect) error {
	var rs []*mysql.Result
	var err error

	conn, err := fc.getInfoSingleBackendConn(sql)
	defer fc.closeConn(conn, false)
	if err != nil {
		return err
	}
	//execute.sql may be rewritten in getShowExecDB
	rs, err = fc.executeInNode(conn, sql, nil)
	if err != nil {
		return err
	}

	monitor.Monitor.IncrClientQPS("")

	if len(rs) == 0 {
		msg := fmt.Sprintf("result is empty")
		return mysql.NewError(mysql.ER_UNKNOWN_ERROR, msg)
	}

	fc.lastInsertId = int64(rs[0].InsertId)
	fc.affectedRows = int64(rs[0].AffectedRows)

	if rs[0].Resultset != nil {
		err = fc.writeResultset(fc.status, rs[0].Resultset)
	} else {
		err = fc.writeOK(rs[0])
	}

	if err != nil {
		return err
	}

	return nil
}

//build select result with group by opt
func (fc *FrontendConn) buildSelectGroupByResult(rs []*mysql.Result,
	stmt *sqlparser.Select) (*mysql.Result, error) {
	var err error
	var r *mysql.Result
	var groupByIndexs []int

	fieldLen := len(rs[0].Fields)
	startIndex := fieldLen - len(stmt.GroupBy)
	for startIndex < fieldLen {
		groupByIndexs = append(groupByIndexs, startIndex)
		startIndex++
	}

	funcExprs := fc.getFuncExprs(stmt)
	if len(funcExprs) == 0 {
		r, err = fc.mergeGroupByWithoutFunc(rs, groupByIndexs)
	} else {
		r, err = fc.mergeGroupByWithFunc(rs, groupByIndexs, funcExprs)
	}
	if err != nil {
		return nil, err
	}

	//build result
	names := make([]string, 0, 2)
	if 0 < len(r.Values) {
		r.Fields = r.Fields[:groupByIndexs[0]]
		for i := 0; i < len(r.Fields) && i < groupByIndexs[0]; i++ {
			names = append(names, string(r.Fields[i].Name))
		}
		//delete group by columns in Values
		for i := 0; i < len(r.Values); i++ {
			r.Values[i] = r.Values[i][:groupByIndexs[0]]
		}
		r.Resultset, err = fc.buildResultset(r.Fields, names, r.Values)
		if err != nil {
			return nil, err
		}
	} else {
		r.Resultset = fc.newEmptyResultset(stmt)
	}

	return r, nil
}
