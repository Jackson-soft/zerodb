package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"strings"
)

var nstring = sqlparser.String

func (fc *FrontendConn) handleSet(stmt *sqlparser.Set, sql string) (err error) {
	if len(stmt.Exprs) != 1 && len(stmt.Exprs) != 2 {
		return fmt.Errorf("must set one item once, not %s", nstring(stmt))
	}

	k := string(stmt.Exprs[0].Name.Name)
	switch strings.ToUpper(k) {
	case `AUTOCOMMIT`, `@@AUTOCOMMIT`, `@@SESSION.AUTOCOMMIT`:
		return fc.handleSetAutoCommit(stmt.Exprs[0].Expr)
	case `MULTI_ROUTE_PERMIT`:
		return fc.handleMultiRoutePermit(stmt.Exprs[0].Expr)
	case `NAMES`,
		`CHARACTER_SET_RESULTS`, `@@CHARACTER_SET_RESULTS`, `@@SESSION.CHARACTER_SET_RESULTS`,
		`CHARACTER_SET_CLIENT`, `@@CHARACTER_SET_CLIENT`, `@@SESSION.CHARACTER_SET_CLIENT`,
		`CHARACTER_SET_CONNECTION`, `@@CHARACTER_SET_CONNECTION`, `@@SESSION.CHARACTER_SET_CONNECTION`:
		if len(stmt.Exprs) == 2 {
			//SET NAMES 'charset_name' COLLATE 'collation_name'
			return fc.handleSetNames(stmt.Exprs[0].Expr, stmt.Exprs[1].Expr)
		}
		return fc.handleSetNames(stmt.Exprs[0].Expr, nil)
	default:
		glog.Glog.Warnf("connectionID:[%v], command is not executed. sql:%s", fc.connectionId, sql)
		return fc.writeOK(nil)
	}
}

func (fc *FrontendConn) handleSetAutoCommit(val sqlparser.ValExpr) error {
	flag := sqlparser.String(val)
	flag = strings.Trim(flag, "'`\"")
	// autocommit允许为 0, 1, ON, OFF, "ON", "OFF", 不允许"0", "1"
	if flag == `0` || flag == `1` {
		_, ok := val.(sqlparser.NumVal)
		if !ok {
			return fmt.Errorf("set autocommit error")
		}
	}
	switch strings.ToUpper(flag) {
	case `1`, `ON`:
		// 如果设置了auto commit，需要把txConns(事务session)中的所有对应关系删除
		fc.status |= mysql.SERVER_STATUS_AUTOCOMMIT
		if fc.status&mysql.SERVER_STATUS_IN_TRANS > 0 {
			fc.status &= ^mysql.SERVER_STATUS_IN_TRANS
		}
		var pool *backend.DBPool
		for _, co := range fc.txConns {
			pool = co.GetDBPool()
			if e := co.SetAutoCommit(1); e != nil {
				pool.ReturnConn(co)
				fc.txConns = make(map[int]*backend.BackendConn)
				return fmt.Errorf("set autocommit error, %v", e)
			}
			pool.ReturnConn(co)
		}
		fc.txConns = make(map[int]*backend.BackendConn)
	case `0`, `OFF`:
		fc.status &= ^mysql.SERVER_STATUS_AUTOCOMMIT
	default:
		return fmt.Errorf("invalid autocommit flag %s", flag)
	}

	return fc.writeOK(nil)
}

func (fc *FrontendConn) handleMultiRoutePermit(val sqlparser.ValExpr) error {
	flag := sqlparser.String(val)
	flag = strings.Trim(flag, "'`\"")
	// MULTI_ROUTE_PERMIT 允许为 0, 1, ON, OFF, "ON", "OFF", 不允许"0", "1"
	if flag == `0` || flag == `1` {
		_, ok := val.(sqlparser.NumVal)
		if !ok {
			return fmt.Errorf("set MULTI_ROUTE_PERMIT error")
		}
	}
	switch strings.ToUpper(flag) {
	case `1`, `ON`:
		fc.status |= mysql.SERVER_STATUS_PERMIT_MULTI_ROUTE
	case `0`, `OFF`:
		fc.status &= ^mysql.SERVER_STATUS_PERMIT_MULTI_ROUTE
	default:
		return fmt.Errorf("invalid MULTI_ROUTE_PERMIT flag %s", flag)
	}

	return fc.writeOK(nil)
}

func (fc *FrontendConn) handleSetNames(ch, ci sqlparser.ValExpr) error {
	var cid mysql.CollationId
	var ok bool

	value := sqlparser.String(ch)
	value = strings.Trim(value, "'`\"")

	charset := strings.ToLower(value)
	if charset == "null" {
		return fc.writeOK(nil)
	}
	if ci == nil {
		if charset == "default" {
			charset = mysql.DEFAULT_CHARSET
		}
		cid, ok = mysql.CharsetIds[charset]
		if !ok {
			return fmt.Errorf("invalid charset %s", charset)
		}
	} else {
		collate := sqlparser.String(ci)
		collate = strings.Trim(collate, "'`\"")
		collate = strings.ToLower(collate)
		cid, ok = mysql.CollationNames[collate]
		if !ok {
			return fmt.Errorf("invalid collation %s", collate)
		}
	}
	fc.charset = charset
	fc.collation = cid

	return fc.writeOK(nil)
}
