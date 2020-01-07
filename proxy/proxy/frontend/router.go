package frontend

import (
	"fmt"
	"strings"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"strconv"
	"sync"
)

var (
	HashRuleIntType    = "int"
	HashRuleStringType = "string"
)

const (
	NonshardingTableIndex  = -1
	NonshardingSchemaIndex = 0

	TypeSharding = 1
	TypeCustody  = 2
)

type SchemaRouteRule struct {
	TableRouteRules sync.Map //map[string]*TableRouteRule

	DB string

	HasNonshardingHostGroup bool
	NonShardingIndex        int
	//一个库中只有一张分表
	OneTablePerSchema bool

	// 单独分库的
	ShardingSchemaCount int
	ShardingTableCount  int

	FullSchemaIndexes []int //分库角标，类似于 order002
	FullTableIndexes  []int //分表角标，类似于 instance012

	// 类型：分库，托管
	Type int

	// 索引
	TableToHostGroupNode  map[int]string //key is table index, and value is host group node name
	SchemaToHostGroupNode map[int]string
	TableToSchemaIndex    map[int]int

	HostGroupNodes []string
	Shard          map[string]Shard
}

type TableRouteRule struct {
	Table string
	Key   string
	// 类型
	Rule string
}

type Router struct {
	SchemaRules sync.Map //map[string]*SchemaRouteRule
}

/*func (r *SchemaRouteRule) FindNodeIndex(key interface{}) (int, error) {
    tableIndex, err := r.Shard.FindForKey(key)
    if err != nil {
        return -1, err
    }
    return r.TableToHostGroupNode[tableIndex], nil
}*/

func (srr *SchemaRouteRule) FindTableIndex(rule string, key interface{}) (int, error) {
	return srr.Shard[rule].FindForKey(key)
}

//UpdateExprs is the expression after set
func checkUpdateExprs(plan *Plan, exprs sqlparser.UpdateExprs) error {
	// TODO nanxing
	//if len(r.Rule) == 0 {
	//    return nil
	//} else if len(r.HostGroupNodes) == 1 {
	//    return nil
	//}

	for _, e := range exprs {
		if string(e.Name.Name) == plan.TableRouteRule.Key {
			return errcode.ErrUpdateKey
		}
	}
	return nil
}
func (srr *SchemaRouteRule) AddTableRouteRule(tableConfig config.TableConfig) error {
	if srr.GetTableRule(tableConfig.Name) != nil {
		return fmt.Errorf("duplicate table route:[%s]", tableConfig.Name)
	}

	tableRule := new(TableRouteRule)
	tableRule.Table = tableConfig.Name
	tableRule.Key = tableConfig.ShardingKey
	if tableConfig.Rule != HashRuleIntType && tableConfig.Rule != HashRuleStringType {
		return fmt.Errorf("invalid table rule:%s", tableConfig.Rule)
	}
	tableRule.Rule = tableConfig.Rule

	srr.TableRouteRules.Store(tableConfig.Name, tableRule)
	return nil
}
func (srr *SchemaRouteRule) GetTableRule(table string) *TableRouteRule {
	result, ok := srr.TableRouteRules.Load(table)
	if ok {
		return result.(*TableRouteRule)
	} else {
		return nil
	}
}

//NewRouter build router according to the config file
func NewRouter() (*Router, error) {
	rt := new(Router)

	return rt, nil
}

// 返回nil也是正常的
func (r *Router) GetSchemaRule(db string) *SchemaRouteRule {
	result, ok := r.SchemaRules.Load(db)
	if ok {
		return result.(*SchemaRouteRule)
	} else {
		return nil
	}
}

func (r *Router) LoadSchemaRule(schemaConfig *config.SchemaConfig, hostGroupClusterConf *config.HostGroupCluster) (*SchemaRouteRule, error) {
	if schemaConfig.Custody {
		return r.LoadCustodySchemaRule(schemaConfig, hostGroupClusterConf)
	} else {
		return r.LoadShardingSchemaRule(schemaConfig, hostGroupClusterConf)
	}
}

