package collector

import (
//  "github.com/prometheus/common/log"
//  ur "github.com/gbl08ma/uptimerobot-api"
  "github.com/prometheus/client_golang/prometheus"
)

var (
  // Namespace is the namespace of the metrics
  namespace = "uptimerobot"
)

const (
  // Subsystem(s),
  exporter = "exporter"
)

type Collect struct {
  Up              bool
  Status          string
  ResponseTime    string
  ScrapeDuration  string
}

type Exporter struct {
  apiKey          string
  collect         Collect
  up              prometheus.Gauge
  status          prometheus.Gauge
  errorDesc       prometheus.Gauge
  responseTime    prometheus.Gauge
  scrapeDuration  prometheus.Gauge
  totalScrapes    prometheus.Counter
  scrapeErrors    *prometheus.CounterVec
}

func New(apiKey string, collect Collect) *Exporter {
  return &Exporter{
    apiKey: apiKey,
    collect: collect,
    up: prometheus.NewGauge(prometheus.GaugeOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "up",
      Help:      "Indicates if the monitor is up",
    }),
    status: prometheus.NewGauge(prometheus.GaugeOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "status",
      Help:      "Numeric status of the monitor",
    }),
    errorDesc: prometheus.NewGauge(prometheus.GaugeOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "last_scrape_error",
      Help:      "Total number of times an error occurred scraping UptimeRobot.",
    }),
    responseTime: prometheus.NewGauge(prometheus.GaugeOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "responsetime",
      Help:      "Most recent monitor response time",
    }),
    scrapeDuration: prometheus.NewGauge(prometheus.GaugeOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "scrape_duration_seconds",
      Help:      "Duration of uptimerobot scrape",
    }),
    scrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "scrape_errors_total",
      Help:      "Total number of times an error occurred scraping UptimeRobot.",
    }, []string{"collector"}),
    totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
      Namespace: namespace,
      Subsystem: exporter,
      Name:      "scrapes_total",
      Help:      "Total number of times UptimeRobot was scraped for metrics.",
    }),
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
  ch <- e.status
  ch <- e.errorDesc
  ch <- e.responseTime
  ch <- e.scrapeDuration
  ch <- e.totalScrapes
  e.scrapeErrors.Collect(ch)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
  e.totalScrapes.Inc()
/*
  var err error

  scrapeTime := time.Now()
  db, err := sql.Open("mysql", e.dsn)
  if err != nil {
    log.Errorln("Error opening connection to database:", err) e.error.Set(1)
    return
  }

  e.up.Set(1)

  ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "connection")

  if e.collect.GlobalStatus {
    scrapeTime = time.Now()
    if err = ScrapeGlobalStatus(db, ch); err != nil {
      log.Errorln("Error scraping for collect.global_status:", err)
      e.scrapeErrors.WithLabelValues("collect.global_status").Inc()
      e.error.Set(1)
    }
    ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.global_status")
  }
*/
}
