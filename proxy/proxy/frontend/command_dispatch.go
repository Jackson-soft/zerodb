package frontend

import (
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
)

func (fc *FrontendConn) dispatch(data []byte) error {
	cmdType := data[0]
	data = data[1:]

	switch cmdType {
	case mysql.COM_QUIT:
		fc.handleRollback()
		fc.Close(true)
		return nil
	case mysql.COM_QUERY:
		return fc.handleQuery(util.String(data))
	case mysql.COM_PING:
		return fc.writeOK(nil)
	case mysql.COM_INIT_DB:
		return fc.handleUseDB(util.String(data))
	//case mysql.COM_FIELD_LIST: //show column
	//	return c.handleFieldList(data)
	case mysql.COM_STMT_PREPARE:
		return fc.handleStmtPrepare(util.String(data))
	case mysql.COM_STMT_EXECUTE:
		return fc.handleStmtExecute(data)
	case mysql.COM_STMT_CLOSE:
		return fc.handleStmtClose(data)
	case mysql.COM_STMT_SEND_LONG_DATA:
		return fc.handleStmtSendLongData(data)
	case mysql.COM_STMT_RESET:
		return fc.handleStmtReset(data)
	case mysql.COM_SET_OPTION:
		return fc.writeEOF(0)
	default:
		//msg := fmt.Sprintf("command dispatch failed. Unsupported command type : [%d]", cmdType)
		//return mysql.NewError(mysql.ER_UNKNOWN_ERROR, msg)
		return fc.writeEOF(0)
	}

	return nil
}
