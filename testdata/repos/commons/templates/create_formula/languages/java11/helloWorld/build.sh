#!/bin/sh
# shellcheck disable=SC2086

BIN_FOLDER=bin
SH=$BIN_FOLDER/run.sh
BAT=$BIN_FOLDER/run.bat

JAR_FILE=Main.jar
TARGET=target

# Check Dependencies
	checkCommand () {
		if ! $1 $2 | grep $3 >/dev/null; then
            echo "$1 $3x required" >&2;
			exit 1
        fi
	}

	checkCommand mvn --version " 3."
	checkCommand java --version " 11."
# Build
	mvn clean

#java-build:
	mkdir -p $BIN_FOLDER
	mvn clean install
	mv $TARGET/$JAR_FILE $BIN_FOLDER
	rm -Rf $TARGET

#sh-unix:
	{
	echo "#!/bin/sh"
	echo "java -jar \$(dirname \"\$0\")/$JAR_FILE"
	} >> $SH
	chmod +x $SH

#bat-windows:
	{
	echo "@ECHO OFF"
	echo "java -jar $JAR_FILE"
	} >> $BAT

#test:
	mvn clean install

#docker:
	cp Dockerfile set_umask.sh $BIN_FOLDER
