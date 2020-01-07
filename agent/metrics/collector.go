package metrics

import (
	"git.2dfire.net/zerodb/agent/glog"
	"git.2dfire.net/zerodb/common/system"
	"git.2dfire.net/zerodb/common/system/load"
	"github.com/prometheus/client_golang/prometheus"
)

type AgentCollector struct {
	Zone      string
	CpuTemp   *prometheus.Desc
	Memory    *prometheus.Desc
	Load1avg  *prometheus.Desc
	Load5avg  *prometheus.Desc
	Load15avg *prometheus.Desc
}

func NewAgentCollector(zone string) *AgentCollector {
	a := new(AgentCollector)

	a.Zone = zone
	a.CpuTemp = prometheus.NewDesc(
		zone+"_cpu",
		"Current utilization of the CPU.",
		nil, nil,
	)

	a.Memory = prometheus.NewDesc(
		zone+"_memory",
		"Current utilization of the memory.",
		nil, nil,
	)

	a.Load1avg = prometheus.NewDesc(zone+"_load1", "1m load average.", nil, nil)

	a.Load5avg = prometheus.NewDesc(zone+"_load5", "5m load average.", nil, nil)

	a.Load15avg = prometheus.NewDesc(zone+"_load15", "15m load average.", nil, nil)
	return a
}

func (c *AgentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CpuTemp
	ch <- c.Memory
	ch <- c.Load1avg
	ch <- c.Load5avg
	ch <- c.Load15avg
}

func (c *AgentCollector) Collect(ch chan<- prometheus.Metric) {
	cpu, err := system.CPULoad()
	if err != nil {
		glog.Glog.Errorln(err)
	}
	ch <- prometheus.MustNewConstMetric(c.CpuTemp, prometheus.GaugeValue, cpu)

	mem, err := system.MemLoad()
	if err != nil {
		glog.Glog.Errorln(err)
	}

	ch <- prometheus.MustNewConstMetric(c.Memory, prometheus.GaugeValue, mem)

	load, err := load.LoadAvg()
	if err != nil {
		glog.Glog.Errorln(err)
	}

	ch <- prometheus.MustNewConstMetric(c.Load1avg, prometheus.GaugeValue, load.Avg1min)

	ch <- prometheus.MustNewConstMetric(c.Load5avg, prometheus.GaugeValue, load.Avg5min)

	ch <- prometheus.MustNewConstMetric(c.Load15avg, prometheus.GaugeValue, load.Avg15min)
}
