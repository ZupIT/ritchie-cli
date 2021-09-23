#!/bin/sh

BIN_FOLDER=bin
BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
ENTRY_POINT_UNIX=main.sh
ENTRY_POINT_WINDOWS=main.bat

# check-dependencies
	checkCommand () {
		if ! command -v "$1" >/dev/null; then
    		echo "$1 required" >&2;
			exit 1;
		fi
	}

	checkCommand terraform

# bash-build:
	mkdir -p $BIN_FOLDER
	cp -r src/* $BIN_FOLDER
	mv $BIN_FOLDER/$ENTRY_POINT_UNIX $BIN_FOLDER/$BINARY_NAME_UNIX
	chmod +x $BIN_FOLDER/$BINARY_NAME_UNIX

# bat-build:
	mv $BIN_FOLDER/$ENTRY_POINT_WINDOWS $BIN_FOLDER/$BINARY_NAME_WINDOWS

# docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER
