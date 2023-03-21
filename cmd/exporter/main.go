package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		promPort = flag.Int("prom.port", 9150, "Port to expose prometheus metrics")
	)

	flag.Parse()

	reg := prometheus.NewRegistry()

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	// start listenning for http connection

	port := fmt.Sprintf(":%d", *promPort)

	log.Printf("Starting nginx exporter on %q/metrics", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Cannot start nginx exporter: %s", err)
	}
}
