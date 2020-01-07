package metrics

import (
	"net/http"

	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reg       = prometheus.NewPedanticRegistry()
	namespace = "zero_proxy"

	Manager *ClusterManager
)

func Run(address string) {
	Manager = NewClusterManager(namespace, reg)

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(address, nil); err != nil {
		glog.Glog.Fatalln(err)
	}
}
