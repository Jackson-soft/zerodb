package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	//
	ConnGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: zone,
			Subsystem: "server",
			Name:      "connections",
			Help:      "Number of connections.",
		})

	// agent存活数
	AgentGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: zone,
		Subsystem: "server",
		Name:      "agent_number",
		Help:      "Number of agent alive.",
	})

	// proxy 存活数
	ProxyGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: zone,
		Subsystem: "server",
		Name:      "proxy_number",
		Help:      "Number of proxy alive.",
	})
)

func init() {
	reg.MustRegister(ConnGauge, AgentGauge, ProxyGauge)
}
