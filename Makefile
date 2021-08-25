SHELL=/bin/bash
# Go aliases
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_TOOL_COVER=$(GO_CMD) tool cover
GO_GET=$(GO_CMD) get
# Binary related variables, like where to put stuff in bucket, etc
BINARY_NAME=rit
CMD_PATH=./cmd/main.go
BIN=bin
DIST=dist
DIST_MAC=$(DIST)/darwin
DIST_LINUX=$(DIST)/linux
DIST_WIN=$(DIST)/windows
# Variables used everywhere
MODULE=$(shell go list -m)
DATE=$(shell date +%D_%H:%M)
# Routing stuff
METRIC_SERVER_URL=$(shell ./.github/scripts/routing.sh metric_server)
# Build Params
BUILD_ENVS='-X $(MODULE)/pkg/metric.BasicUser=$(METRIC_BASIC_USER) -X $(MODULE)/pkg/metric.BasicPass=$(METRIC_BASIC_PASS) -X $(MODULE)/pkg/metric.ServerRestURL=$(METRIC_SERVER_URL) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)'

build-linux:
	mkdir -p $(DIST_LINUX)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -ldflags $(BUILD_ENVS) -o ./$(DIST_LINUX)/$(BINARY_NAME) -v $(CMD_PATH)

build-mac:
	mkdir -p $(DIST_MAC)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -ldflags $(BUILD_ENVS) -o ./$(DIST_MAC)/$(BINARY_NAME) -v $(CMD_PATH)

build-windows:
	mkdir -p $(DIST_WIN)
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -ldflags $(BUILD_ENVS) -o ./$(DIST_WIN)/$(BINARY_NAME).exe -v $(CMD_PATH)

build: build-linux build-mac build-windows

clean:
	rm -rf $(DIST)
	rm -rf $(BIN)

unit-test:
	./run-tests.sh

functional-test:
	mkdir -p $(BIN)
	$(GO_TEST) -v -count=1 -p 1 `go list ./functional/... | grep -v vendor/ | sort -r `

generate-translation:
	go get github.com/go-bindata/go-bindata/v3/...
	~/go/bin/go-bindata -pkg i18n -o ./internal/pkg/i18n/translations.go ./resources/i18n/...
	go mod tidy
