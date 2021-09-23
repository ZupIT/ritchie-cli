#!/bin/sh

BINARY_NAME_UNIX=run.sh
BINARY_NAME_WINDOWS=run.bat
BIN_FOLDER=bin

checkCommand () {
  if ! command -v "$1" >/dev/null; then
    echo "$1 required" >&2;
		exit 1;
	fi
}

#node-build:
	checkCommand npm
	checkCommand node
	mkdir -p $BIN_FOLDER
	cp -r src/* $BIN_FOLDER
	npm install --silent --no-progress --prefix $BIN_FOLDER
	npm run build --silent --prefix  $BIN_FOLDER

	# Unix
	{
	echo "#!/bin/sh"
	echo "node \$(dirname \"\$0\")/dist/index.js"
	} >>  $BIN_FOLDER/$BINARY_NAME_UNIX
	chmod +x $BIN_FOLDER/$BINARY_NAME_UNIX

	# Windows
	echo "node /dist/index.js" > $BIN_FOLDER/$BINARY_NAME_WINDOWS


#docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER
