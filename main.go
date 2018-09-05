package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr     = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	gauges   = flag.Int("gauges", 100, "Number of gauges to generate")
	counters = flag.Int("counters", 100, "Number of counters to generate")
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

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
