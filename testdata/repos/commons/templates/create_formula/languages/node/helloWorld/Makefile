# The Makefile file will be deprecated in March 2021. It will be replaced by the current build.sh file
# Make Run Node
BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin

build: node-build docker

node-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	npm install --silent --no-progress --prefix $(BIN_FOLDER)
	echo '#!/bin/sh' > $(BIN_FOLDER)/$(BINARY_NAME_UNIX)
	echo 'node "$$(dirname "$$0")"/index.js' >>  $(BIN_FOLDER)/$(BINARY_NAME_UNIX)
	echo 'node index.js' > $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)
	chmod +x $(BIN_FOLDER)/$(BINARY_NAME_UNIX)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
