package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rateLimitHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(rateLimitHits)
}

func SetupMetrics() {
	http.Handle("/metrics", promhttp.Handler())
}

func RecordRateLimitHit(status string) {
	rateLimitHits.WithLabelValues(status).Inc()
}
