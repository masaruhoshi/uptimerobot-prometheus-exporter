package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const prog = "uptimerobot-exporter"

var (
	version string
	commit  string
	date    string
	branch  string
)

// Version - Returns the application version
func Version() string {
	return fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
}

// RootCmd Cobra Command instance
var RootCmd = &cobra.Command{
	Use:   prog,
	Short: fmt.Sprintln("A Prometheus exporter for Uptimerobot metrics."),
}
