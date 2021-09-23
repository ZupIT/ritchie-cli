# The Makefile file will be deprecated in March 2021. It will be replaced by the current build.sh file
# SH
BIN_FOLDER=bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
BIN_NAME=main.ps1

build: bash-build sh-unix bat-windows docker

bash-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	chmod +x $(BIN_FOLDER)/$(BINARY_NAME)

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'pwsh "$$(dirname "$$0")"/$(BIN_NAME)' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'SET mypath=%~dp0' >> $(BAT)
	echo 'start /B /D "%%mypath%%" /WAIT $(BIN_NAME)' >> $(BAT)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
