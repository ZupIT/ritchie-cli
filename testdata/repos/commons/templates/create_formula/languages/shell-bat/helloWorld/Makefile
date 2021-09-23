# The Makefile file will be deprecated in March 2021. It will be replaced by the current build.sh file
# SH

BIN_FOLDER=bin
BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
ENTRY_POINT_UNIX=main.sh
ENTRY_POINT_WINDOWS=main.bat

build: bash-build bat-build docker

bash-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	mv $(BIN_FOLDER)/$(ENTRY_POINT_UNIX) $(BIN_FOLDER)/$(BINARY_NAME_UNIX)
	chmod +x $(BIN_FOLDER)/$(BINARY_NAME_UNIX)

bat-build:
	mv $(BIN_FOLDER)/$(ENTRY_POINT_WINDOWS) $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
