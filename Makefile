# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_TOOL_COVER=$(GO_CMD) tool cover
GO_GET=$(GO_CMD) get
BINARY_NAME=rit
CMD_PATH=./cmd/main.go
BIN=bin
DIST=dist
DIST_MAC=$(DIST)/darwin
DIST_LINUX=$(DIST)/linux
DIST_WIN=$(DIST)/windows
VERSION=$(RELEASE_VERSION)
GIT_REMOTE=https://$(GIT_USERNAME):$(GIT_PASSWORD)@github.com/ZupIT/ritchie-cli
MODULE=$(shell go list -m)
DATE=$(shell date +%D_%H:%M)
BUCKET=$(shell VERSION=$(VERSION) ./.circleci/scripts/bucket.sh)
IS_RELEASE=$(shell echo ${VERSION} | egrep "^([0-9]{1,}\.)+[0-9]{1,}$")
IS_BETA=$(shell echo $(VERSION) | egrep "*.pre.*")
IS_QA=$(shell echo $(VERSION) | egrep "*qa.*")
IS_STG=$(shell echo $(VERSION) | egrep "*stg.*")
IS_NIGHTLY=$(shell echo $(VERSION) | egrep "*.nightly.*")
GONNA_RELEASE=$(shell ./.circleci/scripts/gonna_release.sh)
NEXT_VERSION=$(shell ./.circleci/scripts/next_version.sh)
METRIC_SERVER_URL=$(shell VERSION=$(VERSION) ./.circleci/scripts/ritchie_metric_server.sh)
BUILD_ENVS='-X $(MODULE)/pkg/metric.BasicUser=$(METRIC_BASIC_USER) -X $(MODULE)/pkg/metric.BasicPass=$(METRIC_BASIC_PASS) -X $(MODULE)/pkg/metric.ServerRestURL=$(METRIC_SERVER_URL) -X $(MODULE)/pkg/cmd.Version=$(VERSION) -X $(MODULE)/pkg/cmd.BuildDate=$(DATE)'

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
ifneq "$(BUCKET)" ""
	echo $(BUCKET)
	aws s3 sync dist s3://$(BUCKET)/$(RELEASE_VERSION) --include "*"
ifneq "$(IS_NIGHTLY)" ""
	echo -n "$(RELEASE_VERSION)" > nightly.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "nightly.txt"
endif
ifneq "$(IS_BETA)" ""
	echo -n "$(RELEASE_VERSION)" > beta.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "beta.txt"
endif
ifneq "$(IS_RELEASE)" ""
	echo -n "$(RELEASE_VERSION)" > stable.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
endif
ifneq "$(IS_QA)" ""
	echo -n "$(RELEASE_VERSION)" > stable.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
endif
ifneq "$(IS_STG)" ""
	echo -n "$(RELEASE_VERSION)" > stable.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
endif
else
	echo "NOT GONNA PUBLISH"
endif

build-circle: build-linux build-mac build-windows

release:
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"
	git tag -a $(RELEASE_VERSION) -m "CHANGELOG: https://github.com/ZupIT/ritchie-cli/blob/master/CHANGELOG.md"
	git push $(GIT_REMOTE) $(RELEASE_VERSION)
	gem install github_changelog_generator
	github_changelog_generator -u zupit -p ritchie-cli --token $(GIT_PASSWORD) --enhancement-labels feature,Feature --exclude-labels duplicate,question,invalid,wontfix
	git add .
	git commit --allow-empty -m "[ci skip] release"
	git push $(GIT_REMOTE) HEAD:release-$(RELEASE_VERSION)
	curl --user $(GIT_USERNAME):$(GIT_PASSWORD) -X POST https://api.github.com/repos/ZupIT/ritchie-cli/pulls -H 'Content-Type: application/json' -d '{ "title": "Release $(RELEASE_VERSION) merge", "body": "Release $(RELEASE_VERSION) merge with master", "head": "release-$(RELEASE_VERSION)", "base": "master" }'

delivery:
	@echo $(VERSION)
ifneq "$(BUCKET)" ""
	aws s3 sync dist s3://$(BUCKET)/$(RELEASE_VERSION) --include "*"
ifneq "$(IS_NIGHTLY)" ""
	echo -n "$(RELEASE_VERSION)" > nightly.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "nightly.txt"
endif
ifneq "$(IS_BETA)" ""
	echo -n "$(RELEASE_VERSION)" > beta.txt
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "beta.txt"
endif
ifneq "$(IS_RELEASE)" ""
	echo -n "$(RELEASE_VERSION)" > stable.txt
	mkdir latest
	cp dist/installer/ritchiecli.msi latest/
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "latest/ritchiecli.msi"
endif
ifneq "$(IS_QA)" ""
	echo -n "$(RELEASE_VERSION)" > stable.txt
	mkdir latest
	cp dist/installer/ritchiecli.msi latest/
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "latest/ritchiecli.msi"
endif
ifneq "$(IS_STG)" ""
	echo -n "$(RELEASE_VERSION)" > stable.txt
	mkdir latest
	cp dist/installer/ritchiecli.msi latest/
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "stable.txt"
	aws s3 sync . s3://$(BUCKET)/ --exclude "*" --include "latest/ritchiecli.msi"
endif
else
	echo "NOT GONNA PUBLISH"
endif

publish:
	echo "Do nothing"

clean:
	rm -rf $(DIST)
	rm -rf $(BIN)

unit-test:
	./run-tests.sh

functional-test:
	mkdir -p $(BIN)
	$(GO_TEST) -v -count=1 -p 1 `go list ./functional/... | grep -v vendor/ | sort -r `

rebase-nightly:
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"
	git push $(GIT_REMOTE) --delete nightly | true
	git checkout -b nightly
	git reset --hard master
	git add .
	git commit --allow-empty -m "nightly"
	git push $(GIT_REMOTE) HEAD:nightly

rebase-beta:
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"
	git push $(GIT_REMOTE) --delete beta | true
	git checkout -b beta
	git reset --hard nightly
	git add .
	git commit --allow-empty -m "beta"
	git push $(GIT_REMOTE) HEAD:beta

release-creator:
ifeq "$(GONNA_RELEASE)" "RELEASE"
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"
	git checkout -b "release-$(NEXT_VERSION)"
	git add .
	git commit --allow-empty -m "release-$(NEXT_VERSION)"
	git push $(GIT_REMOTE) HEAD:release-$(NEXT_VERSION)
else
	echo "NOT GONNA RELEASE"
endif