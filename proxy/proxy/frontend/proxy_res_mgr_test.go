package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/common/config"
	"sync"
	"testing"
)

func TestHostGroupAddDelete(t *testing.T) {
	var addWg sync.WaitGroup
	addLen := 10
	addWg.Add(addLen)

	for i := 0; i < addLen; i++ {
		hostGroupCfg := new(config.HostGroup)
		hostGroupCfg.Name = fmt.Sprintf("hostGroup10%d", i)
		hostGroupCfg.InitConn = 1
		hostGroupCfg.MaxConn = 1024
		hostGroupCfg.User = "zerodb"
		hostGroupCfg.Password = "zerodb@552208"
		hostGroupCfg.Write = "10.1.22.1:3306,10.1.21.79:3306"

		addHG := func() {
			if _, err := testProxyEngine.AddHostGroup(*hostGroupCfg); err != nil {
				t.Errorf(err.Error())
			}
			addWg.Done()
		}

		go addHG()
	}
	addWg.Wait()

	var DelWg sync.WaitGroup
	DelWg.Add(addLen)

	delHG := func(j int) {
		n := fmt.Sprintf("hostGroup10%d", j)
		if _, err := testProxyEngine.DeleteHostGroup(n); err != nil {
			t.Errorf(err.Error())
		}
		DelWg.Done()
	}

	for j := 0; j < addLen; j++ {
		go delHG(j)
	}

	DelWg.Wait()
}

func TestSchemaUsingHostGroupCluster(t *testing.T) {
	hostGroupClusterCfg := new(config.HostGroupCluster)
	hostGroupClusterCfg.Name = fmt.Sprintf("TestSchemaUsingHostGroupCluster")
	hostGroupClusterCfg.ShardingHostGroups = append(hostGroupClusterCfg.ShardingHostGroups, "hostGroup1")
	hostGroupClusterCfg.ShardingHostGroups = append(hostGroupClusterCfg.ShardingHostGroups, "hostGroup2")

	if _, err := testProxyEngine.AddHostGroupCluster(*hostGroupClusterCfg); err != nil {
		t.Errorf(err.Error())
	}

	schemaCfg := new(config.SchemaConfig)
	schemaCfg.Name = fmt.Sprintf("unit_test_host_group_cluster_schema")
	schemaCfg.HostGroupCluster = "TestSchemaUsingHostGroupCluster"
	schemaCfg.SchemaSharding = 4
	schemaCfg.TableSharding = 8
	if _, err := testProxyEngine.AddSchema(*schemaCfg); err != nil {
		t.Errorf(err.Error())
	}

	if _, err := testProxyEngine.DeleteHostGroupCluster("TestSchemaUsingHostGroupCluster"); err == nil {
		t.Errorf("should have error here")
	}

	hostGroupClusterCfg.ShardingHostGroups = append(hostGroupClusterCfg.ShardingHostGroups, "hostGroup3")
	hostGroupClusterCfg.ShardingHostGroups = append(hostGroupClusterCfg.ShardingHostGroups, "hostGroup4")

	if _, err := testProxyEngine.UpdateHostGroupCluster(*hostGroupClusterCfg); err != nil {
		t.Errorf(err.Error())
	}

	if _, err := testProxyEngine.DeleteSchema(schemaCfg.Name); err != nil {
		t.Errorf(err.Error())
	}

	if _, err := testProxyEngine.DeleteHostGroupCluster(hostGroupClusterCfg.Name); err != nil {
		t.Errorf(err.Error())
	}
}

