package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reg  = prometheus.NewPedanticRegistry()
	zone = "zero_agent"
)

func Run(address string) error {
	agent := NewAgentCollector(zone)

	reg.MustRegister(agent)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return http.ListenAndServe(address, nil)
}
