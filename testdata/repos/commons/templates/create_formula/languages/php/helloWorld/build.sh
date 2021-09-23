#!/bin/sh

BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin

# check-dependencies
	checkCommand () {
		if ! command -v "$1" >/dev/null; then
    		echo "$1 required" >&2;
			exit 1;
		fi
	}

	checkCommand composer
	checkCommand php

# php-build:
	mkdir -p $BIN_FOLDER
	cp -r src/* $BIN_FOLDER
	composer install -q -d $BIN_FOLDER

	# Unix
	{
	echo "#!/bin/sh"
	echo "php -f \$(dirname \"\$0\")/index.php"
	} >>  $BIN_FOLDER/$BINARY_NAME_UNIX
	chmod +x $BIN_FOLDER/$BINARY_NAME_UNIX

	# Windows
	{
	echo "@echo off"
	echo "php -f index.php"
	} >> $BIN_FOLDER/$BINARY_NAME_WINDOWS

# docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER
