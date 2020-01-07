package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"testing"
)

func TestUseShardingSchema(t *testing.T) {
	schema := "zerodb"
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", ""); err != nil {
		t.Errorf(err.Error())
	}
	srr := testProxyEngine.GetRouter().GetSchemaRule(schema)
	if srr.Type != TypeSharding {
		t.Errorf("should be TypeSharding")
	}

	if _, err := c.Execute("use " + schema); err != nil {
		t.Error(err.Error())
	}
}

func TestUseCustodySchema(t *testing.T) {
	schema := "custody"
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", ""); err != nil {
		t.Errorf(err.Error())
	}
	srr := testProxyEngine.GetRouter().GetSchemaRule(schema)
	if srr.Type != TypeCustody {
		t.Errorf("should be TypeCustody")
	}
	if _, err := c.Execute("use " + schema); err != nil {
		t.Error(err.Error())
	}
}
