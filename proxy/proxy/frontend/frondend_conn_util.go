package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"sync"
)

func (fc *FrontendConn) readPacket() ([]byte, error) {
	return fc.pkg.ReadPacket()
}

func (fc *FrontendConn) writePacket(data []byte) error {
	return fc.pkg.WritePacket(data)
}

func (fc *FrontendConn) writePacketBatch(total, data []byte, direct bool) ([]byte, error) {
	return fc.pkg.WritePacketBatch(total, data, direct)
}

func (fc *FrontendConn) executeInNode(conn *backend.BackendConn, sql string, args []interface{}) ([]*mysql.Result, error) {
	r, err := conn.Execute(sql, args...)
	if err != nil {
	} else {
	}

	if err != nil {
		return nil, err
	}

	return []*mysql.Result{r}, err
}

func (fc *FrontendConn) executeInMultiNodes(conns map[int]*backend.BackendConn, shardingNodes map[int]*backend.ShardingNode, args []interface{}, logSQL bool) ([]*mysql.Result, error) {
	if len(conns) != len(shardingNodes) {
		return nil, fmt.Errorf("length of conn and shardingNode should be 1 : 1, now is %d : %d", len(conns), len(shardingNodes))
	}

	var wg sync.WaitGroup

	if len(conns) == 0 {
		return nil, errcode.ErrNoConn
	}

	wg.Add(len(conns))

	resultCount := len(shardingNodes)

	rs := make([]interface{}, resultCount)

	f := func(rs []interface{}, offset int, shardingNode *backend.ShardingNode, backendConn *backend.BackendConn, alloc *util.FastYetSafeAllocator) {
		// TODO nanxing 超时机制，否则对于大SQL，会占用太多时间和连接资源
		r, err := backendConn.ExecuteWithAlloc(alloc, shardingNode.ShardingSQL, args...)
		if logSQL {
			glog.Glog.Infof("DBPool: %s, SQL: %s", backendConn.GetDBPool().GetName(), shardingNode.ShardingSQL)
		}
		if err != nil {
			rs[offset] = err
		} else {
			rs[offset] = r
		}

		wg.Done()
	}

	offsert := 0
	for shardingIndex, co := range conns {
		shardingNode := shardingNodes[shardingIndex]
		go f(rs, offsert, shardingNode, co, fc.executeArena)
		offsert += 1
	}

	wg.Wait()

	var err error
	results := make([]*mysql.Result, resultCount)
	for i, v := range rs {
		if e, ok := v.(error); ok {
			err = e
			break
		}
		if rs[i] != nil {
			results[i] = rs[i].(*mysql.Result)
		}
	}

	return results, err
}
