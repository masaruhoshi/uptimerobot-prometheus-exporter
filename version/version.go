package version

import (
	"fmt"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Version - overwritten at build-time
	Version string
	// Revision - overwritten at build-time
	Revision string
	// Branch - overwritten at build-time
	Branch string
	// BuildTime - overwritten at build-time
	BuildTime string

	// GoVersion -
	GoVersion = runtime.Version()
)

// NewCollector returns a collector which exports metrics about current version information.
func NewCollector(program string, apiVersion string) *prometheus.GaugeVec {
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: program,
			Name:      "build_info",
			Help: fmt.Sprintf(
				"A metric with a constant '1' value labeled by version, revision, branch, and goversion from which %s was built.",
				program,
			),
		},
		[]string{"version", "revision", "branch", "goversion", "buildTime", "dockerAPIVersion"},
	)
	buildInfo.WithLabelValues(Version, Revision, Branch, GoVersion, BuildTime, apiVersion).Set(1)
	return buildInfo
}

// Info returns version, branch and revision information.
func Info() string {
	return fmt.Sprintf("(version=%s, branch=%s, revision=%s, buildTime=%s)", Version, Branch, Revision, BuildTime)
}
