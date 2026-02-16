package utils

import "github.com/prometheus/client_golang/prometheus"

var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed, labeled by method, path and status code.",
	},
	[]string{"method", "path", "status"},
)

var HttpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of response latency (seconds) for HTTP requests.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)