func (r *Router) LoadCustodySchemaRule(schemaConfig *config.SchemaConfig, hostGroupClusterConf *config.HostGroupCluster) (*SchemaRouteRule, error) {
	if len(hostGroupClusterConf.NonshardingHostGroup) == 0 {
		return nil, fmt.Errorf("custody schema [%s] must have an NonshardingHostGroup", schemaConfig.Name)
	}
	if len(hostGroupClusterConf.ShardingHostGroups) != 0 {
		return nil, fmt.Errorf("custody schema [%s] should not have ShardingHostGroup", schemaConfig.Name)
	}

	srr := new(SchemaRouteRule)
	srr.DB = schemaConfig.Name
	srr.TableToHostGroupNode = make(map[int]string)
	srr.TableToSchemaIndex = make(map[int]int)
	srr.SchemaToHostGroupNode = make(map[int]string)
	srr.Type = TypeCustody

	srr.HasNonshardingHostGroup = len(hostGroupClusterConf.NonshardingHostGroup) != 0

	srr.HostGroupNodes = make([]string, 1)
	srr.SchemaToHostGroupNode[NonshardingSchemaIndex] = hostGroupClusterConf.NonshardingHostGroup
	if srr.HasNonshardingHostGroup {
		// NonshardingIndex只是个标记
		srr.NonShardingIndex = NonshardingTableIndex
		srr.HostGroupNodes = append(srr.HostGroupNodes, hostGroupClusterConf.NonshardingHostGroup)
	}

	if srr.HasNonshardingHostGroup {
		srr.TableToHostGroupNode[NonshardingTableIndex] = hostGroupClusterConf.NonshardingHostGroup
		srr.TableToSchemaIndex[NonshardingTableIndex] = NonshardingSchemaIndex
		srr.FullSchemaIndexes = append(srr.FullSchemaIndexes, NonshardingSchemaIndex)
	}

	//if the database exists in rules
	if r.GetSchemaRule(srr.DB) != nil {
		return nil, fmt.Errorf("schema %s rule duplicate", srr.DB)
	} else {
		r.SchemaRules.Store(srr.DB, srr)
	}

	return srr, nil
}

func (r *Router) LoadShardingSchemaRule(schemaConfig *config.SchemaConfig, hostGroupClusterConf *config.HostGroupCluster) (*SchemaRouteRule, error) {
	if !util.IsPowerOfTwo(schemaConfig.SchemaSharding) {
		return nil, fmt.Errorf("schema:%s, schemaSharding %v is not power of 2", schemaConfig.Name, schemaConfig.SchemaSharding)
	}
	if !util.IsPowerOfTwo(schemaConfig.TableSharding) {
		return nil, fmt.Errorf("schema:%s, tableSharding %v is not power of 2", schemaConfig.Name, schemaConfig.TableSharding)
	}
	if !util.IsPowerOfTwo(len(hostGroupClusterConf.ShardingHostGroups)) {
		return nil, fmt.Errorf("schema:%s, the len of ShardingHostGroups %v is not power of 2", schemaConfig.Name, len(hostGroupClusterConf.ShardingHostGroups))
	}
	if len(hostGroupClusterConf.ShardingHostGroups) > schemaConfig.SchemaSharding {
		return nil, fmt.Errorf("schema:%s, the len of ShardingHostGroups %v is bigger than SchemaSharding %v, which will cause waste",
			schemaConfig.Name,
			len(hostGroupClusterConf.ShardingHostGroups),
			schemaConfig.SchemaSharding)
	}

	if schemaConfig.SchemaSharding*schemaConfig.TableSharding > 1024 {
		return nil, fmt.Errorf("the product:[%d] of schemaSharding:[%d] and tableSharding:[%d] is bigger than 1024, which is not allowed",
			schemaConfig.SchemaSharding*schemaConfig.TableSharding,
			schemaConfig.SchemaSharding,
			schemaConfig.TableSharding)
	}

	srr := new(SchemaRouteRule)
	srr.DB = schemaConfig.Name
	srr.TableToHostGroupNode = make(map[int]string)
	srr.TableToSchemaIndex = make(map[int]int)
	srr.SchemaToHostGroupNode = make(map[int]string)
	srr.Type = TypeSharding

	//// 分库情况下，每个HostGroup中有多少表
	//tablePerGroup := schemaConfig.SchemaSharding * schemaConfig.TableSharding / len(schemaConfig.ShardingHostGroups)
	// 分库的情况下，每个HostGroup中有多少schema
	schemaPerGroup := schemaConfig.SchemaSharding / len(hostGroupClusterConf.ShardingHostGroups)
	// 分库的情况下，是否一个库中只有同一个单表
	srr.OneTablePerSchema = schemaConfig.TableSharding == 1

	srr.ShardingSchemaCount = schemaConfig.SchemaSharding
	srr.ShardingTableCount = schemaConfig.TableSharding * schemaConfig.SchemaSharding

	srr.HasNonshardingHostGroup = len(hostGroupClusterConf.NonshardingHostGroup) != 0

	srr.HostGroupNodes = hostGroupClusterConf.ShardingHostGroups

	if srr.HasNonshardingHostGroup {
		srr.NonShardingIndex = NonshardingTableIndex
		srr.HostGroupNodes = append(srr.HostGroupNodes, hostGroupClusterConf.NonshardingHostGroup)
	}

	// 分成两个部分:
	//1.分库
	tableIndex := 0
	var schemaIndex int
	for i := 0; i < schemaConfig.SchemaSharding; i++ {

		hostGroup := srr.HostGroupNodes[i/schemaPerGroup]

		// 不管分库还是不分库，都从1开始
		schemaIndex = i + 1

		srr.FullSchemaIndexes = append(srr.FullSchemaIndexes, schemaIndex)
		srr.SchemaToHostGroupNode[schemaIndex] = hostGroup

		for j := 0; j < schemaConfig.TableSharding; j++ {
			srr.FullTableIndexes = append(srr.FullTableIndexes, tableIndex)
			srr.TableToHostGroupNode[tableIndex] = hostGroup
			srr.TableToSchemaIndex[tableIndex] = schemaIndex
			tableIndex++
		}
	}
	//2.非分库
	if srr.HasNonshardingHostGroup {
		srr.SchemaToHostGroupNode[NonshardingSchemaIndex] = hostGroupClusterConf.NonshardingHostGroup
		srr.TableToHostGroupNode[NonshardingTableIndex] = srr.HostGroupNodes[0]
		srr.TableToSchemaIndex[NonshardingTableIndex] = NonshardingSchemaIndex
		srr.FullSchemaIndexes = append(srr.FullSchemaIndexes, NonshardingSchemaIndex)
	}

	if err := parseShardMap(srr); err != nil {
		return nil, err
	}

	//if the database exists in rules
	if r.GetSchemaRule(srr.DB) != nil {
		return nil, fmt.Errorf("schema %s rule duplicate", srr.DB)
	} else {
		r.SchemaRules.Store(srr.DB, srr)
	}

	return srr, nil
}

