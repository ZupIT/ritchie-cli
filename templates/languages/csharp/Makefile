# SH
SOURCE_FILE=src/src.csproj
BIN_FOLDER=bin
BIN_FOLDER_LINUX=linux
BIN_FOLDER_WINDOWS=windows
BIN_FOLDER_DARWIN=darwin
BIN_UNIX=src
BIN_WINDOWS=src.exe
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat

build: linux-build windows-build macOS-build sh-unix bat-windows docker

linux-build:
	mkdir -p $(BIN_FOLDER)/$(BIN_FOLDER_LINUX)
	dotnet build $(SOURCE_FILE) -o $(BIN_FOLDER)/$(BIN_FOLDER_LINUX) --configuration Release --runtime linux-x64
windows-build:
	mkdir -p $(BIN_FOLDER)/$(BIN_FOLDER_WINDOWS)
	dotnet build $(SOURCE_FILE) -o $(BIN_FOLDER)/$(BIN_FOLDER_WINDOWS) --configuration Release --runtime win10-x64
macOS-build:
	mkdir -p $(BIN_FOLDER)/$(BIN_FOLDER_DARWIN)
	dotnet build $(SOURCE_FILE) -o $(BIN_FOLDER)/$(BIN_FOLDER_DARWIN) --configuration Release --runtime osx-x64

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'if [ $$(uname) = "Darwin" ]; then' >> $(SH)
	echo '  "$$(dirname "$$0")"/$(BIN_FOLDER_DARWIN)/$(BIN_UNIX)' >> $(SH)
	echo 'else' >> $(SH)
	echo '  "$$(dirname "$$0")"/$(BIN_FOLDER_LINUX)/$(BIN_UNIX)' >> $(SH)
	echo 'fi' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo '$(BIN_FOLDER_WINDOWS)/$(BIN_WINDOWS)' >> $(BAT)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
