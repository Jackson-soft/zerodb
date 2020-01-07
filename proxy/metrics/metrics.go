package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	QueryDurationHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      "handle_query_duration_seconds",
			Namespace: namespace,
			Subsystem: "server",
			Help:      "Bucketed histogram of processing time (s) of handled queries.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 22),
		},
	)
	QueryTotalCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "query_total",
			Help:      "Counter of queries.",
			Subsystem: "server",
		},
	)
	ConnGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "server",
			Name:      "connections",
			Help:      "Number of connections.",
		},
	)

	CPUTemp = prometheus.NewDesc(
		namespace+"_cpu",
		"Current utilization of the CPU.",
		nil, nil,
	)

	Memory = prometheus.NewDesc(
		namespace+"_memory",
		"Current utilization of the memory.",
		nil, nil,
	)

	Load1avg = prometheus.NewDesc(namespace+"_load1", "1m load average.", nil, nil)

	TPS = prometheus.NewDesc(namespace+"_tps", "Transaction per second.", nil, nil)

	QPS = prometheus.NewDesc(namespace+"_qps", "Query per second.", nil, nil)

	FrontendConns = prometheus.NewDesc(namespace+"_frontend_conns", "Frontend connections.", nil, nil)

	ProxyMetrics = prometheus.NewDesc(namespace+"_metrics", "proxy some metrics.", []string{"project"}, nil)
)

func init() {
	reg.MustRegister(
		QueryDurationHistogram,
		QueryTotalCounter,
		ConnGauge,
	)
}
