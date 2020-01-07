package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"testing"
)

func TestShowDatabases(t *testing.T) {
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", ""); err != nil {
		t.Errorf(err.Error())
	}

	r, err := c.Execute("show databases")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(r.Values) != 5 {
		t.Error("database count not right")
	}
}

func TestShowTables(t *testing.T) {
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", ""); err != nil {
		t.Errorf(err.Error())
	}

	_, err := c.Execute("show tables")
	if err == nil {
		t.Error("should have use db error")
		return
	}
	c.Execute("use zerodb")

	r, err := c.Execute("show tables")
	if err != nil {
		t.Error(err.Error())
		return
	}
	expected := 7
	if len(r.Values) != expected {
		t.Error("table count not right")
	}
}
