# Go parameters
BIN_FOLDER=bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
BIN_NAME=main
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
CMD_PATH=main.go
BIN_FOLDER_DARWIN=../$(BIN_FOLDER)/darwin
BIN_DARWIN=$(BIN_FOLDER_DARWIN)/$(BIN_NAME)
BIN_FOLDER_LINUX=../$(BIN_FOLDER)/linux
BIN_LINUX=$(BIN_FOLDER_LINUX)/$(BIN_NAME)
BIN_FOLDER_WINDOWS=../$(BIN_FOLDER)/windows
BIN_WINDOWS=$(BIN_FOLDER_WINDOWS)/$(BIN_NAME).exe


build: go-build sh-unix bat-windows docker

go-build:
	$(shell cd src; mkdir -p $(BIN_FOLDER_DARWIN) $(BIN_FOLDER_LINUX) $(BIN_FOLDER_WINDOWS))
	#LINUX
	$(shell cd src; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o '$(BIN_LINUX)' -v $(CMD_PATH))
	#MAC
	$(shell cd src; GOOS=darwin GOARCH=amd64 $(GOBUILD) -o '$(BIN_DARWIN)' -v $(CMD_PATH))
	#WINDOWS 64
	$(shell cd src; GOOS=windows GOARCH=amd64 $(GOBUILD) -o '$(BIN_WINDOWS)' -v $(CMD_PATH))

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'if [ $$(uname) = "Darwin" ]; then' >> $(SH)
	echo '  ./darwin/$(BIN_NAME)' >> $(SH)
	echo 'else' >> $(SH)
	echo '  ./linux/$(BIN_NAME)' >> $(SH)
	echo 'fi' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'cd windows' >> $(BAT)
	echo 'start /B /WAIT $(BIN_NAME).exe' >> $(BAT)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)

test:
	$(GOTEST) -short `go list ./... | grep -v vendor/`
