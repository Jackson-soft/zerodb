package frontend

import (
	//"fmt"
	//"git.2dfire.net/zerodb/proxy/proxy/backend"
	"strconv"
	"testing"
)

func TestIntShard(t *testing.T) {
	shardingTableCount := 128

	intHashShard := &IntHashShard{SlotNum: 1024, ShardLength: 1024 / shardingTableCount}

	if shard, err := intHashShard.FindForKey(int64(0)); err != nil {
		t.Fatal(err)
	} else {
		if shard != 0 {
			t.Fatal("Sharding err")
		}
	}

	if shard, err := intHashShard.FindForKey(int64(-1)); err != nil {
		t.Fatal(err)
	} else {
		if shard != 127 {
			t.Fatal("Sharding err")
		}
	}

	if shard, err := intHashShard.FindForKey(int64(1023)); err != nil {
		t.Fatal(err)
	} else {
		if shard != 127 {
			t.Fatal("Sharding err")
		}
	}

	if shard, err := intHashShard.FindForKey(int64(1024)); err != nil {
		t.Fatal(err)
	} else {
		if shard != 0 {
			t.Fatal("Sharding err")
		}
	}

	i, err := strconv.ParseInt("00000976", 10, 64)
	if err != nil {
		t.Fatal("ParseInt err")
	} else {
		if shard, err := intHashShard.FindForKey(i); err != nil {
			t.Fatal(err)
		} else {
			if shard != 122 {
				t.Fatal("Sharding err")
			}
		}
	}

}

/*
func TestCobarIntSharding(t *testing.T) {
	cobarConn := newCobarConn()
	zerodbConn := newZerodbConn()

	// int entity_id rule
	cobarSql := fmt.Sprintf("select instance_id, entity_id from instancedetail limit %d", 10)
	if result, err := cobarConn.Execute(cobarSql); err != nil {
		t.Fatal(err)
	} else {
		for _, rd := range result.Resultset.RowDatas {
			a, _ := rd.ParseText(result.Fields)
			id := a[0].(string)
			entity := a[1].(string)

			sql := fmt.Sprintf("select instance_id, entity_id from instancedetail where instance_id = '%s' and entity_id = '%s'", id, entity)

			result, err = zerodbConn.Execute(sql)
			if err != nil {
				t.Fatal(err)
			} else {
				if len(result.Values) != 1 {
					t.Fatal("Different behaviour between cobar and zeroproxy")
				}
			}
		}
	}
}

func TestCobarStringSharding(t *testing.T) {
	cobarConn := newCobarConn()
	zerodbConn := newZerodbConn()

	// int entity_id rule
	cobarSql := fmt.Sprintf("select waitingorder_id, customerregister_id from waitingordercrid limit %d", 10)
	if result, err := cobarConn.Execute(cobarSql); err != nil {
		t.Fatal(err)
	} else {
		for _, rd := range result.Resultset.RowDatas {
			a, _ := rd.ParseText(result.Fields)
			id := a[0].(string)
			customerregister := a[1].(string)

			sql := fmt.Sprintf("select waitingorder_id, customerregister_id from waitingordercrid where waitingorder_id = '%s' and customerregister_id = '%s'", id, customerregister)

			result, err = zerodbConn.Execute(sql)
			if err != nil {
				t.Fatal(err)
			} else {
				if len(result.Values) != 1 {
					t.Fatal("Different behaviour between cobar and zeroproxy")
				}
			}
		}
	}
}

func newCobarConn() *backend.MysqlConn {
	c := backend.NewConnCache()
	if err := c.ConnectMysql("cobar90:8066", "order", "order@552208", "order"); err != nil {
		panic(err)
	}
	c.Execute("set MULTI_ROUTE_PERMIT = 1")
	return c
}

func newZerodbConn() *backend.MysqlConn {
	c := backend.NewConnCache()
	if err := c.ConnectMysql("127.0.0.1:9696", "zerodb", "zerodb", "order"); err != nil {
		panic(err)
	}
	return c
}
*/
