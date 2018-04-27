package collector

import (
	"time"

	"github.com/masaruhoshi/uptimerobot-go/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// Subsystem(s),
	exporter = "exporter"
)

var (
	// Namespace is the namespace of the metrics
	namespace = "uptimerobot"

	// Metric - Duration of last scrape
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

// Collect struct
type Collect struct {
	Up             bool
	Status         string
	ResponseTime   string
	ScrapeDuration string
}

// Exporter signature
type Exporter struct {
	apiKey         string
	collect        Collect
	up             prometheus.Gauge
	status         prometheus.Gauge
	errorDesc      prometheus.Gauge
	responseTime   prometheus.Gauge
	scrapeDuration prometheus.Gauge
	scrapeErrors   *prometheus.CounterVec
}

// New : Creates a new instance of Exporter for scraping metrics
func New(apiKey string, collect Collect) *Exporter {
	return &Exporter{
		apiKey:  apiKey,
		collect: collect,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "up",
			Help:      "Indicates if the monitor is up",
		}),
		errorDesc: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "last_scrape_error",
			Help:      "Last time error occurred scraping UptimeRobot.",
		}),
		scrapeDuration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "scrape_duration_seconds",
			Help:      "How long the last scrape took",
		}),
		scrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "scrape_errors",
			Help:      "Error occurred scraping UptimeRobot",
		}, []string{"collector"}),
	}
}

// Describe implements prometheus.Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)

	ch <- e.up
	ch <- e.errorDesc
	ch <- e.scrapeDuration
	e.scrapeErrors.Collect(ch)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	var err error

	scrapeTime := time.Now()
	client, err := api.NewClient(e.apiKey)
	if err != nil {
		log.Errorln("Error opening connection to UptimeRobot:", err)
		e.errorDesc.Set(1)
		return
	}

	monitors := client.Monitors()
	var request = api.GetMonitorsRequest{
		MonitorId: 1,
	}

	response, err := monitors.Get(request)
	if err != nil {
		log.Errorln("Error getting monitors", err)
		e.errorDesc.Set(1)
		return
	}

	if response == nil {
		log.Errorln("No monitor response: %v", response)
		e.errorDesc.Set(1)
		return
	}

	e.up.Set(1)
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "connection")

	scrapeTime = time.Now()
	if err = ScrapeUptimeRobot(response.Monitors, ch); err != nil {
		log.Errorln("Error scraping for collect.uptimerobot:", err)
		e.scrapeErrors.WithLabelValues("collect.uptimerobot").Inc()
		e.errorDesc.Set(1)
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.uptimerobot")
}
