package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/collector"
	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/log"
	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const prog = "uptimerobot-exporter"

// RootCmd Cobra Command instance
var RootCmd = &cobra.Command{
	Use:   prog,
	Short: fmt.Sprintf("A Prometheus exporter for Uptimerobot metrics. %s", version.Info()),
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("version") {
			fmt.Printf("%s %s\n", cmd.Name(), version.Info())
			return nil
		}

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			// nolint: errcheck
			w.Write([]byte("OK"))
		})

		log.Infof("Starting %s %s", prog, version.Version)

		metricsPath := viper.GetString("web.metrics-path")
		http.Handle(metricsPath, prometheusHandler())

		listenAddress := viper.GetString("web.listen-address")
		log.Infoln("Listening on", listenAddress)
		err := http.ListenAndServe(listenAddress, nil)
		return err
	},
}

func init() {
	RootCmd.Flags().String("web.listen-address", ":9429", "Address on which to expose metrics and web interface")
	viper.BindPFlag("web.listen-address", RootCmd.Flags().Lookup("web.listen-address"))

	RootCmd.Flags().String("web.metrics-path", "/metrics", "Path under which to expose metrics")
	viper.BindPFlag("web.metrics-path", RootCmd.Flags().Lookup("web.metrics-path"))
}

func prometheusHandler() http.Handler {
	apiKey := os.Getenv("UPTIMEROBOT_API_KEY")
	collect := collector.Collect{}

	uptimerobotCollector := collector.New(apiKey, collect)
	prometheus.MustRegister(uptimerobotCollector)

	handler := promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer,
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))

	return handler
}
