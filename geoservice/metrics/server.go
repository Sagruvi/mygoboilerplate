package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// count of requests
var RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "proxy",
	Subsystem: "server",
	Name:      "request_count",
	Help:      "Count of requests.",
}, []string{"method", "endpoint"})
var RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "proxy",
	Subsystem: "server",
	Name:      "request_duration_seconds",
	Help:      "Duration of requests.",
}, []string{"method", "endpoint"})
var CacheAccessDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "proxy",
	Subsystem: "server",
	Name:      "cache_access_duration_seconds",
	Help:      "Duration of cache access.",
}, []string{"method", "endpoint"})
var DBAccessDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "proxy",
	Subsystem: "server",
	Name:      "db_access_duration_seconds",
	Help:      "Duration of db access.",
}, []string{"method", "endpoint"})
var ExternalAPIAccessDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "proxy",
	Subsystem: "server",
	Name:      "external_api_access_duration_seconds",
	Help:      "Duration of external api access.",
}, []string{"method", "endpoint"})

func PrometheusMiddleware() {
	mux := http.NewServeMux()
	prometheus.MustRegister(RequestDuration, CacheAccessDuration, DBAccessDuration, ExternalAPIAccessDuration, RequestCount)
	mux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":19090", mux)
}
