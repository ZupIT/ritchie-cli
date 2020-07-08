package template

const (
	StartFile = "main"

	Main = `#!/usr/bin/python3
import os

from {{bin-name}} import {{bin-name}}

input1 = os.environ.get('SAMPLE_TEXT')
input2 = os.environ.get('SAMPLE_LIST')
input3 = os.environ.get('SAMPLE_BOOL')
{{bin-name}}.Run(input1, input2, input3)
`

	Dockerfile = `
FROM python:3

COPY . .

RUN chmod +x set_umask.sh

WORKDIR /app

ENTRYPOINT ["/set_umask.sh"]
CMD ["python3 /main.py"]`

	File = `#!/usr/bin/python3
import time

def Run(input1, input2, input3):
    print("Hello World!")
    print("You receive {} in text.".format(input1))
    print("You receive {} in list.".format(input2))
    print("You receive {} in boolean.".format(input3))
`

	Makefile = `# Make Run Python
BINARY_NAME={{bin-name}}.py
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin
build:
	mkdir -p $(DIST_DIR)
	cp main.py $(DIST_DIR) && cp -r {{bin-name}} Dockerfile set_umask.sh $(DIST_DIR)
	chmod +x $(DIST_DIR)/main.py
	echo 'python main.py' >> $(DIST_DIR)/$(BINARY_NAME_WINDOWS)`

	WindowsBuild = `:: Python parameters
echo off
SETLOCAL
SET BINARY_NAME_UNIX={{bin-name}}.sh
SET BINARY_NAME_WINDOWS={{bin-name}}.bat
SET DIST=..\dist
SET DIST_DIR=%DIST%\commons\bin
:build
    mkdir %DIST_DIR%
    echo python main.py >> %DIST_DIR%\%BINARY_NAME_WINDOWS%
    xcopy {{bin-name}} %DIST_DIR%\{{bin-name}} /E /C /I
    for %%i in (main.py Dockerfile set_umask.sh) do copy %%i %DIST_DIR%
    GOTO DONE
:DONE`
)
