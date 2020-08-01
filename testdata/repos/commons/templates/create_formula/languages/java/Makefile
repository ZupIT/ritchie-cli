# JAVA
BIN_FOLDER=bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
BUILD=mvn clean install
CLEAN=mvn clean
JAR_FILE=Main-jar-with-dependencies.jar
TARGET=target

build: java-build sh-unix bat-windows
	$(CLEAN)

java-build:
	mkdir -p $(BIN_FOLDER)
	$(BUILD)
	mv $(TARGET)/$(JAR_FILE) $(BIN_FOLDER)

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'java -jar $(JAR_FILE)' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'java -jar $(JAR_FILE)' >> $(BAT)
