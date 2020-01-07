package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
)

func (fc *FrontendConn) isInTransaction() bool {
	return fc.status&mysql.SERVER_STATUS_IN_TRANS > 0 ||
		!fc.isAutoCommit()
}

func (fc *FrontendConn) isAutoCommit() bool {
	return fc.status&mysql.SERVER_STATUS_AUTOCOMMIT > 0
}

func (fc *FrontendConn) handleBegin() error {
	for _, co := range fc.txConns {
		if err := co.Begin(); err != nil {
			return err
		}
	}
	fc.status |= mysql.SERVER_STATUS_IN_TRANS
	return fc.writeOK(nil)
}

func (fc *FrontendConn) handleCommit() (err error) {
	if err := fc.commit(); err != nil {
		return err
	} else {
		return fc.writeOK(nil)
	}
}

func (fc *FrontendConn) handleRollback() (err error) {
	if err := fc.rollback(); err != nil {
		return err
	} else {
		return fc.writeOK(nil)
	}
}

// 涉及到一致性问题
func (fc *FrontendConn) commit() (err error) {
	fc.status &= ^mysql.SERVER_STATUS_IN_TRANS

	var pool *backend.DBPool
	for _, co := range fc.txConns {
		pool = co.GetDBPool()
		if e := co.Commit(); e != nil {
			err = e
		}
		pool.ReturnConn(co)
	}

	fc.txConns = make(map[int]*backend.BackendConn)
	return
}

// 涉及到一致性问题
func (fc *FrontendConn) rollback() (err error) {
	fc.status &= ^mysql.SERVER_STATUS_IN_TRANS

	var pool *backend.DBPool
	for _, co := range fc.txConns {
		pool = co.GetDBPool()
		if e := co.Rollback(); e != nil {
			err = e
		}
		pool.ReturnConn(co)
	}

	fc.txConns = make(map[int]*backend.BackendConn)
	return
}
