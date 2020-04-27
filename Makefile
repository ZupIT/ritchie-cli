# Go parameters
TEAM=team
SINGLE=single
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOLCOVER=$(GOCMD) tool cover
GOGET=$(GOCMD) get
BINARY_NAME=rit
SINGLE_CMD_PATH=./cmd/$(SINGLE)/main.go
TEAM_CMD_PATH=./cmd/$(TEAM)/main.go
BIN=bin
DIST=dist
DIST_MAC=$(DIST)/mac
DIST_MAC_TEAM=$(DIST_MAC)/$(TEAM)
DIST_MAC_SINGLE=$(DIST_MAC)/$(SINGLE)
DIST_LINUX=$(DIST)/linux
DIST_LINUX_TEAM=$(DIST_LINUX)/$(TEAM)
DIST_LINUX_SINGLE=$(DIST_LINUX)/$(SINGLE)
DIST_WIN=$(DIST)/windows
DIST_WIN_TEAM=$(DIST_WIN)/$(TEAM)
DIST_WIN_SINGLE=$(DIST_WIN)/$(SINGLE)
VERSION=$(RELEASE_VERSION)
GIT_REMOTE=https://$(GIT_USERNAME):$(GIT_PASSWORD)@github.com/ZupIT/ritchie-cli
MODULE=$(shell go list -m)
DATE=$(shell date +%D_%H:%M)
BUCKET=$(shell VERSION=$(VERSION) ./bucket.sh)
RITCHIE_ENV=$(shell VERSION=$(VERSION) ./ritchie_env.sh)

build:
	mkdir -p $(DIST_MAC_TEAM) $(DIST_MAC_SINGLE) $(DIST_LINUX_TEAM) $(DIST_LINUX_SINGLE) $(DIST_WIN_TEAM) $(DIST_WIN_SINGLE)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_LINUX_TEAM)/$(BINARY_NAME) -v $(TEAM_CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_MAC_TEAM)/$(BINARY_NAME) -v $(TEAM_CMD_PATH)
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_WIN_TEAM)/$(BINARY_NAME).exe -v $(TEAM_CMD_PATH)
	#LINUX SINGLE
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_LINUX_SINGLE)/$(BINARY_NAME) -v $(SINGLE_CMD_PATH)
	#MAC SINGLE
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_MAC_SINGLE)/$(BINARY_NAME) -v $(SINGLE_CMD_PATH)
	#WINDOWS 64 SINGLE
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_WIN_SINGLE)/$(BINARY_NAME).exe -v $(SINGLE_CMD_PATH)
ifneq "$(BUCKET)" ""
	echo $(BUCKET)
	aws s3 sync dist s3://$(BUCKET)/$(RELEASE_VERSION) --include "*"
	echo -n "$(RELEASE_VERSION)" > stable.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
else
	echo "NOT GONNA PUBLISH"
endif

build-circle:
	mkdir -p $(DIST_MAC_TEAM) $(DIST_MAC_SINGLE) $(DIST_LINUX_TEAM) $(DIST_LINUX_SINGLE) $(DIST_WIN_TEAM) $(DIST_WIN_SINGLE)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_LINUX_TEAM)/$(BINARY_NAME) -v $(TEAM_CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_MAC_TEAM)/$(BINARY_NAME) -v $(TEAM_CMD_PATH)
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_WIN_TEAM)/$(BINARY_NAME).exe -v $(TEAM_CMD_PATH)
	#LINUX SINGLE
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_LINUX_SINGLE)/$(BINARY_NAME) -v $(SINGLE_CMD_PATH)
	#MAC SINGLE
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_MAC_SINGLE)/$(BINARY_NAME) -v $(SINGLE_CMD_PATH)
	#WINDOWS 64 SINGLE
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags '-X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)' -o ./$(DIST_WIN_SINGLE)/$(BINARY_NAME).exe -v $(SINGLE_CMD_PATH)

release:
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"
	git add .
	git commit --allow-empty -m "release"
	git push $(GIT_REMOTE) HEAD:release-$(RELEASE_VERSION)
	git tag -a $(RELEASE_VERSION) -m "release"
	git push $(GIT_REMOTE) $(RELEASE_VERSION)
	curl --user $(GIT_USERNAME):$(GIT_PASSWORD) -X POST https://api.github.com/repos/ZupIT/ritchie-cli/pulls -H 'Content-Type: application/json' -d '{ "title": "Release $(RELEASE_VERSION) merge", "body": "Release $(RELEASE_VERSION) merge with master", "head": "release-$(RELEASE_VERSION)", "base": "master" }'

delivery:
ifneq "$(BUCKET)" ""
	echo $(BUCKET)
	aws s3 sync dist s3://$(BUCKET)/$(RELEASE_VERSION) --include "*"
	echo -n "$(RELEASE_VERSION)" > stable.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
else
	echo "NOT GONNA PUBLISH"
endif

publish:
	echo "Do nothing"

clean:
	rm -rf $(DIST)
	rm -rf $(BIN)

test:
	mkdir -p $(BIN)
	$(GOTEST) -short -coverprofile=$(BIN)/cov.out `go list ./... | grep -v vendor/`
	$(GOTOOLCOVER) -func=$(BIN)/cov.out