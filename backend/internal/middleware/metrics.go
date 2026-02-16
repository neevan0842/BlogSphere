package middleware

import (
	"fmt"
	"net/http"

	"github.com/neevan0842/BlogSphere/backend/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *responseRecorder) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		statusStr := fmt.Sprintf("%d", recorder.status)

		// record duration and status after handler completes
		timer := prometheus.NewTimer(
			utils.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, statusStr),
		)
		timer.ObserveDuration()

		// increment counter with status
		utils.HttpRequestsTotal.
			WithLabelValues(r.Method, r.URL.Path, statusStr).
			Inc()
	})
}
