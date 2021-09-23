#!/bin/sh

SOURCE_FILE=src.csproj
BIN_FOLDER=bin
BIN_FOLDER_LINUX=linux
BIN_FOLDER_DARWIN=darwin
BIN_UNIX=src.dll
SH=$BIN_FOLDER/run.sh

checkCommand () {
  if ! command -v "$1" >/dev/null; then
    echo "$1 required" >&2;
		exit 1;
	fi
}

#linux-build:
	checkCommand dotnet
	mkdir -p $BIN_FOLDER
	cp -r src/* $BIN_FOLDER
	dotnet build $BIN_FOLDER/$SOURCE_FILE -o $BIN_FOLDER/$BIN_FOLDER_LINUX --configuration Release

#macOS-build:
	mkdir -p $BIN_FOLDER/$BIN_FOLDER_DARWIN
	dotnet build $BIN_FOLDER/$SOURCE_FILE -o $BIN_FOLDER/$BIN_FOLDER_DARWIN --configuration Release --runtime osx-x64

#sh-unix:
	{
	echo "#!/bin/sh"
	echo "if [ $(uname) = \"Darwin\" ]; then"
	echo "	dotnet \$(dirname \"\$0\")/$BIN_FOLDER_DARWIN/$BIN_UNIX"
	echo 'else'
	echo "	dotnet \$(dirname \"\$0\")/$BIN_FOLDER_LINUX/$BIN_UNIX"
	echo "fi"
	} >> $SH

	chmod +x $SH

#docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER
