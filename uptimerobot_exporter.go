package main

import (
	"flag"
	"fmt"
	slog "log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddressFlag = flag.String("web.listen-address", ":9429", "Address on which to expose metrics and web interface.")
	metricsPathFlag   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
	version           = flag.Bool("version", false, "Print uptimerobot_exporter version")
)

func prometheusHandler() http.Handler {
	handler := prometheus.Handler()
	return handler
}

func startWebServer() {
	handler := prometheusHandler()

	registerCollector()

	http.Handle(*metricsPathFlag, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>Uptimerobot Exporter</title></head>
<body>
<h2>Prometheus Uptimerobot Exporter</h2>
<p><a href='` + *metricsPathFlag + `'>Metrics</a></p>
</body>
</html>`))
	})

	server := &http.Server{
		Addr:     *listenAddressFlag,
		ErrorLog: createHTTPServerLogWrapper(),
	}

	var err error
	fmt.Printf("Listening on %s (scheme=HTTP, secured=no, clientValidation=no)\n", server.Addr)
	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func registerCollector() {
	apiKey := os.Getenv("UPTIMEROBOT_API_KEY")
	collect := collector.Collect{}

	uptimerobotCollector := collector.New(apiKey, collect)
	prometheus.MustRegister(uptimerobotCollector)
}

type bufferedLogWriter struct {
	buf []byte
}

func (w *bufferedLogWriter) Write(p []byte) (n int, err error) {
	glog.Info(strings.TrimSpace(strings.Replace(string(p), "\n", " ", -1)))
	return len(p), nil
}

func createHTTPServerLogWrapper() *slog.Logger {
	return slog.New(&bufferedLogWriter{}, "", 0)
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println("uptimerobot_exporter version: {{VERSION}}")
		return
	}

	startWebServer()
}
