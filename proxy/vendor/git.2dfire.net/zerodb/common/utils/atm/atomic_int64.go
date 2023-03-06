package atm

import "sync/atomic"

type AtomicInt64 struct {
	num int64
}

func NewAtomicInt64() *AtomicInt64 {
	ai := new(AtomicInt64)
	ai.num = 1
	return ai
}

func (i *AtomicInt64) Add(n int64) {
	i.num = atomic.AddInt64(&i.num, n)
}

func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&i.num)
}

func (i *AtomicInt64) CompareAndSwap(oldval, newval int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&i.num, oldval, newval)
}
