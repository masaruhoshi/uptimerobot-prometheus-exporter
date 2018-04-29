package = github.com/masaruhoshi/uptimerobot-prometheus-exporter
VERSION := $(shell cat VERSION)

all: deps build

clean:
	@rm -f *.lock uptimerobot_exporter
	@rm -f release/*
	@rm -fr vendor/*

deps:
	glide install

build: deps
	go build uptimerobot_exporter.go

release: deps
	mkdir -p release
	perl -p -i -e 's/{{VERSION}}/$(VERSION)/g' uptimerobot_exporter.go
	GOOS=darwin GOARCH=amd64 go build -o release/uptimerobot_exporter-$(VERSION).darwin-amd64 $(package)
	GOOS=linux GOARCH=amd64 go build -o release/uptimerobot_exporter-$(VERSION).linux-amd64 $(package)
	perl -p -i -e 's/$(VERSION)/{{VERSION}}/g' uptimerobot_exporter.go
	tar -czf release/uptimerobot_exporter-$(VERSION).darwin-amd64.tar.gz \
		release/uptimerobot_exporter-$(VERSION).darwin-amd64 && rm release/uptimerobot_exporter-$(VERSION).darwin-amd64
	tar -czf release/uptimerobot_exporter-$(VERSION).linux-amd64.tar.gz \
		release/uptimerobot_exporter-$(VERSION).linux-amd64 && rm release/uptimerobot_exporter-$(VERSION).linux-amd64
