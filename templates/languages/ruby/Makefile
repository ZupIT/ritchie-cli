# Make Run Ruby
BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin

build: ruby-build docker

ruby-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	bundle config set path vendor/bundle
	bundle install --gemfile $(BIN_FOLDER)/Gemfile

	echo '#!/bin/sh' > $(BIN_FOLDER)/$(BINARY_NAME_UNIX)
	echo 'cd "$$(dirname "$$0")"' >> $(BIN_FOLDER)/$(BINARY_NAME_UNIX)
	echo 'bundle config set path vendor/bundle' >> $(BIN_FOLDER)/$(BINARY_NAME_UNIX)
	echo 'ruby ./index.rb' >> $(BIN_FOLDER)/$(BINARY_NAME_UNIX)

	echo '@ECHO OFF' > $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)
	echo 'SET mypath=%~dp0' >> $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)
	echo 'ruby %mypath:~0,-1%/index.rb' >> $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)

	chmod +x $(BIN_FOLDER)/$(BINARY_NAME_UNIX)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