func (r *Router) UnloadSchemaRule(schemaName string) error {
	r.SchemaRules.Delete(schemaName)

	return nil
}

func (r *Router) DeleteSchemaRule(schemaName string) error {
	r.SchemaRules.Delete(schemaName)
	return nil
}

func parseShardMap(r *SchemaRouteRule) error {
	r.Shard = make(map[string]Shard)
	r.Shard[HashRuleIntType] = &IntHashShard{SlotNum: 1024, ShardLength: 1024 / r.ShardingTableCount}
	r.Shard[HashRuleStringType] = &StringHashShard{SlotNum: 1024, ShardLength: 1024 / r.ShardingTableCount}

	return nil
}

func containsNode(nodes []string, node string) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

// build a router plan
// add memAlloc just for the sake of higher performance
func (r *Router) BuildPlan(db string, client string, statement sqlparser.Statement, sql string, connectionMultiRoutePermit bool, schemaMultiRoutePermit bool, allocMem []byte) (*Plan, error) {
	var plan *Plan
	var err error
	//因为实现Statement接口的方法都是指针类型，所以type对应类型也是指针类型
	switch stmt := statement.(type) {
	case *sqlparser.Select:
		plan, err = r.buildSelectPlan(db, stmt, allocMem)
	case *sqlparser.Update:
		plan, err = r.buildUpdatePlan(db, stmt, allocMem)
	case *sqlparser.Insert:
		plan, err = r.buildInsertPlan(db, stmt, allocMem)
	case *sqlparser.Replace:
		plan, err = r.buildReplacePlan(db, stmt, allocMem)
	case *sqlparser.Delete:
		plan, err = r.buildDeletePlan(db, stmt, allocMem)
	case *sqlparser.Truncate:
		plan, err = r.buildTruncatePlan(db, stmt, allocMem)
	case *sqlparser.DDL:
		plan, err = r.buildDdlPlan(db, stmt, sql, allocMem)
	default:
		err = errcode.ErrSQLNotSupport
	}

	// add execute limitation for multiple nodes sql (no sharding key)
	if plan != nil && !plan.hasShardingKey && len(plan.ShardingNodes) > 1 { // 理论上99%的SQL不会走入该if块，所以放在第一层，保证执行速度
		shardingInfo := fmt.Sprintf(
			"No sharding key is provided. Schema's multiRoutePermitted: %v, Connection's multiRoutePermitted: %v,",
			schemaMultiRoutePermit,
			connectionMultiRoutePermit)
		clientInfo := fmt.Sprintf("schema: [%s], client: [%s], sql: [%s].", db, client, sql)
		if !(schemaMultiRoutePermit && connectionMultiRoutePermit) {
			glog.Glog.Errorf("%s, %s, Execution is rejected by proxy due to current config.", shardingInfo, clientInfo)
			return nil, errcode.BuildError(errcode.ErrMultiRouteNotPermitted, schemaMultiRoutePermit, connectionMultiRoutePermit)
		} else {
			glog.Glog.Warnf("%s, %s", shardingInfo, clientInfo)
		}
	}

	return plan, err
}

