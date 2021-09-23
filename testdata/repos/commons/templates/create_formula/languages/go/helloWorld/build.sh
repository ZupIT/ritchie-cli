#!/bin/sh
# shellcheck disable=SC2046

BIN_FOLDER=bin
SH=$BIN_FOLDER/run.sh
BAT=$BIN_FOLDER/run.bat
BIN_NAME=main
CMD_PATH=main.go
BIN_FOLDER_DARWIN=../$BIN_FOLDER/darwin
BIN_DARWIN=$BIN_FOLDER_DARWIN/$BIN_NAME
BIN_FOLDER_LINUX=../$BIN_FOLDER/linux
BIN_LINUX=$BIN_FOLDER_LINUX/$BIN_NAME
BIN_FOLDER_WINDOWS=../$BIN_FOLDER/windows
BIN_WINDOWS=$BIN_FOLDER_WINDOWS/$BIN_NAME.exe



#go-build:
	cd src || exit
	mkdir -p $BIN_FOLDER_DARWIN $BIN_FOLDER_LINUX $BIN_FOLDER_WINDOWS
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $BIN_LINUX $CMD_PATH
	#MAC
	GOOS=darwin GOARCH=amd64 go build -o $BIN_DARWIN $CMD_PATH
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 go build -o $BIN_WINDOWS $CMD_PATH
	cd ..

#sh-unix:
	{
	echo "#!/bin/sh"
	echo "if [ $(uname) = \"Darwin\" ]; then"
	echo "  \$(dirname \"\$0\")/darwin/$BIN_NAME"
	echo "else"
	echo "  \$(dirname \"\$0\")/linux/$BIN_NAME"
	echo "fi"
	} >> $SH
	chmod +x $SH

#bat-windows:
	{
	echo "@ECHO OFF"
	echo "SET mypath=%~dp0"
	echo "start /B /WAIT %mypath:~0,-1%/windows/main.exe"
	} >> $BAT

#docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER

#test:
	cd src || exit
	#go test -short $(go list ./... | grep -v vendor/)
