package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"testing"
)

func TestPreshard(t *testing.T) {
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", ""); err != nil {
		t.Errorf(err.Error())
	}

	_, err := c.Execute("select database()")
	if err != nil {
		t.Error(err.Error())
		return
	}

}
