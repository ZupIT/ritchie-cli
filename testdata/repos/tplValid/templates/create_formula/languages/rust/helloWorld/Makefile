# The Makefile file will be deprecated in March 2021. It will be replaced by the current build.sh file
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

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
