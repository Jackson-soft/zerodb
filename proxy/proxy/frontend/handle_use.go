package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"github.com/pkg/errors"
)

// 如果分库情况: order0, order1, order2
// use schema 则会"物理上"切换到order0，"逻辑上"切换到order库
func (fc *FrontendConn) handleUseDB(dbName string) error {
	var co *backend.BackendConn
	var err error

	if len(dbName) == 0 {
		return fmt.Errorf("database name must not be empty")
	}

	r := fc.proxy.router
	schemaRouteRule := r.GetSchemaRule(dbName)

	if schemaRouteRule == nil {
		return errcode.BuildError(errcode.SchemaRuleNotExist, dbName)
	}

	// 默认选取第一个schemaIndex为0或者1的hostGroup来代表整个schema
	hostGroup := schemaRouteRule.SchemaToHostGroupNode[0]
	if hostGroup == "" {
		hostGroup = schemaRouteRule.SchemaToHostGroupNode[1]
	}

	defaultNode := backend.NewShardingNode(0, 0, backend.ShardingTypeSchema, "", hostGroup)

	var physicalDb string
	if schemaRouteRule.Type == TypeCustody {
		physicalDb = dbName
	} else {
		physicalDb = fmt.Sprintf("%s%d", dbName, 0)
	}

	//get the connection from slave preferentially
	co, err = fc.getBackendConn(defaultNode, false, true)
	defer fc.closeConn(co, false)
	if err != nil {
		return err
	}

	fc.logicDb = dbName
	fc.physicalDb = physicalDb

	if err = co.UseDB(physicalDb); err != nil {
		return errors.WithMessage(err, "switch to "+fc.logicDb+", but")
	}
	return fc.writeOK(nil)
}
