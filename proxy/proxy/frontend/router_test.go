package frontend

import (
	"testing"

	"fmt"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
)

func TestIntRoute(t *testing.T) {
	sql := "select * from instancedetail where entity_id = 0"
	expected := make(map[int]*backend.ShardingNode)
	expected[0] = &backend.ShardingNode{
		SchemaIndex:  1,
		TableIndex:   0,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order1.instancedetail where entity_id = 0",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 8"
	expected = make(map[int]*backend.ShardingNode)
	expected[1] = &backend.ShardingNode{
		SchemaIndex:  2,
		TableIndex:   1,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order2.instancedetail where entity_id = 8",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 255"
	expected = make(map[int]*backend.ShardingNode)
	expected[31] = &backend.ShardingNode{
		SchemaIndex:  32,
		TableIndex:   31,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order32.instancedetail where entity_id = 255",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 256"
	expected = make(map[int]*backend.ShardingNode)
	expected[32] = &backend.ShardingNode{
		SchemaIndex:  33,
		TableIndex:   32,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order33.instancedetail where entity_id = 256",
		HostGroup:    "hostGroup2",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id in (255, 254)"
	expected = make(map[int]*backend.ShardingNode)
	expected[31] = &backend.ShardingNode{
		SchemaIndex:  32,
		TableIndex:   31,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order32.instancedetail where entity_id in (255, 254)",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id in (255, 256)"
	expected = make(map[int]*backend.ShardingNode)
	expected[31] = &backend.ShardingNode{
		SchemaIndex:  32,
		TableIndex:   31,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order32.instancedetail where entity_id in (255)",
		HostGroup:    "hostGroup1",
	}
	expected[32] = &backend.ShardingNode{
		SchemaIndex:  33,
		TableIndex:   32,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order33.instancedetail where entity_id in (256)",
		HostGroup:    "hostGroup2",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 511"
	expected = make(map[int]*backend.ShardingNode)
	expected[63] = &backend.ShardingNode{
		SchemaIndex:  64,
		TableIndex:   63,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order64.instancedetail where entity_id = 511",
		HostGroup:    "hostGroup2",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 512"
	expected = make(map[int]*backend.ShardingNode)
	expected[64] = &backend.ShardingNode{
		SchemaIndex:  65,
		TableIndex:   64,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order65.instancedetail where entity_id = 512",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 1023"
	expected = make(map[int]*backend.ShardingNode)
	expected[127] = &backend.ShardingNode{
		SchemaIndex:  128,
		TableIndex:   127,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order128.instancedetail where entity_id = 1023",
		HostGroup:    "hostGroup4",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = 1024"
	expected = make(map[int]*backend.ShardingNode)
	expected[0] = &backend.ShardingNode{
		SchemaIndex:  1,
		TableIndex:   0,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order1.instancedetail where entity_id = 1024",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from instancedetail where entity_id = -1"
	expected = make(map[int]*backend.ShardingNode)
	expected[127] = &backend.ShardingNode{
		SchemaIndex:  128,
		TableIndex:   127,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order128.instancedetail where entity_id = -1",
		HostGroup:    "hostGroup4",
	}
	checkPlan(t, testRouter, "order", sql, expected)
}

func TestStringRoute(t *testing.T) {
	sql := "select * from waitingordercrid where customerregister_id = '6c2a2360a76b4d3489405bb8a6e95286'"
	expected := make(map[int]*backend.ShardingNode)
	expected[85] = &backend.ShardingNode{
		SchemaIndex:  86,
		TableIndex:   85,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order86.waitingordercrid where customerregister_id = '6c2a2360a76b4d3489405bb8a6e95286'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from waitingordercrid where customerregister_id = '6c2a2361a76b4d3489405bb8a6e95286'"
	expected = make(map[int]*backend.ShardingNode)
	expected[117] = &backend.ShardingNode{
		SchemaIndex:  118,
		TableIndex:   117,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order118.waitingordercrid where customerregister_id = '6c2a2361a76b4d3489405bb8a6e95286'",
		HostGroup:    "hostGroup4",
	}
	checkPlan(t, testRouter, "order", sql, expected)

	sql = "select * from waitingordercrid where customerregister_id = '6c2a2361a26b4d3489405bb8a6e15217'"
	expected = make(map[int]*backend.ShardingNode)
	expected[17] = &backend.ShardingNode{
		SchemaIndex:  18,
		TableIndex:   17,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from order18.waitingordercrid where customerregister_id = '6c2a2361a26b4d3489405bb8a6e15217'",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "order", sql, expected)
}

func TestMultiRoute(t *testing.T) {
	sql := "select * from instancedetail"
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		t.Fatal(err.Error())
	}
	var connMultiRoutePermit bool
	var schemaMultiRoutePermit bool
	connMultiRoutePermit = true
	schemaMultiRoutePermit = true

	client := "local test"
	_, err = testRouter.BuildPlan("order", client, stmt, sql, connMultiRoutePermit, schemaMultiRoutePermit, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	connMultiRoutePermit = false
	schemaMultiRoutePermit = true
	_, err = testRouter.BuildPlan("order", client, stmt, sql, connMultiRoutePermit, schemaMultiRoutePermit, nil)
	if err == nil {
		t.Fatal("MultiRoute should not be permitted")
	}

	connMultiRoutePermit = true
	schemaMultiRoutePermit = false
	_, err = testRouter.BuildPlan("order", client, stmt, sql, connMultiRoutePermit, schemaMultiRoutePermit, nil)
	if err == nil {
		t.Fatal("MultiRoute should not be permitted")
	}

	connMultiRoutePermit = false
	schemaMultiRoutePermit = false
	_, err = testRouter.BuildPlan("order", client, stmt, sql, connMultiRoutePermit, schemaMultiRoutePermit, nil)
	if err == nil {
		t.Fatal("MultiRoute should not be permitted")
	}
}

func TestSelect1(t *testing.T) {
	// ##################### schema: zerodb #########################
	sql := "select * from order_ins where entity_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from zerodb153.order_ins304 where entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)

	// ##################### schema: member #########################
	sql = "select * from card where entity_id = 'whatever'"
	expected = nil
	expected = make(map[int]*backend.ShardingNode)
	expected[60] = &backend.ShardingNode{
		SchemaIndex:  61,
		TableIndex:   60,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from member61.card where entity_id = 'whatever'",
		HostGroup:    "hostGroup2",
	}
	checkPlan(t, testRouter, "member", sql, expected)

	sql = "select count(1) from card where entity_id = 'whatever'"
	expected = nil
	expected = make(map[int]*backend.ShardingNode)
	expected[60] = &backend.ShardingNode{
		SchemaIndex:  61,
		TableIndex:   60,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select count(1) from member61.card where entity_id = 'whatever'",
		HostGroup:    "hostGroup2",
	}
	checkPlan(t, testRouter, "member", sql, expected)
}

func TestSelect2(t *testing.T) {
	// ##################### schema: zerodb #########################
	sql := "select * from order_sth where order_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[16] = &backend.ShardingNode{
		SchemaIndex:  9,
		TableIndex:   16,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from zerodb9.order_sth16 where order_id = '33'",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestSelectCustody(t *testing.T) {
	// ##################### schema: account #########################
	// 测试托管
	sql := "select * from account_info"
	expected := make(map[int]*backend.ShardingNode)
	expected[NonshardingTableIndex] = &backend.ShardingNode{
		SchemaIndex:  0,
		TableIndex:   NonshardingTableIndex,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select * from account.account_info",
		HostGroup:    "hostGroup5",
	}
	checkPlan(t, testRouter, "account", sql, expected)
}

func TestInsertCustody(t *testing.T) {
	sql := "INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world')"
	expected := make(map[int]*backend.ShardingNode)
	expected[NonshardingTableIndex] = &backend.ShardingNode{
		SchemaIndex:  0,
		TableIndex:   NonshardingTableIndex,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "insert  into account.order_ins(id, entity_id, name) values (6, 123125, 'hello world')",
		HostGroup:    "hostGroup5",
	}
	checkPlan(t, testRouter, "account", sql, expected)
}

func TestSelectAliased(t *testing.T) {
	sql := "select oi.id from order_ins as oi where oi.entity_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select oi.id from zerodb153.order_ins304 as oi where oi.entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestSelectJoin(t *testing.T) {
	sql := "SELECT oi.* FROM order_ins AS oi INNER JOIN test_table AS t on t.id = oi.id  WHERE oi.entity_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select oi.* from zerodb153.order_ins304 as oi join zerodb153.test_table304 as t on t.id = oi.id where oi.entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestSelectWhereJoin(t *testing.T) {
	sql := "select id.id, id.detail from ins_detail id, order_ins oi where oi.entity_id = id.entity_id and oi.entity_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "select id.id, id.detail from zerodb153.ins_detail304 as id, zerodb153.order_ins304 as oi where oi.entity_id = id.entity_id and oi.entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestFullTable(t *testing.T) {
	sql := "select * from order_ins"
	expected := constructFullShardingNodes(256*2, backend.ShardingTypeTable)
	checkFullNodePlan(t, testRouter, "zerodb", sql, expected)

	sql = "CREATE TABLE `order_ins` (" +
		"`id`        INT(11) NOT NULL," +
		"`entity_id` BIGINT(20)  DEFAULT NULL," +
		"`name`      VARCHAR(45) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE = InnoDB DEFAULT CHARSET = latin1"
	expected = constructFullShardingNodes(256*2, backend.ShardingTypeTable)
	checkFullNodePlan(t, testRouter, "zerodb", sql, expected)

	sql = "drop table order_ins"
	expected = constructFullShardingNodes(256*2, backend.ShardingTypeTable)
	checkFullNodePlan(t, testRouter, "zerodb", sql, expected)

	sql = "TRUNCATE order_ins"
	expected = constructFullShardingNodes(256*2, backend.ShardingTypeTable)
	checkFullNodePlan(t, testRouter, "zerodb", sql, expected)
}

func TestFullSchema(t *testing.T) {
	sql := "create database zerodb"
	expected := constructFullShardingNodes(256+1, backend.ShardingTypeSchema)
	checkFullNodePlan(t, testRouter, "zerodb", sql, expected)

	sql = "drop database zerodb"
	expected = constructFullShardingNodes(256+1, backend.ShardingTypeSchema)
	checkFullNodePlan(t, testRouter, "zerodb", sql, expected)
}

func TestUpdate(t *testing.T) {
	sql := "update order_ins set name = 'test' where entity_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "update zerodb153.order_ins304 set name = 'test' where entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestReplace(t *testing.T) {
	sql := "REPLACE INTO order_ins (id, entity_id, name) values(10086, '33', 'test')"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "replace into zerodb153.order_ins304(id, entity_id, name) values (10086, '33', 'test')",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestInsert(t *testing.T) {
	sql := "insert into order_ins (id, entity_id, name) values(10086, '33', 'test')"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "insert  into zerodb153.order_ins304(id, entity_id, name) values (10086, '33', 'test')",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)

	sql = "insert into order_ins (id, entity_id, name) values(10086, '33', 'test1'), (10087, '33', 'test2')"
	expected = nil
	expected = make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "insert  into zerodb153.order_ins304(id, entity_id, name) values (10086, '33', 'test1'), (10087, '33', 'test2')",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)

	sql = "insert into order_ins (id, entity_id, name) values(10086, '33', 'test1'), (10087, '1', 'test2')"
	expected = nil
	expected = make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "insert  into zerodb153.order_ins304(id, entity_id, name) values (10086, '33', 'test1')",
		HostGroup:    "hostGroup3",
	}
	expected[24] = &backend.ShardingNode{
		SchemaIndex:  13,
		TableIndex:   24,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "insert  into zerodb13.order_ins24(id, entity_id, name) values (10087, '1', 'test2')",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestDelete(t *testing.T) {
	sql := "delete from order_ins where entity_id = '33'"
	expected := make(map[int]*backend.ShardingNode)
	expected[304] = &backend.ShardingNode{
		SchemaIndex:  153,
		TableIndex:   304,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "delete from zerodb153.order_ins304 where entity_id = '33'",
		HostGroup:    "hostGroup3",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func TestNonsharding(t *testing.T) {
	sql := "drop table hello_world"
	expected := make(map[int]*backend.ShardingNode)
	expected[-1] = &backend.ShardingNode{
		SchemaIndex:  0,
		TableIndex:   -1,
		ShardingType: backend.ShardingTypeTable,
		ShardingSQL:  "drop table zerodb0.hello_world",
		HostGroup:    "hostGroup1",
	}
	checkPlan(t, testRouter, "zerodb", sql, expected)
}

func constructFullShardingNodes(nodeCount int, shardingType int) map[int]*backend.ShardingNode {
	expected := make(map[int]*backend.ShardingNode)

	if shardingType == backend.ShardingTypeTable {
		for i := 0; i < nodeCount; i++ {
			expected[i] = &backend.ShardingNode{
				SchemaIndex:  i/2 + 1,
				TableIndex:   i,
				ShardingType: shardingType,
				HostGroup:    "hostGroup3",
			}
		}
	} else {
		for i := 0; i < nodeCount; i++ {
			expected[i] = &backend.ShardingNode{
				SchemaIndex:  i,
				TableIndex:   0,
				ShardingType: shardingType,
				HostGroup:    "hostGroup3",
			}
		}
	}
	return expected
}

func assertEntirelyEqual(sql string, actual, expected map[int]*backend.ShardingNode) (bool, string) {
	if len(actual) != len(expected) {
		return false, fmt.Sprintf("sql:[%s] \n[%d] actual len \n[%d] expected len", sql, len(actual), len(expected))
	}
	if len(actual) == 0 {
		return true, ""
	}

	for key, val1 := range actual {
		val2 := expected[key]
		if val2 == nil {
			return false, fmt.Sprintf("key:[%d] sql:[%s] key [%d] doesn't exist.", key, sql, key)
		}
		if val1.SchemaIndex != val2.SchemaIndex {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%d] actual.SchemaIndex \n[%d] expected.SchemaIndex", key, sql, val1.SchemaIndex, val2.SchemaIndex)
		}

		if val1.TableIndex != val2.TableIndex {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%d] actual.TableIndex \n[%d] expected.TableIndex", key, sql, val1.TableIndex, val2.TableIndex)
		}

		if val1.ShardingSQL != val2.ShardingSQL {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%s] actual.ShardingSQL \n[%s] expected.ShardingSQL", key, sql, val1.ShardingSQL, val2.ShardingSQL)
		}

		if val1.HostGroup != val2.HostGroup {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%s] actual.HostGroup \n[%s] expected.HostGroup", key, sql, val1.HostGroup, val2.HostGroup)
		}

		if val1.ShardingType != val2.ShardingType {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%d] actual.ShardingType \n[%d] expected.ShardingType", key, sql, val1.ShardingType, val2.ShardingType)
		}
	}

	return true, ""
}

func assertPartiallyEqual(sql string, actual, expected map[int]*backend.ShardingNode) (bool, string) {
	if len(actual) != len(expected) {
		return false, fmt.Sprintf("sql:[%s] \n[%d] actual len \n[%d] expected len", sql, len(actual), len(expected))
	}
	if len(actual) == 0 {
		return true, ""
	}

	for key, val1 := range actual {
		val2 := expected[key]
		if val2 == nil {
			return false, fmt.Sprintf("key:[%d] sql:[%s] key [%d] doesn't exist.", key, sql, key)
		}
		if val1.SchemaIndex != val2.SchemaIndex {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%d] actual.SchemaIndex \n[%d] expected.SchemaIndex", key, sql, val1.SchemaIndex, val2.SchemaIndex)
		}
		if val1.TableIndex != val2.TableIndex {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%d] actual.TableIndex \n[%d] expected.TableIndex", key, sql, val1.SchemaIndex, val2.SchemaIndex)
		}

		//if val1.HostGroup != val2.HostGroup {
		//    return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%s] actual.HostGroup \n[%s] expected.HostGroup", key, sql, val1.HostGroup, val2.HostGroup)
		//}

		if val1.ShardingType != val2.ShardingType {
			return false, fmt.Sprintf("key:[%d] sql:[%s] \n[%d] actual.ShardingType \n[%d] expected.ShardingType", key, sql, val1.ShardingType, val2.ShardingType)
		}
	}

	return true, ""
}

func checkPlan(t *testing.T, r *Router, db string, sql string, expectedShardingNodes map[int]*backend.ShardingNode) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		t.Fatal(err.Error())
	}
	status := true
	plan, err := r.BuildPlan(db, "local test", stmt, sql, status, status, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	equal, diff := assertEntirelyEqual(sql, plan.ShardingNodes, expectedShardingNodes)

	if !equal {
		//err := fmt.Errorf("RouteTableIndexs=%v but tableIndexs=%v",
		//    plan.ShardingNodes, expectedShardingNodes)
		t.Fatal(diff)

	}

}

func checkFullNodePlan(t *testing.T, r *Router, db string, sql string, expectedShardingNodes map[int]*backend.ShardingNode) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		t.Fatal(err.Error())
	}
	status := true
	plan, err := r.BuildPlan(db, "local test", stmt, sql, status, status, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	equal, diff := assertPartiallyEqual(sql, plan.ShardingNodes, expectedShardingNodes)

	if !equal {
		//err := fmt.Errorf("RouteTableIndexs=%v but tableIndexs=%v",
		//    plan.ShardingNodes, expectedShardingNodes)
		t.Fatal(diff)

	}

}
