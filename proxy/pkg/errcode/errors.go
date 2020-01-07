package errcode

import (
	"bytes"
	"errors"
	"fmt"
)

type ErrorTuple struct {
	ErrMsg  string
	ErrCode string
}

var (
	ErrNoDBPool = errors.New("no db pool")

	ErrDBPoolClosed = errors.New("db pool is closed")
	ErrConnIsNil    = errors.New("connection is nil")
	ErrBadConn      = errors.New("connection was bad")
	ErrSQLIgnored   = errors.New("this sql is ignored")

	ErrAddressNull     = errors.New("address is nil")
	ErrInvalidArgument = errors.New("argument is invalid")

	ErrCalcRoute               = errors.New("calc route failed")
	ErrNoShardingNode          = errors.New("no sharding node")
	ErrSumColumnType           = errors.New("sum column type error")
	ErrSelectInInsert          = errors.New("select in insert not allowed")
	ErrInsertInMulti           = errors.New("insert in multi node")
	ErrUpdateInMulti           = errors.New("update in multi node")
	ErrDeleteInMulti           = errors.New("delete in multi node")
	ErrReplaceInMulti          = errors.New("replace in multi node")
	ErrExecInMulti             = errors.New("exec in multi node")
	ErrMultiNodeTranNotSupport = errors.New("transaction in multi node is not allowed")

	ErrSQLNotSupport    = errors.New("statement is not supported")
	ErrNoPlanRule       = errors.New("statement have no plan rule")
	ErrUpdateKey        = errors.New("routing key in update expression")
	ErrStmtConvert      = errors.New("statement fail to convert")
	ErrKeyOutOfRange    = errors.New("shard key not in key range")
	ErrMultiShard       = errors.New("insert or replace has multiple shard targets")
	ErrIRNoColumns      = errors.New("insert or replace must specify columns")
	ErrIRNoShardingKey  = errors.New("insert or replace not contain sharding key")
	ErrColsLenNotMatch  = errors.New("insert or replace cols and values length not match")
	ErrDateIllegal      = errors.New("date format illegal")
	ErrDateRangeIllegal = errors.New("date range format illegal")
	ErrDateRangeCount   = errors.New("date range count is not equal")
	ErrReadExist        = errors.New("read has exist")
	ErrSlaveNotExist    = errors.New("slave has not exist")
	ErrBlackSqlExist    = errors.New("black sql has exist")
	ErrBlackSqlNotExist = errors.New("black sql has not exist")
	ErrInsertTooComplex = errors.New("insert is too complex")
	ErrSQLNULL          = errors.New("sql is null")
	ErrSQLNotSupported  = errors.New("sql is not supported")
	ErrProxyNumLittle   = errors.New("proxy count two little ")
	ObjectEmptyErr      = ErrorTuple{ErrMsg: "[%s] should not be empty", ErrCode: "1000"}

	// HostGroup
	HostNotWritable           = ErrorTuple{ErrMsg: "host group [%s] is not writable", ErrCode: "1000"}
	HostGroupNotExist         = ErrorTuple{ErrMsg: "host group [%s] doesn't exist", ErrCode: "1000"}
	ErrMultiRouteNotPermitted = ErrorTuple{ErrMsg: "sharding key not provided, schema's multiRoutePermitted: %v, connection's multiRoutePermitted: %v", ErrCode: "1000"}

	ErrNoReadDBPool         = errors.New("there is no read db pool")
	ErrReadConfNotExist     = errors.New("read db pool config")
	ErrNoWriteDBPool        = errors.New("there is no write db pool")
	ShardingKeyTypeInvalid  = ErrorTuple{ErrMsg: "sharding key value [%v] is not supported", ErrCode: "1000"}
	NoDBUsed                = ErrorTuple{ErrMsg: "database is not used", ErrCode: "1000"}
	DBNotExist              = ErrorTuple{ErrMsg: "database [%s] doesn't exist", ErrCode: "1000"}
	ErrNoConn               = errors.New("there is no conn")
	NonshardingConfNotExist = ErrorTuple{ErrMsg: "there is no nonsharding config for schema [%s], please check table name [%s]", ErrCode: "1000"}
	SchemaRuleNotExist      = ErrorTuple{ErrMsg: "the rule of schema [%s] doesn't exist, please check your schema config", ErrCode: "1000"}
	InvalidCharset          = ErrorTuple{ErrMsg: "charset [%s] is invalid", ErrCode: "1000"}
	NoShardingKeyExplain    = ErrorTuple{ErrMsg: "no sharding key in explain", ErrCode: "1000"}
)

func BuildError(err ErrorTuple, args ...interface{}) error {
	var buf bytes.Buffer
	buf.WriteString("code: ")
	buf.WriteString(err.ErrCode)
	buf.WriteString(", ")
	buf.WriteString(err.ErrMsg)
	return fmt.Errorf(buf.String(), args...)
}
