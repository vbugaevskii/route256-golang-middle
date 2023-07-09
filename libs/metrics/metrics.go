package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MetricRequestCounter              prometheus.Counter
	MetricResponseCounter             *prometheus.CounterVec
	MetricResponseTimeHistogram       *prometheus.HistogramVec
	MetricClientResponseTimeHistogram *prometheus.HistogramVec
)

func Init(serviceName string) {
	MetricRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: serviceName,
		Name:      "server_request_counter",
	})

	MetricResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: serviceName,
		Name:      "server_response_counter",
	},
		[]string{"status"},
	)

	MetricResponseTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: serviceName,
		Name:      "server_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status"},
	)

	MetricClientResponseTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: serviceName,
		Name:      "client_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status"},
	)
}

func ListenAndServeMetrics(port int) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
