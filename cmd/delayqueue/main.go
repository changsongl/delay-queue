package main

import "github.com/prometheus/client_golang/prometheus"

// 普罗米修斯指标
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
			Name:    "delay_queue_http_request_status",
			Help:    "Counter of requests' HTTP code.",
		},
		[]string{"handler", "method", "code"},
	)
)


func init(){
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestStatus)
}

func main(){

}