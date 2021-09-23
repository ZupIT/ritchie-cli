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

	checkCommand bundle
	checkCommand ruby

#ruby-build:
	mkdir -p $BIN_FOLDER
	cp -r src/* $BIN_FOLDER
	bundle config set path vendor/bundle
	bundle install --gemfile $BIN_FOLDER/Gemfile

	{
	echo "#!/bin/sh"
	echo "cd \$(dirname \"\$0\")"
	echo "bundle config set path vendor/bundle"
	echo "ruby ./index.rb"
	} >> $BIN_FOLDER/$BINARY_NAME_UNIX
	chmod +x $BIN_FOLDER/$BINARY_NAME_UNIX

	{
	echo "@ECHO OFF"
	echo "SET mypath=%~dp0"
	echo "ruby %mypath:~0,-1%/index.rb"
	} >> $BIN_FOLDER/$BINARY_NAME_WINDOWS


#docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER
