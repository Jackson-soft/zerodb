package util

import (
	"sync"

	"git.2dfire.net/zerodb/proxy/pkg/glog"
)

// Allocator pre-allocates memory to reduce memory allocation cost.
type Allocator interface {
	// Alloc allocates memory with 0 len and capacity cap.
	Alloc(capacity int) []byte

	// AllocWithLen allocates memory with length and capacity.
	AllocWithLen(length int, capacity int) []byte

	// Reset resets arena offset.
	// Make sure all the allocated memory are not used any more.
	Reset()
}

// SimpleAllocator is a simple implementation of ArenaAllocator.
// It is not goroutine-safe.
type FastAllocator struct {
	arena []byte
	off   int
}

// NewFastAllocator creates an Allocator with a specified capacity.
func NewFastAllocator(capacity int) *FastAllocator {
	return &FastAllocator{arena: make([]byte, 0, capacity)}
}

// Alloc implements Allocator.AllocBytes interface.
func (s *FastAllocator) Alloc(capacity int) []byte {
	if s.off+capacity < cap(s.arena) {
		slice := s.arena[s.off : s.off : s.off+capacity]
		s.off += capacity
		return slice
	}

	glog.Glog.Warnf("insufficient memory in arena, better use a bigger capacity to init arena.")
	return make([]byte, 0, capacity)
}

// AllocWithLen implements Allocator.AllocWithLen interface.
func (s *FastAllocator) AllocWithLen(length int, capacity int) []byte {
	slice := s.Alloc(capacity)
	return slice[:length:capacity]
}

// Reset implements Allocator.Reset interface.
func (s *FastAllocator) Reset() {
	s.off = 0
}

// 协程安全
type FastYetSafeAllocator struct {
	arena []byte
	off   int
	sync.Mutex
}

func NewFastYetSafeAllocator(capacity int) *FastYetSafeAllocator {
	return &FastYetSafeAllocator{arena: make([]byte, 0, capacity)}
}

// 不要直接用这个，没有上锁
func (c *FastYetSafeAllocator) Alloc(capacity int) []byte {
	if c.off+capacity < cap(c.arena) {
		slice := c.arena[c.off : c.off : c.off+capacity]
		c.off += capacity
		return slice
	}

	glog.Glog.Warnf("insufficient memory in arena, better use a bigger capacity to init arena.")
	return make([]byte, 0, capacity)
}

// AllocWithLen implements Allocator.AllocWithLen interface.
func (c *FastYetSafeAllocator) AllocWithLen(length int, capacity int) []byte {
	defer c.Unlock()

	c.Lock()
	slice := c.Alloc(capacity)
	return slice[:length:capacity]
}

// Reset implements Allocator.Reset interface.
func (c *FastYetSafeAllocator) Reset() {
	defer c.Unlock()
	c.Lock()
	c.off = 0
}