func TestHostGroupClusterAddUpdateAndDelete(t *testing.T) {
	var addWg sync.WaitGroup
	addLen := 10
	addWg.Add(addLen)

	for i := 0; i < addLen; i++ {
		hostGroupClsuterCfg := new(config.HostGroupCluster)
		hostGroupClsuterCfg.Name = fmt.Sprintf("hostGroupCluster10%d", i)
		hostGroupClsuterCfg.ShardingHostGroups = append(hostGroupClsuterCfg.ShardingHostGroups, "hostGroup5")

		addHGC := func() {
			if _, err := testProxyEngine.AddHostGroupCluster(*hostGroupClsuterCfg); err != nil {
				t.Errorf(err.Error())
			}
			addWg.Done()
		}

		go addHGC()
	}
	addWg.Wait()

	schemaCfg := new(config.SchemaConfig)
	schemaCfg.Name = fmt.Sprintf("unit_test_host_group_cluster")
	schemaCfg.HostGroupCluster = "hostGroupCluster1"
	schemaCfg.SchemaSharding = 4
	schemaCfg.TableSharding = 8
	if _, err := testProxyEngine.AddSchema(*schemaCfg); err != nil {
		t.Errorf(err.Error())
	}

	var DelWg sync.WaitGroup
	DelWg.Add(addLen)

	delHGC := func(j int) {
		n := fmt.Sprintf("hostGroupCluster10%d", j)
		if _, err := testProxyEngine.DeleteHostGroupCluster(n); err != nil {
			t.Errorf(err.Error())
		}
		DelWg.Done()
	}

	for j := 0; j < addLen; j++ {
		go delHGC(j)
	}

	DelWg.Wait()
}

func TestHostGroupDelete(t *testing.T) {
	if _, err := testProxyEngine.DeleteHostGroup("hostGroup1"); err == nil {
		t.Errorf("should have error")
	} else {
		t.Logf("normal err:%v", err)
	}
}

func TestSchemaAddDelete(t *testing.T) {
	var addWg sync.WaitGroup
	addLen := 100
	addWg.Add(addLen)

	for i := 0; i < addLen; i++ {
		schemaCfg := new(config.SchemaConfig)
		schemaCfg.Name = fmt.Sprintf("unit_test%d", i)
		schemaCfg.HostGroupCluster = "cluster_5"
		schemaCfg.SchemaSharding = 4
		schemaCfg.TableSharding = 8

		addSC := func() {
			if _, err := testProxyEngine.AddSchema(*schemaCfg); err != nil {
				t.Errorf(err.Error())
			}
			addWg.Done()
		}

		go addSC()
	}
	addWg.Wait()

	var DelWg sync.WaitGroup
	DelWg.Add(addLen)

	delSC := func(j int) {
		n := fmt.Sprintf("unit_test%d", j)
		if _, err := testProxyEngine.DeleteSchema(n); err != nil {
			t.Errorf(err.Error())
		}
		DelWg.Done()
	}

	for j := 0; j < addLen; j++ {
		go delSC(j)
	}

	DelWg.Wait()
}

func TestTableAddDelete(t *testing.T) {
	var addWg sync.WaitGroup
	addLen := 10
	addWg.Add(addLen)
	schemaName := "zerodb"

	for i := 0; i < addLen; i++ {
		tc := new(config.TableConfig)
		tc.Name = fmt.Sprintf("unit_test_table%d", i)
		tc.ShardingKey = "entity_id"
		tc.Rule = "int"

		addSC := func() {
			if _, err := testProxyEngine.AddTable(*tc, schemaName); err != nil {
				t.Errorf(err.Error())
			}
			addWg.Done()
		}

		go addSC()
	}
	addWg.Wait()

	var DelWg sync.WaitGroup
	DelWg.Add(addLen)

	delSC := func(j int) {
		n := fmt.Sprintf("unit_test_table%d", j)
		if _, err := testProxyEngine.DeleteTable(n, schemaName); err != nil {
			t.Errorf(err.Error())
		}
		DelWg.Done()
	}

	for j := 0; j < addLen; j++ {
		go delSC(j)
	}

	DelWg.Wait()
}

func TestUpdateSchemaRWSplit(t *testing.T) {
	if _, err := testProxyEngine.UpdateSchemaRWSplit("zerodb", false); err != nil {
		t.Errorf(err.Error())
	}

	if _, err := testProxyEngine.UpdateSchemaRWSplit("zerodb000", false); err == nil {
		t.Errorf("should have err here")
	}
}

func TestApplyLatestConfig(t *testing.T) {
	cfg, err := getCfgFromFile("../test-conf/proxy_conf_test.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}
	if err = testProxyEngine.ApplyLatestConfig(*cfg); err != nil {
		t.Errorf(err.Error())
	}

}

func TestStopWritingAbility(*testing.T) {

}

func TestRecoverWritingAbility(*testing.T) {

}
