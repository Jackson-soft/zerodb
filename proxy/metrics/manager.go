package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type ClusterManager struct {
	Zone string
	lock sync.RWMutex
	Data map[string]float64
}

func NewClusterManager(zone string, reg prometheus.Registerer) *ClusterManager {
	c := new(ClusterManager)
	c.Zone = zone
	c.Data = make(map[string]float64)

	cc := ProxyCollector{
		Manager: c,
	}
	prometheus.WrapRegistererWith(prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
	return c
}

func (cm *ClusterManager) Export() map[string]float64 {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	return cm.Data
}

func (cm *ClusterManager) Load(key string, val float64) {
	cm.lock.Lock()
	cm.Data[key] = val
	cm.lock.Unlock()
}
