# The Makefile file will be deprecated in March 2021. It will be replaced by the current build.sh file
# JAVA
BIN_FOLDER=bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
BUILD=mvn clean install
CLEAN=mvn clean
JAR_FILE=Main.jar
TARGET=target

build: java-build sh-unix bat-windows docker
	$(CLEAN)

java-build:
	mkdir -p $(BIN_FOLDER)
	$(BUILD)
	mv $(TARGET)/$(JAR_FILE) $(BIN_FOLDER)
	rm -Rf $(TARGET)

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'java -jar "$$(dirname "$$0")"/$(JAR_FILE)' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'java -jar $(JAR_FILE)' >> $(BAT)

test:
	mvn clean install

docker:
	cp Dockerfile set_umask.sh $(BIN_FOLDER)
