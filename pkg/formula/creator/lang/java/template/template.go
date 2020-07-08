package template

const (
	StartFile = "Main"

	Main = `package {{final-pkg}};

import {{final-pkg}}.{{bin-name}}.{{bin-name-first-upper}};

public class Main {

    public static void main(String[] args) throws Exception {
        String input1 = System.getenv("SAMPLE_TEXT");
        String input2 = System.getenv("SAMPLE_LIST");
        boolean input3 = Boolean.parseBoolean(System.getenv("SAMPLE_BOOL"));
        {{bin-name-first-upper}} {{bin-name}} = new {{bin-name-first-upper}}(input1, input2, input3);
        {{bin-name}}.Run();
    }
}`

	Dockerfile = `
FROM alpine:latest
USER root

COPY . .
    
RUN apk update
RUN apk fetch openjdk8
RUN apk add openjdk8

ENV JAVA_HOME=/usr/lib/jvm/java-1.8-openjdk
ENV PATH="$JAVA_HOME/bin:${PATH}"

RUN chmod +x set_umask.sh

WORKDIR /app

ENTRYPOINT ["../set_umask.sh"]


CMD ["java -jar ../Main.jar"]`

	File = `package {{final-pkg}}.{{bin-name}};

public class {{bin-name-first-upper}} {

    private String input1;
    private String input2;
    private boolean input3;

    public void Run() throws Exception {
        System.out.printf("Hello World!\n");
        System.out.printf("You receive %s in text.\n", input1);
        System.out.printf("You receive %s in list.\n", input2);
        System.out.printf("You receive %s in boolean.\n", input3);
    }

    public {{bin-name-first-upper}}(String input1, String input2, boolean input3) {
        this.input1 = input1;
        this.input2 = input2;
        this.input3 = input3;
    }

    public String getInput1() {
        return input1;
    }

    public void setInput1(String input1) {
        this.input1 = input1;
    }

    public String getInput2() {
        return input2;
    }

    public void setInput2(String input2) {
        this.input2 = input2;
    }

    public boolean isInput3() {
        return input3;
    }

    public void setInput3(boolean input3) {
        this.input3 = input3;
    }
}`

	Makefile = `# Go parameters
BIN_FOLDER=../bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
JAR_NAME={{bin-name}}.jar

build: mvn-build sh-unix bat-windows

mvn-build:
	mkdir -p $(BIN_FOLDER)
	mvn clean install
	cp target/$(JAR_NAME) $(BIN_FOLDER)/$(JAR_NAME)
	#Clean files
	rm -Rf target

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'java -jar $(JAR_NAME)' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'java -jar $(JAR_NAME)' >> $(BAT)`

	WindowsBuild = `:: Java parameters
echo off
SETLOCAL
SET BINARY_NAME_UNIX={{bin-name}}.sh
SET BINARY_NAME_WINDOWS={{bin-name}}.bat
SET DIST=..\dist
SET DIST_DIR=%DIST%\commons\bin
:build
    mkdir %DIST_DIR%
	javac -source 1.8 -target 1.8 *.java
    echo Main-Class: Main > manifest.txt
    jar cvfm Main.jar manifest.txt *.class {{bin-name}}/*.class
    more +1 run_template > %BINARY_NAME_WINDOWS%
    copy run_template %BINARY_NAME_UNIX%
    for %%i in (Main.jar %BINARY_NAME_WINDOWS% %BINARY_NAME_UNIX% Dockerfile set_umask.sh) do copy %%i %DIST_DIR%
    erase Main.jar manifest.txt *.class {{bin-name}}\*.class %BINARY_NAME_WINDOWS% %BINARY_NAME_UNIX%
    GOTO DONE
:DONE`
	Pom = `
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>#rit{{groupId}}</groupId>
    <artifactId>#rit{{artifactId}}</artifactId>
    <version>1.0-SNAPSHOT</version>

    <properties>
        <maven.compiler.source>1.8</maven.compiler.source>
        <maven.compiler.target>1.8</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <maven-jar-plugin.version>3.2.0</maven-jar-plugin.version>
    </properties>

    <build>
        <finalName>${project.artifactId}</finalName>
        <plugins>
            <plugin>
                <!-- Build an executable JAR -->
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-jar-plugin</artifactId>
                <version>${maven-jar-plugin.version}</version>
                <configuration>
                    <archive>
                        <manifest>
                            <!-- <addClasspath>true</addClasspath> -->
                            <mainClass>#rit{{groupId}}.Main</mainClass>
                        </manifest>
                    </archive>
                </configuration>
            </plugin>
        </plugins>
    </build>

    <dependencies>
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>4.12</version>
            <scope>test</scope>
        </dependency>
    </dependencies>
</project>
`
)
