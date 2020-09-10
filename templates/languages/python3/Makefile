# SH
BINARY_NAME=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin
BIN_CONFIG_VENV=config_env.sh

build: python-build sh_unix bat_windows docker

python-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)

	# Config VENV
	echo '#!/bin/bash' > $(BIN_FOLDER)/$(BIN_CONFIG_VENV)
	echo 'python3 -m venv --clear "$$(dirname "$$0")"/venv' >> $(BIN_FOLDER)/$(BIN_CONFIG_VENV)
	echo 'source ./"$$(dirname "$$0")"/venv/bin/activate' >> $(BIN_FOLDER)/$(BIN_CONFIG_VENV)
	echo 'pip3 install -r "$$(dirname "$$0")"/requirements.txt' >> $(BIN_FOLDER)/$(BIN_CONFIG_VENV)
	chmod +x $(BIN_FOLDER)/$(BIN_CONFIG_VENV)
	./$(BIN_FOLDER)/$(BIN_CONFIG_VENV)

	# Fix lib64
	rm ./$(BIN_FOLDER)/venv/lib64
	cp -r ./$(BIN_FOLDER)/venv/lib/ ./$(BIN_FOLDER)/venv/lib64

sh_unix:
	echo '#!/bin/bash' > $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'source "$$(dirname "$$0")"/venv/bin/activate' >> $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'python3 "$$(dirname "$$0")"/main.py' >> $(BIN_FOLDER)/$(BINARY_NAME)
	echo 'deactivate' >> $(BIN_FOLDER)/$(BINARY_NAME)
	chmod +x $(BIN_FOLDER)/$(BINARY_NAME)

bat_windows:
	echo '"$$(dirname "$$0")"/venv/bin/activate.bat' >> $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)
	echo 'python main.py' >> $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)
	echo 'deactivate.bat' >> $(BIN_FOLDER)/$(BINARY_NAME_WINDOWS)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
