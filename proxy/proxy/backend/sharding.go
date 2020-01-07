package backend

import "fmt"

const (
	ShardingTypeTable  = 0
	ShardingTypeSchema = 1
)

type ShardingNode struct {
	SchemaIndex int
	TableIndex  int
	// 当ShardingType为TABLE时，只读取SchemaIndex的值
	// 当ShardingType为SCHEMA时，要读取SchemaIndex和TableIndex的值
	ShardingType int

	ShardingSQL string // 根据分库分表信息，改写后的SQL
	HostGroup   string
}

func NewShardingNode(schemaIndex int, tableIndex int, shardingType int, sql string, hostGroup string) *ShardingNode {

	shardingNode := new(ShardingNode)
	shardingNode.HostGroup = hostGroup
	shardingNode.SchemaIndex = schemaIndex
	shardingNode.TableIndex = tableIndex
	shardingNode.ShardingType = shardingType
	shardingNode.ShardingSQL = sql

	return shardingNode
}

func (n ShardingNode) String() string {
	return fmt.Sprintf("{SchemaIndex=%d, TableIndex=%d, HostGroup=%s}", n.SchemaIndex, n.TableIndex, n.HostGroup)
}
