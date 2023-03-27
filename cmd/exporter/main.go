package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	exporter "github.com/pratamaizzat/prometheus-nginx-exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		targetHost = flag.String("target.host", "localhost", "Nginx address with basic_status page")
		targetPort = flag.Int("target.port", 8081, "Nginx port with basic_status page")
		targetPath = flag.String("target.path", "/status", "URL path to scrap metrics")
		promPort   = flag.Int("prom.port", 9150, "Port to expose prometheus metrics")
	)

	flag.Parse()

	uri := fmt.Sprintf("http://%s:%d%s", *targetHost, *targetPort, *targetPath)

	// called on each collector.Collect.

	basicStats := func() ([]exporter.NginxStats, error) {

		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := netClient.Get(uri)

		if err != nil {
			log.Fatalf("netClient.Get failed %s: %s", uri, err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("io.ReadAll failed: %s", err)
		}

		r := bytes.NewReader(body)

		return exporter.ScanBasicStats(r)

	}

	// Make prometheus client aware of our collectors.
	bc := exporter.NewBasicCollector(basicStats)

	reg := prometheus.NewRegistry()
	reg.MustRegister(bc)

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
