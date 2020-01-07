package metrics

import (
	"net/http"

	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reg  = prometheus.NewPedanticRegistry()
	zone = "zero_keeper"
)

func Run(address string) error {
	zc := NewZeroCollector(zone)

	reg.MustRegister(
		zc,
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	glog.GLog.Infof("prometheus listen %v", address)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return http.ListenAndServe(address, nil)
}
