package backend

import (
	"errors"
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
)

// 数据源切换
func (node *HostGroupNode) SwitchDatasource(toIndex int) (string, error) {
	if node.activedWriteIndex == toIndex {
		// 不需要切换
		glog.Glog.Infof("hostGroup:%s, no need switching datasource from %v to %v.", node.Cfg.Name, node.activedWriteIndex, toIndex)
		return fmt.Sprintf("hostGroup:%s is on %d, no need switching.", node.Cfg.Name, toIndex), nil
	}

	fromIndex := node.activedWriteIndex

	newDbPool := node.Write[toIndex]
	if newDbPool == nil {
		return "", errors.New("writeDBPool for index:[" + string(toIndex) + "] doesn't exist.")
	}

	// 互斥锁，这时不能获取连接
	node.Lock()
	defer node.Unlock()

	oldDbPool := node.Write[node.activedWriteIndex]

	// 在关闭之前，需要将数据源设置为「关闭」状态，防止数据源在「CloseActiveBackendConns」的时候连接被借走。
	// 将旧的数据源的连接全部关闭掉（让proxy向mysql发送quit命令），否则数据源切换后，事务将不能被完整执行
	oldDbPool.CloseActiveBackendConns()

	// 设置心跳检测关闭，然后将用于心跳的detectConn关闭
	oldDbPool.NeedHeartbeat = false

	// 设置心跳检测开启，然后将用于心跳的detectConn打开
	newDbPool.NeedHeartbeat = true

	node.activedWriteIndex = toIndex

	glog.Glog.Infof("hostGroup:%s, switching datasource from %v to %v.", node.Cfg.Name, fromIndex, toIndex)
	return fmt.Sprintf("hostGroup:%s from %v to %v.", node.Cfg.Name, fromIndex, toIndex), nil
}

func (node *HostGroupNode) Heartbeat() error {
	dbPool := node.Write[node.activedWriteIndex]
	if !dbPool.NeedHeartbeat {
		return nil
	}

	return dbPool.Ping()
}

func (node *HostGroupNode) StopWriting() {
	node.writable = false
}

func (node *HostGroupNode) EnableWriting() {
	node.writable = true
}

func (node *HostGroupNode) IsWritable() bool {
	return node.writable
}

func (node *HostGroupNode) GetNextWriteIndex() int {
	return (node.activedWriteIndex + 1) % len(node.Write)
}

func (node *HostGroupNode) GetActivedWriteIndex() int {
	return node.activedWriteIndex
}

func (node *HostGroupNode) InitActivedWriteIndex(activedWriteIndex int) {
	node.activedWriteIndex = activedWriteIndex
}
