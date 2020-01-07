package frontend

import (
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"testing"
)

func TestTxShardingTables(t *testing.T) {
	// Rollback ###########################################################
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", "zerodb"); err != nil {
		t.Errorf(err.Error())
	}
	// clear data
	clearData(c)
	c.Execute("set autocommit = 0")
	tx(c, t)
	c.Rollback()

	if r, err := c.Execute("select * from order_ins where entity_id = 123125"); err != nil {
		t.Errorf(err.Error())
	} else {
		if len(r.Values) != 0 {
			t.Errorf("Result rows should be 0")
		}
	}
	if r, err := c.Execute("select * from ins_detail where entity_id = 123125"); err != nil {
		t.Errorf(err.Error())
	} else {
		if len(r.Values) != 0 {
			t.Errorf("Result rows should be 0")
		}
	}
	c.Close()

	// Commit ###########################################################
	c = backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", "zerodb"); err != nil {
		t.Errorf(err.Error())
	}
	// clear data
	clearData(c)
	c.Execute("set autocommit = 0")
	tx(c, t)
	c.Commit()
	if r, err := c.Execute("select * from order_ins where entity_id = 123125"); err != nil {
		t.Errorf(err.Error())
	} else {
		if len(r.Values) != 1 {
			t.Errorf("Result rows should be 1")
		}
	}
	if r, err := c.Execute("select * from ins_detail where entity_id = 123125"); err != nil {
		t.Errorf(err.Error())
	} else {
		if len(r.Values) != 1 {
			t.Errorf("Result rows should be 1")
		}
	}
	c.Close()
}

func TestTxNonshardingTables(*testing.T) {
	// commit
	// rollback
}

func TestTxCustodyTables(*testing.T) {
	// commit
	// rollback
}

func clearData(c *backend.MysqlConn) {
	c.Execute("delete from order_ins where entity_id = 123125")
	c.Execute("delete from ins_detail where entity_id = 123125")
}

func tx(c *backend.MysqlConn, t *testing.T) {
	if r, err := c.Execute("INSERT INTO order_ins (id, entity_id, name) VALUES (1, 123125, 'hello world🌀order_ins')"); err != nil {
		t.Errorf(err.Error())
	} else {
		if r.AffectedRows != 1 {
			t.Errorf("AffectedRows should be 1")
		}
	}
	if r, err := c.Execute("select * from order_ins where entity_id = 123125"); err != nil {
		t.Errorf(err.Error())
	} else {
		if len(r.Values) != 1 {
			t.Errorf("Result rows should be 1")
		}
	}
	if r, err := c.Execute("INSERT INTO ins_detail (id, entity_id, detail) VALUES (1, 123125, 'hello world👿')"); err != nil {
		t.Errorf(err.Error())
	} else {
		if r.AffectedRows != 1 {
			t.Errorf("AffectedRows should be 1")
		}
	}
	if r, err := c.Execute("select * from ins_detail where entity_id = 123125"); err != nil {
		t.Errorf(err.Error())
	} else {
		if len(r.Values) != 1 {
			t.Errorf("Result rows should be 1")
		}
	}
}
