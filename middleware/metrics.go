package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

var (
	RequestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "campaign_request_total",
			Help: "Total number of requests to /campaign",
		},
	)

	MatchSuccessCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "campaign_match_success_total",
			Help: "Total number of successful campaign matches",
		},
	)

	MatchMissCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "campaign_match_miss_total",
			Help: "Total number of campaign matches miss",
		},
	)

	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "campaign_request_duration_seconds",
			Help:    "Histogram of request durations for /campaign",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(MatchSuccessCount)
	prometheus.MustRegister(MatchMissCount)
	prometheus.MustRegister(RequestDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		path := c.Request.URL.Path
		log.Println("📊 Prometheus middleware path:", path)

		if path == "/campaign" {
			log.Println("📈 Incrementing campaign_requests_total")
			RequestCount.Inc()
			RequestDuration.Observe(duration)
			LogMetricValue("campaign_requests_count", RequestCount)
		}

	}
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func LogMetricValue(name string, c prometheus.Counter) {
	metric := &dto.Metric{}
	err := c.Write(metric)

	if err != nil {
		log.Printf("⚠️ Could not read %s metric: %v", name, err)
		return
	}

	log.Printf("📊 %s = %v", name, metric.GetCounter().GetValue())
}