func (r *Router) buildSelectPlan(db string, statement sqlparser.Statement, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	var where *sqlparser.Where
	var err error
	var tableName string

	stmt := statement.(*sqlparser.Select)
	switch v := (stmt.From[0]).(type) {
	case *sqlparser.AliasedTableExpr:
		tableName = sqlparser.String(v.Expr)
	case *sqlparser.JoinTableExpr:
		if ate, ok := (v.LeftExpr).(*sqlparser.AliasedTableExpr); ok {
			tableName = sqlparser.String(ate.Expr)
		} else {
			tableName = sqlparser.String(v)
		}
	default:
		tableName = sqlparser.String(v)
	}

	plan.SchemaRouteRule = r.GetSchemaRule(db) //根据表名获得分表规则
	if plan.SchemaRouteRule == nil {
		// error
		return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
	}

	where = stmt.Where
	var routeTableIndexes []int

	plan.TableRouteRule = plan.SchemaRouteRule.GetTableRule(tableName)
	if plan.TableRouteRule == nil {
		// nonsharding
		// 不是分库表，如果改schema配置了nonsharding的Hostgroup，那就放入改hostGroup中，如果没有配置，那么就是表名出错了
		if !plan.SchemaRouteRule.HasNonshardingHostGroup {
			return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, tableName)
		}
		routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)
	} else {
		// sharding
		if where != nil {
			plan.Criteria = where.Expr //路由条件
			routeTableIndexes, err = plan.calRouteIndexes()
			if err != nil {
				return nil, err
			}
		} else {
			//if shard select without where, send to all nodes and all tables
			routeTableIndexes = plan.SchemaRouteRule.FullTableIndexes
		}
	}

	if len(routeTableIndexes) == 0 {
		return nil, errcode.ErrCalcRoute
	}
	//generate sql,如果routeTableindexs为空则表示不分表，不分表则发default node
	err = r.generateSelectShardingNodes(routeTableIndexes, plan, tableName, stmt, allocMem)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *Router) buildInsertPlan(db string, statement sqlparser.Statement, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	plan.Rows = make(map[int]sqlparser.Values)
	stmt := statement.(*sqlparser.Insert)
	if _, ok := stmt.Rows.(sqlparser.SelectStatement); ok {
		return nil, errcode.ErrSelectInInsert
	}

	if stmt.Columns == nil {
		return nil, errcode.ErrIRNoColumns
	}

	//根据sql语句的表，获得对应的分片规则
	plan.SchemaRouteRule = r.GetSchemaRule(db)
	if plan.SchemaRouteRule == nil {
		// error
		return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
	}
	table := sqlparser.String(stmt.Table)
	plan.TableRouteRule = plan.SchemaRouteRule.GetTableRule(table)

	var routeTableIndexes []int
	if plan.TableRouteRule == nil {
		// nonsharding
		if !plan.SchemaRouteRule.HasNonshardingHostGroup {
			return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, table)
		}
		plan.Rows[plan.SchemaRouteRule.NonShardingIndex] = stmt.Rows.(sqlparser.Values)
		routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)
	} else {
		// sharding
		err := plan.GetIRKeyIndex(stmt.Columns)
		if err != nil {
			return nil, err
		}

		if stmt.OnDup != nil {
			err := checkUpdateExprs(plan, sqlparser.UpdateExprs(stmt.OnDup))
			if err != nil {
				return nil, err
			}
		}

		plan.Criteria, err = plan.checkValuesType(stmt.Rows.(sqlparser.Values))
		if err != nil {
			return nil, err
		}

		routeTableIndexes, err = plan.calRouteIndexes()
		if err != nil {
			return nil, err
		}
	}

	err := r.generateInsertShardingNodes(routeTableIndexes, plan, table, stmt, allocMem)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *Router) buildUpdatePlan(db string, statement sqlparser.Statement, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	var where *sqlparser.Where

	stmt := statement.(*sqlparser.Update)
	plan.SchemaRouteRule = r.GetSchemaRule(db)
	if plan.SchemaRouteRule == nil {
		// error
		return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
	}

	table := sqlparser.String(stmt.Table)
	plan.TableRouteRule = plan.SchemaRouteRule.GetTableRule(table)

	var routeTableIndexes []int
	if plan.TableRouteRule == nil {
		// nonsharding
		if !plan.SchemaRouteRule.HasNonshardingHostGroup {
			return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, table)
		}
		routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)
	} else {
		// sharding
		err := checkUpdateExprs(plan, stmt.Exprs)
		if err != nil {
			return nil, err
		}

		where = stmt.Where
		if where != nil {
			plan.Criteria = where.Expr //路由条件
			routeTableIndexes, err = plan.calRouteIndexes()
			if err != nil {
				return nil, err
			}
		} else {
			//if shard update without where,send to all nodes and all tables
			routeTableIndexes = plan.SchemaRouteRule.FullTableIndexes
		}

		if len(routeTableIndexes) == 0 {
			return nil, errcode.ErrCalcRoute
		}
	}

	//generate sql,如果routeTableindexs为空则表示不分表，不分表则发default node
	err := r.generateUpdateShardingNodes(routeTableIndexes, plan, table, stmt, allocMem)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *Router) buildDeletePlan(db string, statement sqlparser.Statement, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	var where *sqlparser.Where
	var err error

	stmt := statement.(*sqlparser.Delete)
	plan.SchemaRouteRule = r.GetSchemaRule(db)
	if plan.SchemaRouteRule == nil {
		// error
		return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
	}
	table := sqlparser.String(stmt.Table)
	plan.TableRouteRule = plan.SchemaRouteRule.GetTableRule(table)
	where = stmt.Where

	var routeTableIndexes []int
	if plan.TableRouteRule == nil {
		// nonsharding
		if !plan.SchemaRouteRule.HasNonshardingHostGroup {
			return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, table)
		}
		routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)
	} else {
		// sharding
		if where != nil {
			plan.Criteria = where.Expr //路由条件
			routeTableIndexes, err = plan.calRouteIndexes()
			if err != nil {
				return nil, err
			}
		} else {
			routeTableIndexes = plan.SchemaRouteRule.FullTableIndexes
		}

		if len(routeTableIndexes) == 0 {
			return nil, errcode.ErrCalcRoute
		}
	}

	//generate sql,如果routeTableindexs为空则表示不分表，不分表则发default node
	err = r.generateDeleteShardingNodes(routeTableIndexes, plan, table, stmt, allocMem)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *Router) buildTruncatePlan(db string, statement sqlparser.Statement, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	var err error

	stmt := statement.(*sqlparser.Truncate)
	plan.SchemaRouteRule = r.GetSchemaRule(db)
	if plan.SchemaRouteRule == nil {
		// error
		return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
	}

	//send to all nodes and all tables
	routeTableIndexes := plan.SchemaRouteRule.FullTableIndexes

	if len(routeTableIndexes) == 0 {
		return nil, errcode.ErrCalcRoute
	}
	//generate sql,如果routeTableindexs为空则表示不分表，不分表则发default node
	err = r.generateTruncateShardingNodes(routeTableIndexes, plan, sqlparser.String(stmt.Table), stmt, allocMem)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *Router) buildDdlPlan(db string, statement sqlparser.Statement, sql string, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	var err error

	stmt := statement.(*sqlparser.DDL)

	switch stmt.Action {
	case sqlparser.AST_TABLE_CREATE:
		plan.SchemaRouteRule = r.GetSchemaRule(db)
		if plan.SchemaRouteRule == nil {
			// error
			return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
		}

		tableRule := plan.SchemaRouteRule.GetTableRule(string(stmt.Table))

		var routeTableIndexes []int
		if tableRule == nil {
			if !plan.SchemaRouteRule.HasNonshardingHostGroup {
				return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, string(stmt.Table))
			}
			routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)

		} else {
			// send to all tables
			routeTableIndexes = plan.SchemaRouteRule.FullTableIndexes

			if len(routeTableIndexes) == 0 {
				return nil, errcode.ErrCalcRoute
			}
		}

		err = r.generateCreateTableShardingNodes(routeTableIndexes, plan, string(stmt.Table), stmt, sql, allocMem)
		break
	case sqlparser.AST_INDEX_CREATE:
		plan.SchemaRouteRule = r.GetSchemaRule(db)
		if plan.SchemaRouteRule == nil {
			// error
			return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
		}

		tableRule := plan.SchemaRouteRule.GetTableRule(string(stmt.Table))

		var routeTableIndexes []int
		if tableRule == nil {
			if !plan.SchemaRouteRule.HasNonshardingHostGroup {
				return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, string(stmt.Table))
			}
			routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)

		} else {
			// send to all tables
			routeTableIndexes = plan.SchemaRouteRule.FullTableIndexes

			if len(routeTableIndexes) == 0 {
				return nil, errcode.ErrCalcRoute
			}
		}

		err = r.generateCreateIndexShardingNodes(routeTableIndexes, plan, string(stmt.Table), string(stmt.NewName), stmt, sql, allocMem)
		break
	case sqlparser.AST_DATABASE_CREATE:
		plan.SchemaRouteRule = r.GetSchemaRule(string(stmt.Database))

		if plan.SchemaRouteRule == nil {
			// error
			return nil, errcode.BuildError(errcode.SchemaRuleNotExist, string(stmt.Database))
		}

		// send to all databases
		routeSchemaIndexes := plan.SchemaRouteRule.FullSchemaIndexes

		if len(routeSchemaIndexes) == 0 {
			return nil, errcode.ErrCalcRoute
		}
		err = r.generateCreateDatabaseShardingNodes(routeSchemaIndexes, plan, stmt, sql, allocMem)
		break
	case sqlparser.AST_TABLE_DROP:
		plan.SchemaRouteRule = r.GetSchemaRule(db)
		if plan.SchemaRouteRule == nil {
			// error
			return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
		}
		tableRule := plan.SchemaRouteRule.GetTableRule(string(stmt.Table))
		var routeTableIndexes []int
		if tableRule == nil {
			if !plan.SchemaRouteRule.HasNonshardingHostGroup {
				return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, string(stmt.Table))
			}
			routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)

		} else {
			// send to all tables
			routeTableIndexes = plan.SchemaRouteRule.FullTableIndexes

			if len(routeTableIndexes) == 0 {
				return nil, errcode.ErrCalcRoute
			}
		}

		err = r.generateDropTableShardingNodes(routeTableIndexes, plan, string(stmt.Table), stmt, sql, allocMem)
		break
	case sqlparser.AST_DATABASE_DROP:
		plan.SchemaRouteRule = r.GetSchemaRule(string(stmt.Database))
		if plan.SchemaRouteRule == nil {
			// error
			return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
		}
		// send to all databases
		routeSchemaIndexes := plan.SchemaRouteRule.FullSchemaIndexes

		if len(routeSchemaIndexes) == 0 {
			return nil, errcode.ErrCalcRoute
		}
		err = r.generateDropDatabaseShardingNodes(routeSchemaIndexes, plan, stmt, sql, allocMem)
		break
	default:
		return nil, errcode.ErrSQLNotSupported
	}

	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *Router) buildReplacePlan(db string, statement sqlparser.Statement, allocMem []byte) (*Plan, error) {
	plan := &Plan{}
	plan.Rows = make(map[int]sqlparser.Values)

	stmt := statement.(*sqlparser.Replace)
	if _, ok := stmt.Rows.(sqlparser.SelectStatement); ok {
		panic(sqlparser.NewParserError("select in replace not allowed"))
	}

	if stmt.Columns == nil {
		return nil, errcode.ErrIRNoColumns
	}

	plan.SchemaRouteRule = r.GetSchemaRule(db)
	if plan.SchemaRouteRule == nil {
		// error
		return nil, errcode.BuildError(errcode.SchemaRuleNotExist, db)
	}
	table := sqlparser.String(stmt.Table)
	plan.TableRouteRule = plan.SchemaRouteRule.GetTableRule(table)

	var routeTableIndexes []int
	if plan.TableRouteRule == nil {
		// nonsharding
		if !plan.SchemaRouteRule.HasNonshardingHostGroup {
			return nil, errcode.BuildError(errcode.NonshardingConfNotExist, db, table)
		}
		routeTableIndexes = append(routeTableIndexes, plan.SchemaRouteRule.NonShardingIndex)
	} else {
		err := plan.GetIRKeyIndex(stmt.Columns)
		if err != nil {
			return nil, err
		}

		plan.Criteria, err = plan.checkValuesType(stmt.Rows.(sqlparser.Values))
		if err != nil {
			return nil, err
		}

		routeTableIndexes, err = plan.calRouteIndexes()
		if err != nil {
			return nil, err
		}
	}

	err := r.generateReplaceShardingNodes(routeTableIndexes, plan, table, stmt, allocMem)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

//rewrite select sql
func (r *Router) rewriteSelectSql(plan *Plan, node *sqlparser.Select, schemaIndex int, tableIndex int, originalTable string, allocMem []byte) string {
	var prefix string

	buf := sqlparser.NewTrackedBuffer(allocMem, nil)

	buf.Fprintf("select %v%s",
		node.Comments,
		node.Distinct,
	)

	schema := getPrefixedSchema(plan, schemaIndex)
	table := getSuffixedTable(plan, tableIndex, originalTable)

	//rewrite select expr
	for _, expr := range node.SelectExprs {
		switch v := expr.(type) {
		case *sqlparser.StarExpr:
			//for shardTable.*,need replace table into shardTable_xxxx.
			if string(v.TableName) == originalTable {
				fmt.Fprintf(buf, "%s%s.%s.*",
					prefix,
					schema,
					table,
				)
			} else {
				buf.Fprintf("%s%v", prefix, expr)
			}
		case *sqlparser.NonStarExpr:
			//rewrite shardTable.column as a
			//into shardTable_xxxx.column as a
			if colName, ok := v.Expr.(*sqlparser.ColName); ok {
				if string(colName.Qualifier) == originalTable {
					fmt.Fprintf(buf, "%s%s.%s.%s",
						prefix,
						schema,
						table,
						string(colName.Name),
					)
				} else {
					buf.Fprintf("%s%v", prefix, colName)
				}
				//if expr has as
				if v.As != nil {
					buf.Fprintf(" as %s", v.As)
				}
			} else {
				buf.Fprintf("%s%v", prefix, expr)
			}
		default:
			buf.Fprintf("%s%v", prefix, expr)
		}
		prefix = ", "
	}
	//insert the group columns in the first of select cloumns
	if len(node.GroupBy) != 0 {
		prefix = ","
		for _, n := range node.GroupBy {
			buf.Fprintf("%s%v", prefix, n)
		}
	}
	buf.Fprintf(" from ")
	switch v := (node.From[0]).(type) {
	case *sqlparser.AliasedTableExpr:
		if len(v.As) != 0 {
			fmt.Fprintf(buf, "%s.%s as %s",
				schema,
				table,
				string(v.As),
			)
		} else {
			fmt.Fprintf(buf, "%s.%s",
				schema,
				table,
			)
		}
	case *sqlparser.JoinTableExpr:
		if ate, ok := (v.LeftExpr).(*sqlparser.AliasedTableExpr); ok {
			if len(ate.As) != 0 {
				fmt.Fprintf(buf, "%s.%s as %s",
					schema,
					getSuffixedTable(plan, tableIndex, sqlparser.String(ate.Expr)),
					string(ate.As),
				)
			} else {
				fmt.Fprintf(buf, "%s.%s",
					schema,
					getSuffixedTable(plan, tableIndex, sqlparser.String(ate.Expr)),
				)
			}
		} else {
			fmt.Fprintf(buf, "%s.%s",
				schema,
				getSuffixedTable(plan, tableIndex, sqlparser.String(v.LeftExpr)),
			)
		}
		//buf.Fprintf(" %s %v", v.Join, v.RightExpr)
		if ate, ok := (v.RightExpr).(*sqlparser.AliasedTableExpr); ok {
			if len(ate.As) != 0 {
				fmt.Fprintf(buf, " %s %s.%s as %s",
					v.Join,
					schema,
					getSuffixedTable(plan, tableIndex, sqlparser.String(ate.Expr)),
					string(ate.As),
				)
			}
		}
		if v.On != nil {
			buf.Fprintf(" on %v", v.On)
		}
	default:
		fmt.Fprintf(buf, "%s.%s",
			schema,
			getSuffixedTable(plan, tableIndex, sqlparser.String(node.From[0])),
		)
	}

	prefix = ", "
	for i := 1; i < len(node.From); i++ {
		if ate, ok := (node.From[i]).(*sqlparser.AliasedTableExpr); ok {
			if len(ate.As) != 0 {
				fmt.Fprintf(buf, ", %s.%s as %s",
					schema,
					getSuffixedTable(plan, tableIndex, sqlparser.String(ate.Expr)),
					string(ate.As),
				)
			}
		} else {
			buf.Fprintf("%s%s %v",
				prefix,
				schema,
				node.From[i])
		}
	}

	newLimit, err := node.Limit.RewriteLimit()
	if err != nil {
		//do not change limit
		newLimit = node.Limit
	}
	//rewrite where
	plan.rewriteWhereIn(tableIndex)

	buf.Fprintf("%v%v%v%v%v%s",
		node.Where,
		node.GroupBy,
		node.Having,
		node.OrderBy,
		newLimit,
		node.Lock,
	)
	return buf.String()
}

func (r *Router) generateSelectShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, allocMem []byte) error {
	node, ok := stmt.(*sqlparser.Select)
	if ok == false {
		return errcode.ErrStmtConvert
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		// 由于原生的SQL中，schema和table没有携带分库分表信息，需要解析出[schema]和[table]部分，加入相应的分库分表后缀
		selectSql := r.rewriteSelectSql(plan, node, schemaIndex, tableIndex, table, allocMem)

		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, selectSql, hostGroupName)
	}

	return nil
}

