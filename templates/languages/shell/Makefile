# SH
BINARY_NAME=run.sh
BIN_FOLDER=bin

build: bash-build docker

bash-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	mv $(BIN_FOLDER)/main.sh $(BIN_FOLDER)/$(BINARY_NAME)
	chmod +x $(BIN_FOLDER)/$(BINARY_NAME)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
