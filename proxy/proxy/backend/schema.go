package backend

import (
	"git.2dfire.net/zerodb/common/config"
	"sync"
)

type SchemaNode struct {
	Name string

	Conf *config.SchemaConfig

	// 部分分库分表
	HostGroupCluster string

	ShardingHostGroupNodes   map[string]*HostGroupNode
	NonshardingHostGroupNode *HostGroupNode

	TableNodes sync.Map //map[string]*TableNode

	// other settings
	RwSplit                     bool
	MultiRoutePermitted         bool
	InitConnMultiRoutePermitted bool
}

func (s *SchemaNode) GetTableNode(name string) *TableNode {
	result, ok := s.TableNodes.Load(name)
	if ok {
		return result.(*TableNode)
	} else {
		return nil
	}
}
