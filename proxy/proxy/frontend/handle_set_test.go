package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"testing"
)

func TestSet(t *testing.T) {
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", ""); err != nil {
		t.Errorf(err.Error())
	}

	_, err := c.Execute("set AUTOCOMMIT = 0")
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = c.Execute("set AUTOCOMMIT = 0")
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = c.Execute("set NAMES utf8")
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = c.Execute("set xxxx = 0")
	if err != nil {
		t.Error(err.Error())
		return
	}
}
