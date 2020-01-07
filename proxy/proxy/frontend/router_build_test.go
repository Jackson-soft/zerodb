package frontend

import (
	"testing"
)

func TestRouteBuild1(t *testing.T) {
	schemaRule := testRouter.GetSchemaRule("test_schema_with_nonsharding_host_group")
	if len(schemaRule.FullSchemaIndexes) != 5 {
		t.Error("FullSchemaIndexes length is not right")
		return
	}
	if schemaRule.SchemaToHostGroupNode[NonshardingSchemaIndex] != "hostGroup1" {
		t.Error("the hostGroup of NonshardingSchemaIndex schema is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[NonshardingTableIndex] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[0] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[1] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[2] != "hostGroup2" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[3] != "hostGroup2" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[4] != "hostGroup3" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[5] != "hostGroup3" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[6] != "hostGroup4" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[7] != "hostGroup4" {
		t.Error("the hostGroup of table is not right")
		return
	}

	schemaRule = testRouter.GetSchemaRule("test_schema_without_nonsharding_host_group")
	if len(schemaRule.FullSchemaIndexes) != 4 {
		t.Error("FullSchemaIndexes length is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[0] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[1] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[2] != "hostGroup2" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[3] != "hostGroup2" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[4] != "hostGroup3" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[5] != "hostGroup3" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[6] != "hostGroup4" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[7] != "hostGroup4" {
		t.Error("the hostGroup of table is not right")
		return
	}

	schemaRule = testRouter.GetSchemaRule("test_schema_one_table_sharding")
	if !schemaRule.OneTablePerSchema {
		t.Error("OneTablePerSchema should be true")
		return
	}
}

func TestRouteBuild2(t *testing.T) {
	schemaRule := testRouter.GetSchemaRule("test_schema_without_nonsharding_host_group")
	if len(schemaRule.FullSchemaIndexes) != 4 {
		t.Error("FullSchemaIndexes length is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[0] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[1] != "hostGroup1" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[2] != "hostGroup2" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[3] != "hostGroup2" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[4] != "hostGroup3" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[5] != "hostGroup3" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[6] != "hostGroup4" {
		t.Error("the hostGroup of table is not right")
		return
	}

	if schemaRule.TableToHostGroupNode[7] != "hostGroup4" {
		t.Error("the hostGroup of table is not right")
		return
	}

}

func TestRouteBuild3(t *testing.T) {
	schemaRule := testRouter.GetSchemaRule("test_schema_one_table_sharding")
	if !schemaRule.OneTablePerSchema {
		t.Error("OneTablePerSchema should be true")
		return
	}
}
