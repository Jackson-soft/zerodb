package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"runtime"
	"runtime/debug"
	"strings"
)

/*处理query语句*/
func (fc *FrontendConn) handleQuery(sql string) (err error) {
	defer func() {
		// 停止panic往caller传播
		if e := recover(); e != nil {
			glog.Glog.Errorf("error:%s, sql:%s", debug.Stack(), sql)
			if err, ok := e.(error); ok {
				const size = 4096
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				glog.Glog.Errorf("handleQuery error:%v, stack:%s, sql:%s", err, string(buf), sql)
			}
			return
		}
	}()

	sql = strings.TrimRight(sql, ";") //删除sql语句最后的分号

	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(sql) //解析sql语句,得到的stmt是一个interface
	if err != nil {
		glog.Glog.Warnf("Parse SQL failed, try to preHandle. err: %v, sql: [%s]", err, sql)
		// 很多语句是 类SQL 或者 命令 。这些语句是sqlparser不支持的，这些语句让默认的节点来处理
		hasHandled, err := fc.preHandleShard(sql)
		if err != nil {
			glog.Glog.Warnf("connectionID:[%v]. SQL: %s, preHandleShard failed, err: %s", fc.connectionId, sql, err.Error())
			return err
		}
		if hasHandled {
			return nil
		}

		return err
	}

	switch real := stmt.(type) {
	// ################# Read ###################
	case *sqlparser.Select:
		err = fc.handleReadSQL(sql, real, nil)
	case *sqlparser.SimpleSelect:
		// TODO nanxing need more test cases
		err = fc.handleSimpleSelect(sql, real)
	// ################# Read ###################

	// ################# Write 可以跨库写，但是没有分布式事务机制，只有弱事务 ###################
	// 不支持多值插入
	case *sqlparser.Insert:
		err = fc.handleWriteSQL(sql, stmt, nil)
	case *sqlparser.Update:
		err = fc.handleWriteSQL(sql, stmt, nil)
	case *sqlparser.Delete:
		err = fc.handleWriteSQL(sql, stmt, nil)
	case *sqlparser.Replace:
		err = fc.handleWriteSQL(sql, stmt, nil)
	case *sqlparser.Truncate:
		err = fc.handleWriteSQL(sql, stmt, nil)
	// ################# Write ###################
	case *sqlparser.DDL:
		err = fc.handleDdl(stmt, sql)
	case *sqlparser.Set:
		err = fc.handleSet(real, sql)
	case *sqlparser.Begin:
		err = fc.handleBegin()
	case *sqlparser.Commit:
		err = fc.handleCommit()
	case *sqlparser.Rollback:
		err = fc.handleRollback()
	case *sqlparser.Show:
		err = fc.handleShow(real)
	case *sqlparser.UseDB:
		err = fc.handleUseDB(real.DB)
	// ###########################################
	case *sqlparser.Route:
		err = fc.handleRoute(sql, real)
	case *sqlparser.Explain:
		err = fc.handleExplain(sql, real)
	default:
		err = fmt.Errorf("statement %T not support now", stmt)
	}

	if err != nil {
		glog.Glog.Errorf("connectionID:[%v]. handleQuery failed. SQL:%s, Err:%v", fc.connectionId, sql, err)
	}

	return err
}

func (fc *FrontendConn) closeConn(conn *backend.BackendConn, rollback bool) {
	if fc.isInTransaction() {
		return
	}

	if rollback {
		conn.Rollback()
	}

	conn.GetDBPool().ReturnConn(conn)
}

func (fc *FrontendConn) closeShardConns(conns map[int]*backend.BackendConn, rollback bool) {
	if fc.isInTransaction() {
		return
	}

	var pool *backend.DBPool
	for _, co := range conns {
		pool = co.GetDBPool()
		if rollback {
			co.Rollback()
		}
		pool.ReturnConn(co)
	}
}
