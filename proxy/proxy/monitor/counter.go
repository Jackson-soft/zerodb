package monitor

import (
	"git.2dfire.net/zerodb/common/utils/atm"
	"sync"
	"sync/atomic"
	"time"
)

type Counter struct {
	sync.RWMutex

	clientQPS       int64 // 瞬时值
	oneSecClientQPS int64 // 上1秒QPS

	clientTPS       int64 // 瞬时值
	oneSecClientTPS int64 // 上1秒TPS

	SlowLogTotal    int64 // 瞬时值
	OldSlowLogTotal int64 // 上1秒SlowQuery

	frontendConns           int64
	schemaedFrontendConns   sync.Map // map[string]AtomicInt64
	hostGroupedBackendConns sync.Map // map[string]AtomicInt64

	schemaedTPS sync.Map // map[string]AtomicInt64
	schemaedQPS sync.Map // map[string]AtomicInt64
}

func newCounter() *Counter {
	counter := new(Counter)
	go counter.flush()
	return counter
}

func (counter *Counter) GetFrontendConns() int64 {
	return counter.frontendConns
}

func (counter *Counter) GetTPS() int64 {
	return counter.oneSecClientTPS
}

func (counter *Counter) GetQPS() int64 {
	return counter.oneSecClientQPS
}

func (counter *Counter) IncrClientConns() {
	atomic.AddInt64(&counter.frontendConns, 1)
}

func (counter *Counter) DecrClientConns() {
	atomic.AddInt64(&counter.frontendConns, -1)
}

func (counter *Counter) IncrFrontendConns(schema string) {
	result, ok := counter.schemaedFrontendConns.Load(schema)
	if ok {
		c := result.(*atm.AtomicInt64)
		c.Add(1)
	} else {
		counter.Lock()
		defer counter.Unlock()
		result, ok = counter.schemaedFrontendConns.Load(schema)
		if !ok {
			c := atm.NewAtomicInt64()
			counter.schemaedFrontendConns.Store(schema, c)
		}
	}
}

func (counter *Counter) IncrBackendConns(schema string) {
	result, ok := counter.schemaedFrontendConns.Load(schema)
	if ok {
		c := result.(*atm.AtomicInt64)
		c.Add(1)
	} else {
		counter.Lock()
		defer counter.Unlock()
		result, ok = counter.schemaedFrontendConns.Load(schema)
		if !ok {
			c := atm.NewAtomicInt64()
			counter.schemaedFrontendConns.Store(schema, c)
		}
	}
}

func (counter *Counter) DecrFrontendConns(schema string) {
	result, ok := counter.schemaedFrontendConns.Load(schema)
	if ok {
		c := result.(*atm.AtomicInt64)
		c.Add(-1)
	} else {
		// 有问题，没有创建连接，怎么降低。
	}
}

func (counter *Counter) IncrClientQPS(schema string) {
	atomic.AddInt64(&counter.clientQPS, 1)
	counter.IncrClientOPS()

}

func (counter *Counter) IncrClientTPS(schema string) {
	atomic.AddInt64(&counter.clientTPS, 1)
	counter.IncrClientOPS()

	result, ok := counter.schemaedTPS.Load(schema)
	if ok {
		c := result.(*atm.AtomicInt64)
		c.Add(1)
	} else {
		counter.Lock()
		defer counter.Unlock()
		result, ok = counter.schemaedTPS.Load(schema)
		if !ok {
			c := atm.NewAtomicInt64()
			counter.schemaedTPS.Store(schema, c)
		}
	}
}

func (counter *Counter) IncrClientOPS() {
	atomic.AddInt64(&counter.clientTPS, 1)
}

func (counter *Counter) IncrSlowLogTotal() {
	atomic.AddInt64(&counter.SlowLogTotal, 1)
}

//flush the count per second
func (counter *Counter) FlushCounter() {
	atomic.StoreInt64(&counter.oneSecClientQPS, counter.clientQPS)
	atomic.StoreInt64(&counter.oneSecClientTPS, counter.clientTPS)
	atomic.StoreInt64(&counter.OldSlowLogTotal, counter.SlowLogTotal)

	atomic.StoreInt64(&counter.clientQPS, 0)
	atomic.StoreInt64(&counter.clientTPS, 0)
}

func (counter *Counter) flush() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			counter.FlushCounter()
		}
	}
}
