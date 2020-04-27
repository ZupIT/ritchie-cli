package tplgo

const (
	TemplateConfig = `{
  "description": "Sample inputs in Ritchie.",
  "inputs" : [
    {
      "name" : "sample_text",
      "type" : "text",
      "label" : "Type : ",
      "cache" : {
        "active": true,
        "qtd" : 6,
        "newLabel" : "Type new value. "
      }
    },
    {
      "name" : "sample_list",
      "type" : "text",
      "default" : "in1",
      "items" : ["in_list1", "in_list2", "in_list3", "in_listN"],
      "label" : "Pick your : "
    },
    {
      "name" : "sample_bool",
      "type" : "bool",
      "default" : "false",
      "items" : ["false", "true"],
      "label" : "Pick: "
    }
  ]
}`
	TemplateCopyBinConfig = `#!/bin/sh

FORMULAS="$1"

create_formulas_dir() {
  mkdir -p formulas/"$formula"
}

find_config_files() {
  files=$(find "$formula" -type f -name "*config.json")
}

copy_config_files() {
  for file in $files; do
    cp "$file" formulas/"$formula"
  done
}

copy_formula_bin() {
  cp -rf "$formula"/bin formulas/"$formula"
}

rm_formula_bin() {
  rm -rf "$formula"/bin
}

create_formula_checksum() {
  find "${formula}"/bin -type f -exec md5sum {} \; | sort -k 2 | md5sum | cut -f1 -d ' ' > formulas/"${formula}.md5"
}
` +
		"\ncompact_formula_bin_and_remove_them() {\n" +
		"for bin_dir in `find formulas \"$formula\" -type d -name \"bin\"` ; do\n" +
		"for binary in `ls -1 $bin_dir`; do\n" +
		"zip -j \"${bin_dir}/${binary}.zip\" \"${bin_dir}/${binary}\"\n" +
		"rm \"${bin_dir}/${binary}\"\n" +
		`done;
  done
}


init() {
  for formula in $FORMULAS; do
    create_formulas_dir
    find_config_files
    copy_config_files
    create_formula_checksum
    copy_formula_bin
    rm_formula_bin
    compact_formula_bin_and_remove_them
  done
}

init
`
	TemplateGoMod = `module {{nameModule}}

go 1.14

require github.com/fatih/color v1.9.0`
	TemplateMain = `package main

import (
    "os"
	"{{nameModule}}/pkg/{{nameModule}}"
)

func main() {
	input1 := os.Getenv("SAMPLE_TEXT")
	input2 := os.Getenv("SAMPLE_LIST")
	input3 := os.Getenv("SAMPLE_BOOL")

	{{nameModule}}.Input{
    	Text:    input1,
    	List:    input2,
    	Boolean: input3,
    }.Run()
}`
	TemplateMakefile = `# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME={{name}}
CMD_PATH=./main.go
DIST=../bin
DIST_MAC=$(DIST)/$(BINARY_NAME)-darwin
DIST_LINUX=$(DIST)/$(BINARY_NAME)-linux
DIST_WIN=$(DIST)/$(BINARY_NAME)-windows.exe

FORM_PATH={{form-path}}
PWD_INITIAL=$(shell pwd)

build:
	mkdir -p $(DIST)
	export MODULE=$(GO111MODULE=on go list -m)
	#LINUX
	GOOS=linux GOARCH=amd64 $(GOBUILD) -tags release -ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' -o ./$(DIST_LINUX) -v $(CMD_PATH)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -tags release -ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' -o ./$(DIST_MAC) -v $(CMD_PATH)
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -tags release -ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' -o ./$(DIST_WIN) -v $(CMD_PATH)
`
	TemplateMakefileMain = `#Makefiles
{{formName}}={{formPath}}
FORMULAS=$({{formName}})

PWD_INITIAL=$(shell pwd)

FORM_TO_UPPER  = $(shell echo $(form) | tr  '[:lower:]' '[:upper:]')
FORM = $($(FORM_TO_UPPER))

build:bin

bin:
	echo "Init pwd: $(PWD_INITIAL)"
	echo "Formulas bin: $(FORMULAS)"
	for formula in $(FORMULAS); do cd $$formula/src && make build && cd $(PWD_INITIAL); done
	./copy-bin-configs.sh "$(FORMULAS)"

test-local:
ifneq ("$(FORM)", "")
	@echo "Using form true: "  $(FORM_TO_UPPER)
	$(MAKE) bin FORMULAS=$(FORM)
	mkdir -p $(HOME)/.rit/formulas
	rm -rf $(HOME)/.rit/formulas/$(FORM)
	./unzip-bin-configs.sh
	cp -r formulas/* $(HOME)/.rit/formulas
	rm -rf formulas
else
	@echo "Use make test-local form=NAME_FORMULA for specific formula."
	@echo "form false: ALL FORMULAS"
	$(MAKE) bin
	rm -rf $(HOME)/.rit/formulas
	./unzip-bin-configs.sh
	mv formulas $(HOME)/.rit
endif
	mkdir -p $(HOME)/.rit/repo/local
	rm -rf $(HOME)/.rit/repo/local/tree.json
	cp tree/tree.json  $(HOME)/.rit/repo/local/tree.json
`
	TemplatePkg = `package {{nameModule}}

import (
	"fmt"
	"github.com/fatih/color"
)

type Input struct {
	Text string
	List string
	Boolean string
}

func(in Input)Run()  {
	fmt.Println("Hello world!")
	color.Green(fmt.Sprintf("You receive %s in text.", in.Text ))
	color.Red(fmt.Sprintf("You receive %s in list.", in.List ))
	color.Yellow(fmt.Sprintf("You receive %s in boolean.", in.Boolean ))
}`
	TemplateUnzipBinConfigs = `#!/bin/sh
find formulas -name "*.zip" | while read filename; do unzip -o -d "` + "`dirname \"$filename\"`\" \"$filename\"; rm -f \"$filename\"; done;"
)
