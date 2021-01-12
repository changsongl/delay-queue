package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

// server prometheus metrics
var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "delay_queue_http_request_duration_seconds",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: []float64{.05, 0.1, .25, .5, .75, 1, 2, 5, 20, 60},
		},
		[]string{"handler", "method"},
	)

	requestStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "delay_queue_http_request_status",
			Help: "Counter of requests' HTTP code.",
		},
		[]string{"handler", "method", "code"},
	)
)

// init all metrics to prometheus default register
func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestStatus)
}

// setServerMetricHandlerAndMiddleware return a function to set
// metrics api and a middleware to save http request statistics.
func setServerMetricHandlerAndMiddleware() func(r *gin.Engine) {
	return func(r *gin.Engine) {
		r.GET("/metrics", func(c *gin.Context) {
			promhttp.Handler().ServeHTTP(c.Writer, c.Request)
		})
		r.Use(saveHttpServerStat)
	}
}

// saveHttpServerStat save http server request stats to metrics.
func saveHttpServerStat(c *gin.Context) {
	startTime := time.Now()
	url, method := c.Request.URL.Path, c.Request.Method
	c.Next()

	requestDuration.With(map[string]string{
		"handler": url,
		"method":  method,
	}).Observe(time.Since(startTime).Seconds())

	status := c.Writer.Status()
	if status != http.StatusNotFound {
		requestStatus.With(map[string]string{
			"handler": url,
			"method":  method,
			"code":    strconv.Itoa(c.Writer.Status()),
		}).Inc()
	}
}
