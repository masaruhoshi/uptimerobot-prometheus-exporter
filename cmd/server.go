package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the HTTP server",
		RunE:  serverRun,
	}
)

func init() {
	RootCmd.AddCommand(serverCmd)

	RootCmd.Flags().String("web.listen-address", ":9429", "Address on which to expose metrics and web interface")
	viper.BindPFlag("web.listen-address", RootCmd.Flags().Lookup("web.listen-address"))

	RootCmd.Flags().String("web.metrics-path", "/metrics", "Path under which to expose metrics")
	viper.BindPFlag("web.metrics-path", RootCmd.Flags().Lookup("web.metrics-path"))
}

var (
	listenAddressFlag = flag.String("web.listen-address", ":9429", "Address on which to expose metrics and web interface.")
	metricsPathFlag   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func prometheusHandler() http.Handler {
	handler := prometheus.Handler()
	return handler
}

func serverRun(cmd *cobra.Command, args []string) error {
	if viper.GetBool("version") {
		fmt.Printf("%s %s\n", cmd.Name(), Version())
		return nil
	}

	registerCollector()

	log.Infof("Starting %s", prog, Version())

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// nolint: errcheck
		w.Write([]byte("OK"))
	})

	metricsPath := viper.GetString("web.metrics-path")
	http.Handle(metricsPath, prometheusHandler())

	listenAddress := viper.GetString("web.listen-address")
	log.Infoln("Listening on", listenAddress)
	err := http.ListenAndServe(listenAddress, nil)
	return err
}

func registerCollector() {
	apiKey := os.Getenv("UPTIMEROBOT_API_KEY")
	collect := collector.Collect{}

	uptimerobotCollector := collector.New(apiKey, collect)
	prometheus.MustRegister(uptimerobotCollector)
}
