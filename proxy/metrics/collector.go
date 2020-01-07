package metrics

import (
	"git.2dfire.net/zerodb/common/system"
	"git.2dfire.net/zerodb/common/system/load"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"

	"github.com/prometheus/client_golang/prometheus"
)

type ProxyCollector struct {
	Manager *ClusterManager
}

func (c ProxyCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c ProxyCollector) Collect(ch chan<- prometheus.Metric) {
	cpu, err := system.CPULoad()
	if err != nil {
		glog.Glog.Errorln(err)
	} else {
		ch <- prometheus.MustNewConstMetric(CPUTemp, prometheus.GaugeValue, cpu)
	}

	mem, err := system.MemLoad()
	if err != nil {
		glog.Glog.Errorln(err)
	} else {
		ch <- prometheus.MustNewConstMetric(Memory, prometheus.GaugeValue, mem)
	}

	avg, err := load.LoadAvg()
	if err != nil {
		glog.Glog.Errorln(err)
	} else {
		ch <- prometheus.MustNewConstMetric(Load1avg, prometheus.GaugeValue, avg.Avg1min)
	}

	ch <- prometheus.MustNewConstMetric(TPS, prometheus.GaugeValue, float64(monitor.Monitor.GetTPS()))

	ch <- prometheus.MustNewConstMetric(QPS, prometheus.GaugeValue, float64(monitor.Monitor.GetQPS()))

	ch <- prometheus.MustNewConstMetric(FrontendConns, prometheus.GaugeValue, float64(monitor.Monitor.GetFrontendConns()))

	myData := c.Manager.Export()
	if len(myData) > 0 {
		for key, val := range myData {
			ch <- prometheus.MustNewConstMetric(ProxyMetrics, prometheus.GaugeValue, val, key)
		}
	}
}
