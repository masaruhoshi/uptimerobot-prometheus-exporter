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
func ScrapeUptimeRobot(client *api.Client, ch chan<- prometheus.Metric) error {
	offset, totalScraped, totalMonitors := 0, 0, 0
	scappedMonitors := make(map[int]bool)

	for {
		xmlMonitors, err := getMonitors(client, offset)
		if err != nil {
			return err
		}
		totalMonitors = xmlMonitors.Pagination.Total
		for _, monitor := range xmlMonitors.Monitors {
			up := 1.0

			if scappedMonitors[monitor.ID] {
				log.Warnf("Trying to scrape a duplicate monitor for %s", monitor.FriendlyName)
				continue
			}
			if monitor.ResponseTimes == nil {
				log.Warnf("No response times collected for %s", monitor.FriendlyName)
				continue
			}
			status, _ := strconv.ParseFloat(monitor.Status, 64)
			responseTime := float64(monitor.ResponseTimes[0].Value)
			if status != statusUp {
				up = 0
			}

			log.Infof("Scrapping metric for %s (%s)", monitor.FriendlyName, monitor.URL)
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
			totalScraped++
			scappedMonitors[monitor.ID] = true
		}
		log.Infof("ScrapeUptimeRobot scraped %d monitors", totalMonitors)
		if totalScraped < totalMonitors {
			offset++
		} else {
			log.Infof("Scraped %d monitors", totalScraped)
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
