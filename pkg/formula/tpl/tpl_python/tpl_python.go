package tpl_python

const (
	Main = `#!/usr/bin/python3
import os

from {{bin-name}} import {{bin-name}}

input1 = os.environ.get('SAMPLE_TEXT')
input2 = os.environ.get('SAMPLE_LIST')
input3 = os.environ.get('SAMPLE_BOOL')
{{bin-name}}.Run(input1, input2, input3)
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

	Dockerfile = `
FROM python:3

WORKDIR /app

COPY . .

RUN chmod +x set_umask.sh

ENTRYPOINT ["/app/set_umask.sh"]
CMD ["python3 main.py"]`

	File = `#!/usr/bin/python3
import time

def Run(input1, input2, input3):
    print("Hello World!")
    print("You receive {} in text.".format(input1))
    print("You receive {} in list.".format(input2))
    print("You receive {} in boolean.".format(input3))
`
)
