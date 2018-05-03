// Scrape heartbeat data.

package collector

import (
	"strconv"

	"github.com/masaruhoshi/uptimerobot-go.v2/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// monitor is the Metric subsystem we use.
	monitor = "monitor"
	// statusPaused
	statusPaused = 0
	// statusNotChecked
	statusNotChecked = 1
	// statusUp
	statusUp = 2
	// statusSeemsDown
	statusSeemsDown = 8
	// statusDown
	statusDown = 9
)

// Metric descriptors.
var (
	MonitorUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, monitor, "up"),
		"Whether the target of the monitor is up.",
		[]string{"name", "type", "url"}, nil,
	)
	MonitorStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, monitor, "status"),
		"The status response of the target.",
		[]string{"name", "type", "url"}, nil,
	)
	MonitorResponseTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, monitor, "responsetime"),
		"Response Time of the monitor",
		[]string{"name", "type", "url"}, nil,
	)
)

// ScrapeUptimeRobot : scrapes from UptimeRobot API.
func ScrapeUptimeRobot(monitors []api.XMLMonitor, ch chan<- prometheus.Metric) error {
	log.Infof("ScrapeUptimeRobot found %d monitors", len(monitors))
	for _, monitor := range monitors {
		up := 1.0
		status, _ := strconv.ParseFloat(monitor.Status, 64)
		responseTime := float64(monitor.ResponseTimes[0].Value)
		if status != statusUp {
			up = 0
		}

		ch <- prometheus.MustNewConstMetric(
			MonitorUpDesc,
			prometheus.GaugeValue,
			up,
			monitor.FriendlyName,
			monitor.Type,
			monitor.URL,
		)
		ch <- prometheus.MustNewConstMetric(
			MonitorStatusDesc,
			prometheus.GaugeValue,
			status,
			monitor.FriendlyName,
			monitor.Type,
			monitor.URL,
		)
		ch <- prometheus.MustNewConstMetric(
			MonitorResponseTimeDesc,
			prometheus.GaugeValue,
			responseTime,
			monitor.FriendlyName,
			monitor.Type,
			monitor.URL,
		)
	}

	return nil
}
