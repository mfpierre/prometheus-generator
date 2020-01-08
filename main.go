package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr                        = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	gauges                      = flag.Int("gauges", 100, "Number of gauges to generate")
	counters                    = flag.Int("counters", 100, "Number of counters to generate")
	httpRequestsResponseTime    prometheus.Histogram
	httpRequestsResponseTimeSum prometheus.Summary
)

func main() {
	flag.Parse()

	for i := 0; i < *gauges; i++ {
		g := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "random_gauge_" + strconv.Itoa(i),
			Help: "random gauge " + strconv.Itoa(i),
		})
		g.Set(float64(i))
		prometheus.MustRegister(g)
	}

	for i := 0; i < *counters; i++ {
		c := prometheus.NewCounter(prometheus.CounterOpts{
			Name: "random_counter_" + strconv.Itoa(i),
			Help: "random counter " + strconv.Itoa(i),
		})
		c.Add(float64(i))
		prometheus.MustRegister(c)
	}

	httpRequestsResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "response_time_seconds_histogram",
		Help:      "Request response times",
	})
	prometheus.MustRegister(httpRequestsResponseTime)

	httpRequestsResponseTimeSum = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "http",
		Name:      "response_time_seconds_summary",
		Help:      "Request response times",
	})
	prometheus.MustRegister(httpRequestsResponseTimeSum)

	// Expose the registered metrics via HTTP.
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	withMetrics := middleware(handler)
	log.Fatal(http.ListenAndServe(*addr, withMetrics))
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// have a better distribution
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		httpRequestsResponseTime.Observe(float64(time.Since(start).Seconds()))
		httpRequestsResponseTimeSum.Observe(float64(time.Since(start).Seconds()))
	})
}
