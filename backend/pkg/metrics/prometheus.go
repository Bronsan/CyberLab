package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cyberlab_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cyberlab_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ContainersRunning = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cyberlab_containers_running",
			Help: "Number of running Docker containers",
		},
	)

	UsersRegistered = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cyberlab_users_registered",
			Help: "Total number of registered users",
		},
	)

	ChallengesTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cyberlab_challenges_total",
			Help: "Total number of challenges",
		},
	)

	SubmissionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cyberlab_submissions_total",
			Help: "Total number of flag submissions",
		},
	)

	FlagCorrectTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cyberlab_flag_correct_total",
			Help: "Total number of correct flag submissions",
		},
	)
)

func init() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
		ContainersRunning,
		UsersRegistered,
		ChallengesTotal,
		SubmissionsTotal,
		FlagCorrectTotal,
	)
}

// Middleware returns a Gin middleware that records HTTP metrics.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			HttpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(v)
		}))
		defer timer.ObserveDuration()

		c.Next()

		status := c.Writer.Status()
		HttpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), httpStatusLabel(status)).Inc()
	}
}

func httpStatusLabel(code int) string {
	switch {
	case code < 200: return "1xx"
	case code < 300: return "2xx"
	case code < 400: return "3xx"
	case code < 500: return "4xx"
	default: return "5xx"
	}
}

// Handler returns the /metrics HTTP handler.
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