func (r *Router) generateInsertShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, allocMem []byte) error {
	node, ok := stmt.(*sqlparser.Insert)
	if ok == false {
		return errcode.ErrStmtConvert
	}
	tableCount := len(routeTableIndexes)
	if tableCount == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)
		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)

		buf.Fprintf("insert %v%s into %s.%s", node.Comments, node.Ignore, schema, table)
		buf.Fprintf("%v %v%v",
			node.Columns,
			plan.Rows[tableIndex],
			node.OnDup)

		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateUpdateShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, allocMem []byte) error {
	node, ok := stmt.(*sqlparser.Update)
	if ok == false {
		return errcode.ErrStmtConvert
	}
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		table := getSuffixedTable(plan, tableIndex, table)
		schema := getPrefixedSchema(plan, schemaIndex)

		buf.Fprintf("update %v%s.%s",
			node.Comments,
			schema,
			table,
		)
		buf.Fprintf(" set %v%v%v%v",
			node.Exprs,
			node.Where,
			node.OrderBy,
			node.Limit,
		)

		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateDeleteShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, allocMem []byte) error {
	node, ok := stmt.(*sqlparser.Delete)
	if ok == false {
		return errcode.ErrStmtConvert
	}
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}

	tableCount := len(routeTableIndexes)
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)

		buf.Fprintf("delete %vfrom %s.%s",
			node.Comments,
			schema,
			table,
		)
		buf.Fprintf("%v%v%v",
			node.Where,
			node.OrderBy,
			node.Limit,
		)

		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateReplaceShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, allocMem []byte) error {
	node, ok := stmt.(*sqlparser.Replace)
	if ok == false {
		return errcode.ErrStmtConvert
	}
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)

		buf := sqlparser.NewTrackedBuffer(allocMem, nil)
		buf.Fprintf("replace %vinto %s.%s",
			node.Comments,
			schema,
			table,
		)
		buf.Fprintf("%v %v",
			node.Columns,
			plan.Rows[tableIndex],
		)
		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateTruncateShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, allocMem []byte) error {
	node, ok := stmt.(*sqlparser.Truncate)
	if ok == false {
		return errcode.ErrStmtConvert
	}
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)

		buf.Fprintf("truncate %v%s%s.%s",
			node.Comments,
			node.TableOpt,
			schema,
			table,
		)
		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateCreateTableShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, sql string, allocMem []byte) error {
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)
		len := len(sql)

		buf.Fprintf("create table %s.%s %s",
			schema,
			table,
			sql[strings.Index(sql, "("):len],
		)
		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateCreateIndexShardingNodes(routeTableIndexes []int, plan *Plan, table string, index string, stmt sqlparser.Statement, sql string, allocMem []byte) error {
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)
		len := len(sql)

		var unique string
		if strings.Index(sql, " index ") == -1 {
			unique = ""
		} else {
			unique = " unique"
		}

		buf.Fprintf("create%s index %s on %s.%s%s",
			unique,
			index,
			schema,
			table,
			sql[strings.Index(sql, "("):len],
		)
		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateCreateDatabaseShardingNodes(routeSchemaIndexes []int, plan *Plan, ddl *sqlparser.DDL, sql string, allocMem []byte) error {
	schemaCount := len(routeSchemaIndexes)
	if schemaCount == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	for i := 0; i < schemaCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		schemaIndex := routeSchemaIndexes[i]

		schema := getPrefixedSchema(plan, schemaIndex)

		buf.Fprintf("create database %s",
			schema,
		)
		hostGroupName := plan.SchemaRouteRule.SchemaToHostGroupNode[schemaIndex]
		plan.ShardingNodes[schemaIndex] = backend.NewShardingNode(schemaIndex, 0, backend.ShardingTypeSchema, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateDropTableShardingNodes(routeTableIndexes []int, plan *Plan, table string, stmt sqlparser.Statement, sql string, allocMem []byte) error {
	if len(routeTableIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	tableCount := len(routeTableIndexes)
	for i := 0; i < tableCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		tableIndex := routeTableIndexes[i]
		schemaIndex := plan.SchemaRouteRule.TableToSchemaIndex[tableIndex]

		schema := getPrefixedSchema(plan, schemaIndex)
		table := getSuffixedTable(plan, tableIndex, table)

		buf.Fprintf("drop table %s.%s",
			schema,
			table,
		)
		hostGroupName := plan.SchemaRouteRule.TableToHostGroupNode[tableIndex]
		plan.ShardingNodes[tableIndex] = backend.NewShardingNode(schemaIndex, tableIndex, backend.ShardingTypeTable, buf.String(), hostGroupName)
	}

	return nil
}

func (r *Router) generateDropDatabaseShardingNodes(routeSchemaIndexes []int, plan *Plan, stmt sqlparser.Statement, sql string, allocMem []byte) error {
	if len(routeSchemaIndexes) == 0 {
		return errcode.ErrNoShardingNode
	}
	plan.ShardingNodes = make(map[int]*backend.ShardingNode)
	schemaCount := len(routeSchemaIndexes)
	for i := 0; i < schemaCount; i++ {
		buf := sqlparser.NewTrackedBuffer(allocMem, nil)

		schemaIndex := routeSchemaIndexes[i]

		db := getPrefixedSchema(plan, schemaIndex)

		buf.Fprintf("drop database %s",
			db,
		)
		hostGroupName := plan.SchemaRouteRule.SchemaToHostGroupNode[schemaIndex]
		plan.ShardingNodes[schemaIndex] = backend.NewShardingNode(schemaIndex, 0, backend.ShardingTypeSchema, buf.String(), hostGroupName)
	}

	return nil
}

func getSuffixedTable(plan *Plan, tableIndex int, table string) string {
	if plan.SchemaRouteRule.Type == TypeCustody {
		return table
	} else {
		// 如果一个库中只有一张表，或者这个表是非分库表，那这张表就不带后缀
		if plan.SchemaRouteRule.OneTablePerSchema || tableIndex == NonshardingTableIndex {
			return table
		} else {
			return fmt.Sprintf("%s%d", table, tableIndex)
		}
	}
}

func getPrefixedSchema(plan *Plan, schemaIndex int) string {
	if plan.SchemaRouteRule.Type == TypeCustody {
		return plan.SchemaRouteRule.DB
	} else {
		//return fmt.Sprintf("%s%d", plan.SchemaRouteRule.DB, schemaIndex)
		var sb strings.Builder
		sb.WriteString(plan.SchemaRouteRule.DB)
		sb.WriteString(strconv.Itoa(schemaIndex))
		return sb.String()
	}
}
