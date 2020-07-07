package template

const (
	StartFile = "index"

	Index = `const run = require("./{{bin-name}}/{{bin-name}}")

const INPUT1 = process.env.SAMPLE_TEXT
const INPUT2 = process.env.SAMPLE_LIST
const INPUT3 = process.env.SAMPLE_BOOL

run(INPUT1, INPUT2, INPUT3)`

	Dockerfile = `
FROM node:10

COPY . .

RUN chmod +x set_umask.sh

WORKDIR /app

ENTRYPOINT ["/set_umask.sh"]
CMD ["node /index.js"]
`

	Run = `#!/bin/sh
npm install && node index.js`

	PackageJson = `{
  "name": "src",
  "version": "1.0.0",
  "description": "Sample formula in node",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "author": "Dennis.Ritchie",
  "license": "ISC"
}`

	File = `function Run(input1, input2, input3) {
    console.log("Hello World!")
    console.log("You receive "+ input1 +" in text.");
    console.log("You receive "+ input2 +" in list.");
    console.log("You receive "+ input3 +" in boolean.");
}

const {{bin-name}} = Run
module.exports = {{bin-name}}`

	Makefile = `# Make Run Node
BINARY_NAME_UNIX={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin

build:
	mkdir -p $(DIST_DIR)
	cp run_template $(BINARY_NAME_UNIX) && chmod +x $(BINARY_NAME_UNIX)
	sed '1d' run_template > $(BINARY_NAME_WINDOWS) && chmod +x $(BINARY_NAME_WINDOWS)

	cp -r . $(DIST_DIR)

	#Clean files
	rm $(BINARY_NAME_UNIX)`

	WindowsBuild = `:: Node parameters
echo off
SETLOCAL
SET BINARY_NAME_UNIX={{bin-name}}.sh
SET BINARY_NAME_WINDOWS={{bin-name}}.bat
SET DIST=..\dist
SET DIST_DIR=%DIST%\commons\bin
:build
    mkdir %DIST_DIR%
	more +1 run_template > %DIST_DIR%\%BINARY_NAME_WINDOWS%
    copy run_template %DIST_DIR%\%BINARY_NAME_UNIX%
    xcopy . %DIST_DIR% /E /H /C /I
    GOTO DONE
:DONE`
)
