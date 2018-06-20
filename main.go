package main

import (
	"log"
	"os"

	"github.com/masaruhoshi/uptimerobot-prometheus-exporter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
