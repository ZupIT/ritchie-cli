package tpl_python

const (
	TemplateMain = `#!/usr/bin/python3
import os

from {{bin-name-first-upper}} import {{bin-name-first-upper}}

input1 = os.environ.get('SAMPLE_TEXT')
input2 = os.environ.get('SAMPLE_LIST')
input3 = os.environ.get('SAMPLE_BOOL')
{{bin-name-first-upper}}.Run(input1, input2, input3)
`

	TemplateMakefile = `# Make Run Python
BINARY_NAME={{bin-name}}.py
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin
build:
	mkdir -p $(DIST_DIR)
	cp main.py $(DIST_DIR)/$(BINARY_NAME) && cp -r {{bin-name-first-upper}} $(DIST_DIR)
	chmod +x $(DIST_DIR)/$(BINARY_NAME)
	echo 'python {{bin-name}}.py' >> $(DIST_DIR)/$(BINARY_NAME_WINDOWS)`

	TemplateFilePython = `#!/usr/bin/python3
import time

def Run(input1, input2, input3):
    print("Hello World!")
    print("You receive {} in text.".format(input1))
    print("You receive {} in list.".format(input2))
    print("You receive {} in boolean.".format(input3))
`
)
