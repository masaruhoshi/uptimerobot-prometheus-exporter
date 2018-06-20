# Copyright 2016 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO    		:= go
GORELEASER  := $(GOPATH)/bin/goreleaser
DEP   		:= $(GOPATH)/bin/dep
pkgs   		 = $(shell $(GO) list ./... | grep -v /vendor/)

PKG_NAME			?= uptimerobot-exporter
PREFIX              ?= $(shell pwd)
BIN_DIR             ?= $(shell pwd)
DOCKER_IMAGE_NAME   ?= uptimerobot-exporter

COMMIT ?= `git rev-parse --short HEAD 2>/dev/null`
BRANCH ?= `git rev-parse --abbrev-ref HEAD 2>/dev/null`
VERSION ?= $(shell cat VERSION)
BUILD_DATE := `date -u +"%Y-%m-%dT%H:%M:%SZ"`

COMMIT_FLAG := -X `go list ./cmd`.commit=$(COMMIT)
VERSION_FLAG := -X `go list ./cmd`.version=$(VERSION)
BRANCH_FLAG := -X `go list ./cmd`.branch=$(BRANCH)
BUILD_DATE_FLAG := -X `go list ./cmd`.date=$(BUILD_DATE)

all: deps style format test build

VERSION:
	./gen_version.sh > VERSION

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

clean:
	@rm -f *.lock *.tar.gz uptimerobot-exporter
	@rm -f release/*
	@rm -fr vendor/*
	@rm -fr .build .tarballs dist bin VERSION

deps:
	@echo ">> Checking dependencies"
	@hash $(GORELEASER) 2>/dev/null || (echo "Unable to find goreleaser"; exit)
	@hash $(DEP) 2>&1 /dev/null || (curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)
	@$(DEP) ensure

$(PREFIX)/bin/$(PKG_NAME): VERSION deps $(shell find $(PREFIX) -type f -name '*.go')
	CGO_ENABLED=0 $(GO) build \
			-ldflags "-w -s $(COMMIT_FLAG) $(VERSION_FLAG) $(BRANCH_FLAG) $(BUILD_DATE_FLAG)" \
			-o $@

test:
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

build: $(PREFIX)/bin/$(PKG_NAME)

docker: deps format $(PREFIX)/bin/$(PKG_NAME)
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME)" .

.PHONY: all style format build test vet docker
