# The Makefile file will be deprecated in March 2021. It will be replaced by the current build.sh file
# SH
BINARY_NAME=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin
BIN_CONFIG_VENV=config_env.sh

build: python-build sh_unix bat_windows docker

python-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	pip3 install -r $(BIN_FOLDER)/requirements.txt --user --disable-pip-version-check

sh_unix:
	echo '#!/bin/bash' > $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'if [ -f /.dockerenv ] ; then' >> $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'pip3 install -r "$$(dirname "$$0")"/requirements.txt --user --disable-pip-version-check >> /dev/null' >> $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'fi' >> $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'python3 "$$(dirname "$$0")"/main.py' >> $(BIN_FOLDER)/$(BINARY_NAME)
	chmod +x $(BIN_FOLDER)/$(BINARY_NAME)

bat_windows:
	echo '@ECHO OFF' > $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)
	echo 'python main.py' >> $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
