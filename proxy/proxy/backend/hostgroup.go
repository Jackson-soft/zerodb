package backend

import (
	"git.2dfire.net/zerodb/common/config"
	"sync"
)

const (
	Master      = "master"
	Slave       = "slave"
	HostSplit   = ","
	WeightSplit = "@"
)

// HostGroup维度管理数据库连接池
type HostGroupNode struct {
	Cfg config.HostGroup

	sync.RWMutex

	Write             []*DBPool
	activedWriteIndex int

	writable bool // 「写」功能

	Read []*DBPool

	// load balance for read
	LastReadIndex int
	RoundRobinQ   []int
	ReadWeights   []int
}

func (node *HostGroupNode) InitPool(addr string, poolType string) (*DBPool, error) {
	// 初始化的连接没有默认的schema的
	dbPool, err := InitDBPool(addr, node.Cfg.Name, poolType, node.Cfg.User, node.Cfg.Password, "", node.Cfg.InitConn, node.Cfg.MaxConn)
	return dbPool, err
}

func (node *HostGroupNode) String() string {
	return node.Cfg.Name
}
