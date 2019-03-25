// Scrape heartbeat data.

package collector

import (
	"strconv"

	"github.com/masaruhoshi/uptimerobot-go.v2/api"
	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/log"
	"github.com/prometheus/client_golang/prometheus"
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
	MonitorSslInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, monitor, "ssl"),
		"Information about the SSL certificate (if present)",
		[]string{"name", "type", "url"}, nil,
	)
)

// ScrapeUptimeRobot : scrapes from UptimeRobot API.
func ScrapeUptimeRobot(client *api.Client, ch chan<- prometheus.Metric) error {
	totalScraped, totalMonitors := 0, 0
	scrappedMonitors := make(map[int]bool)

	for {
		xmlMonitors, err := getMonitors(client, totalScraped)
		if err != nil {
			return err
		}
		totalMonitors = xmlMonitors.Pagination.Total
		monitors := xmlMonitors.Monitors

		// There is no reason to continue
		if len(monitors) == 0 {
			log.Warnf("No monitor returned")
			return nil
		}
		for _, monitor := range monitors {
			up := 1.0
			status := 0.0
			responseTime := 0.0

			if scrappedMonitors[monitor.ID] {
				log.Warnf("Trying to scrape a duplicate monitor for %s", monitor.FriendlyName)
				continue
			}
			if monitor.ResponseTimes == nil {
				log.Warnf("No response times collected for %s", monitor.FriendlyName)
			} else {
				status, _ = strconv.ParseFloat(monitor.Status, 64)
				responseTime = float64(monitor.ResponseTimes[0].Value)
			}
			if status != statusUp {
				up = 0
			}

			log.Infof("Scraping metric for %s (%s)", monitor.FriendlyName, monitor.URL)
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
			if monitor.Ssl.Product != "" && monitor.Ssl.Brand != "" {
				sslExpiration := float64(monitor.Ssl.Expires)
				ch <- prometheus.MustNewConstMetric(
					MonitorSslInfo,
					prometheus.GaugeValue,
					sslExpiration,
					monitor.FriendlyName,
					monitor.Type,
					monitor.URL,
				)
			}
			totalScraped++
			scrappedMonitors[monitor.ID] = true
		}
		log.Infof("ScrapeUptimeRobot scraped %d monitors", totalMonitors)
		if totalScraped >= totalMonitors {
			log.Infof("Scraped %d monitors (out of %d)", totalScraped, totalMonitors)
			return nil
		}
		// If no monitor was scrapped, something is wrong with the API call
		if totalScraped == 0 && totalMonitors > 0 {
			log.Warnf("No monitor scrapped. Check UptimeRobot API")
			return nil
		}
	}
}

func getMonitors(client *api.Client, offset int) (*api.XMLMonitors, error) {
	monitorsRequest := client.Monitors()
	var request = api.GetMonitorsRequest{
		ResponseTimes:      1,
		ResponseTimesLimit: 1,
		Offset:             offset,
	}

	response, err := monitorsRequest.Get(request)
	if err != nil {
		log.Errorln("Error getting monitorsRequest", err)
		return nil, err
	}
	log.Infof("Response from UptimeRobot API", response)

	if response == nil {
		log.Errorln("No monitor response: %v", response)
		return nil, err
	}

	return response, nil
}
