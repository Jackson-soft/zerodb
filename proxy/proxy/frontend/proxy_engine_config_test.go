package frontend

import (
	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	hgclusters = []config.HostGroupCluster{}
	hostGroups = []config.HostGroup{}
)

func TestApplyBackConfig(t *testing.T) {
	cfg, err := getCfgFromFile("../test-conf/proxy_conf_test.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}
	if testProxyEngine.ApplyBackConfig(cfg); err != nil {
		t.Errorf(err.Error())
	}
	var hostGroupClusterNames []string
	for _, hg := range cfg.HostGroupClusters {
		hostGroupClusterNames = append(hostGroupClusterNames, hg.Name)
	}
	testProxyEngine.HostGroupClusters.Range(func(key, value interface{}) bool {
		assert.NotNil(t, value)
		assert.Contains(t, cfg.HostGroupClusters, *value.(*config.HostGroupCluster))
		return true
	})
	testProxyEngine.HostGroupNodes.Range(func(key, value interface{}) bool {
		assert.NotNil(t, value)
		h := *value.(*backend.HostGroupNode)
		assert.Contains(t, cfg.HostGroups, h.Cfg)
		assert.NotNil(t, h.Write)
		return true
	})
	testProxyEngine.schemaNodes.Range(func(key, value interface{}) bool {
		assert.NotNil(t, value)
		s := *value.(*backend.SchemaNode)
		assert.Contains(t, hostGroupClusterNames, s.HostGroupCluster)
		//assert.Contains(t, cfg.SchemaConfigs, s.Name)
		return true
	})
	//TODO add more unit test
	/*	if testProxyEngine.router == nil {
			t.Errorf("rollback router error ")
		}
		if _, ok := testProxyEngine.schemaNodes.Load("zerodb"); !ok {
			t.Errorf("rollback schemaNodes error")
		}
		if _, ok := testProxyEngine.HostGroupNodes.Load("hostGroup1"); !ok {
			t.Errorf("rollback hostgroup error")
		}*/
}
