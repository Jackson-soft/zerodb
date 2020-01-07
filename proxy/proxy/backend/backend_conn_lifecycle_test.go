package backend

import (
	"fmt"
	"testing"

	"git.2dfire.net/zerodb/proxy/proxy/mysql"
)

func newTestConn() *MysqlConn {
	c := NewConnCache()

	if err := c.ConnectMysql("10.1.22.1:3306", "zerodb", "zerodb@552208", ""); err != nil {
		panic(err)
	}

	return c
}

func TestConn_Connect(t *testing.T) {
	c := newTestConn()
	defer c.Close()
}

func TestConn_Ping(t *testing.T) {
	c := newTestConn()
	defer c.Close()

	if err := c.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestConn_CreateDatabase(t *testing.T) {
	c := newTestConn()
	defer c.Close()

	if _, err := c.Execute("create database if not exists backend_conn_test_database"); err != nil {
		t.Fatal(err)
	}
}

func TestConn_DeleteTable(t *testing.T) {
	c := newTestConn()
	defer c.Close()

	c.UseDB("backend_conn_test_database")

	if _, err := c.Execute("drop table if exists backend_conn_test_table"); err != nil {
		t.Fatal(err)
	}
}

func TestConn_CreateTable(t *testing.T) {
	s := `CREATE TABLE IF NOT EXISTS backend_conn_test_table (
          id BIGINT(64) UNSIGNED  NOT NULL,
          str VARCHAR(256),
          f DOUBLE,
          e enum("test1", "test2"),
          u tinyint unsigned,
          i tinyint,
          PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8`

	c := newTestConn()
	defer c.Close()
	c.UseDB("backend_conn_test_database")

	if _, err := c.Execute(s); err != nil {
		t.Fatal(err)
	}
}

func TestConn_Insert(t *testing.T) {
	s := `insert into backend_conn_test_table (id, str, f, e) values(1, "a", 3.14, "test1")`

	c := newTestConn()
	defer c.Close()
	c.UseDB("backend_conn_test_database")

	if pkg, err := c.Execute(s); err != nil {
		t.Fatal(err)
	} else {
		if pkg.AffectedRows != 1 {
			t.Fatal(pkg.AffectedRows)
		}
	}
}

func TestConn_Select(t *testing.T) {
	s := `select str, f, e from backend_conn_test_table where id = 1`

	c := newTestConn()
	defer c.Close()

	c.UseDB("backend_conn_test_database")

	if result, err := c.Execute(s); err != nil {
		t.Fatal(err)
	} else {
		if len(result.Fields) != 3 {
			t.Fatal(len(result.Fields))
		}

		if len(result.Values) != 1 {
			t.Fatal(len(result.Values))
		}

		if str, _ := result.GetString(0, 0); str != "a" {
			t.Fatal("invalid str", str)
		}

		if f, _ := result.GetFloat(0, 1); f != float64(3.14) {
			t.Fatal("invalid f", f)
		}

		if e, _ := result.GetString(0, 2); e != "test1" {
			t.Fatal("invalid e", e)
		}

		if str, _ := result.GetStringByName(0, "str"); str != "a" {
			t.Fatal("invalid str", str)
		}

		if f, _ := result.GetFloatByName(0, "f"); f != float64(3.14) {
			t.Fatal("invalid f", f)
		}

		if e, _ := result.GetStringByName(0, "e"); e != "test1" {
			t.Fatal("invalid e", e)
		}

	}
}

func TestConn_Escape(t *testing.T) {
	c := newTestConn()
	defer c.Close()

	e := `""''\abc`
	s := fmt.Sprintf(`insert into backend_conn_test_table (id, str) values(5, "%s")`,
		mysql.Escape(e))

	c.UseDB("backend_conn_test_database")

	if _, err := c.Execute(s); err != nil {
		t.Fatal(err)
	}

	s = `select str from backend_conn_test_table where id = ?`

	if r, err := c.Execute(s, 5); err != nil {
		t.Fatal(err)
	} else {
		str, _ := r.GetString(0, 0)
		if str != e {
			t.Fatal(str)
		}
	}
}

func TestConn_SetCharset(t *testing.T) {
	c := newTestConn()
	defer c.Close()

	if err := c.SetCharset("gb2312", 24); err != nil {
		t.Fatal(err)
	}
}

func TestConn_SetAutoCommit(t *testing.T) {
	c := newTestConn()
	defer c.Close()

	if err := c.SetAutoCommit(0); err != nil {
		t.Fatal(err)
	}

	if err := c.SetAutoCommit(1); err != nil {
		t.Fatal(err)
	}
}
