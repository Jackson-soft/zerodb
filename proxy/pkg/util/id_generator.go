package util

import "sync/atomic"

type IdGenerator struct {
	currentId int64
}

func (g *IdGenerator) NextId() int64 {
	return atomic.AddInt64(&g.currentId, 1)
}
