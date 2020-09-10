# Make Run Rust
BIN_NAME=formula
BIN_FOLDER=bin
UNIX-SH=run.sh
WINDOWS-BAT=run.bat

build: rust-build docker

rust-build:
	mkdir -p $(BIN_FOLDER)
	cp -r src/* $(BIN_FOLDER)
	cargo build --manifest-path $(BIN_FOLDER)/Cargo.toml --release

	# Unix
	echo '#!/bin/sh' > $(BIN_FOLDER)/$(UNIX-SH)
	echo '"$$(dirname "$$0")"/target/release/$(BIN_NAME)' >>  $(BIN_FOLDER)/$(UNIX-SH)
	chmod +x $(BIN_FOLDER)/$(UNIX-SH)

	# Windows
	echo '@ECHO OFF' > $(BIN_FOLDER)/$(WINDOWS-BAT)
	echo 'SET mypath=%~dp0' >> $(BIN_FOLDER)/$(WINDOWS-BAT)
	echo 'start /B /WAIT %mypath:~0,-1%/target/release/$(BIN_NAME)' >> $(BIN_FOLDER)/$(WINDOWS-BAT)

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
