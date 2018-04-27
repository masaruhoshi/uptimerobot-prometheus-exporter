package = github.com/masaruhoshi/uptimerobot_exporter
TAG := $(shell git tag | sort -r | head -n 1)

#test:
#	go test github.com/masaruhoshi/uptimerobot_exporter/collector -cover -coverprofile=collector_coverage.out -short
#	go tool cover -func=collector_coverage.out
#	@rm *.out

clean:
	@rm -f *.lock uptimerobot_exporter
	@rm -fr vendor/*
	@rm -fr coverage.txt

deps:
	glide install

build: deps
	go build uptimerobot_exporter.go

release: deps
	mkdir -p release
	perl -p -i -e 's/{{VERSION}}/$(TAG)/g' uptimerobot_exporter.go
	GOOS=darwin GOARCH=amd64 go build -o release/uptimerobot_exporter-darwin-amd64 $(package)
	GOOS=linux GOARCH=amd64 go build -o release/uptimerobot_exporter-linux-amd64 $(package)
	perl -p -i -e 's/$(TAG)/{{VERSION}}/g' uptimerobot_exporter.go
