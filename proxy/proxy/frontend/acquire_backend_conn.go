package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"

	"github.com/pkg/errors"
)

// SELECT @@session.tx_isolation 等信息查询类的命令专用
// 如果是正常SQL，用getBackendConn
func (fc *FrontendConn) getInfoSingleBackendConn(sql string) (bc *backend.BackendConn, err error) {
	executeDB := new(ExecuteDB)
	executeDB.sql = sql
	err = fc.SetExecuteNodeForDefault(executeDB)

	if err != nil {
		return nil, err
	}

	hostGroup := fc.proxy.GetHostGroupNode(executeDB.ExecNode.HostGroup)
	if hostGroup == nil {
		return nil, fmt.Errorf(errcode.HostGroupNotExist.ErrMsg,
			errcode.HostGroupNotExist.ErrCode,
			executeDB.ExecNode.HostGroup)
	}
	bc, err = hostGroup.GetWriteConn()
	return
}

//获取shard的conn
func (fc *FrontendConn) getShardingBackendConns(rwSplit bool, readOnly bool, plan *Plan) (map[int]*backend.BackendConn, error) {
	var err error
	shardingCount := len(plan.ShardingNodes)
	if plan == nil || shardingCount == 0 {
		return nil, errcode.ErrNoShardingNode
	}

	// 如果事务是打开的，不支持读写分离
	if fc.isInTransaction() {
		if 1 < shardingCount {
			return nil, errcode.ErrMultiNodeTranNotSupport
		}
		//exec in multi node
		var theOnlyNode *backend.ShardingNode
		for _, n := range plan.ShardingNodes {
			theOnlyNode = n // 这里用遍历只是为了找出map中唯一的一个node
			break
		}
		if len(fc.txConns) == 1 && fc.txConns[theOnlyNode.TableIndex] == nil {
			return nil, errors.WithStack(errcode.ErrMultiNodeTranNotSupport)
		}
	}

	// TODO nanxing 优化方向，并发创建连接
	// 给每一个shardingNode创建一个连接
	conns := make(map[int]*backend.BackendConn)
	var co *backend.BackendConn
	for key, n := range plan.ShardingNodes {
		co, err = fc.getBackendConn(n, rwSplit, readOnly)
		if err != nil {
			break
		}

		conns[key] = co
	}

	return conns, err
}

// TODO nanxing Bug detected. max_conn如果配置得不够大，流程会被阻塞。
func (fc *FrontendConn) getBackendConn(shardingNode *backend.ShardingNode, rwSplit bool, readOnly bool) (bc *backend.BackendConn, err error) {
	hostGroup := fc.proxy.GetHostGroupNode(shardingNode.HostGroup)
	if hostGroup == nil {
		return nil, errcode.BuildError(errcode.HostGroupNotExist, shardingNode.HostGroup)
	}

	if !readOnly {
		// 从当前开始不可写入，让已经在写的连接写完。
		if !hostGroup.IsWritable() {
			return nil, errcode.BuildError(errcode.HostNotWritable, shardingNode.HostGroup)
		}
	}

	if !fc.isInTransaction() {
		if rwSplit {
			bc, err = hostGroup.GetReadConn()
			if err != nil {
				glog.Glog.Errorf("hostGroup:[%s], rwSplit is on, but read conn can't be fetched.", shardingNode.HostGroup)
				bc, err = hostGroup.GetWriteConn()
			}
		} else {
			bc, err = hostGroup.GetWriteConn()
		}
		if err != nil {
			glog.Glog.Errorf("hostGroup:[%s], Conn can't be fetched.", shardingNode.HostGroup)
			return
		}
	} else {
		var ok bool
		bc, ok = fc.txConns[shardingNode.TableIndex]

		if !ok {
			if bc, err = hostGroup.GetWriteConn(); err != nil {
				return
			}

			if !fc.isAutoCommit() {
				if err = bc.SetAutoCommit(0); err != nil {
					return
				}
			} else {
				if err = bc.Begin(); err != nil {
					return
				}
			}

			fc.txConns[shardingNode.TableIndex] = bc
		}
	}

	if err = bc.SetCharset(fc.charset, fc.collation); err != nil {
		return
	}

	return
}
