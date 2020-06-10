package tpl_go

const (
	Config = `{
  "description": "Sample inputs in Ritchie.",
  "inputs" : [
    {
      "name" : "sample_text",
      "type" : "text",
      "label" : "Type : ",
      "cache" : {
        "active": true,
        "qty" : 6,
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
	CopyBinConfig = `#!/bin/sh

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
  cp -rf "$formula"/dist formulas/"$formula"
}

rm_formula_bin() {
  rm -rf "$formula"/dist
}

create_formula_checksum() {
  find "${formula}"/dist -type f -exec md5sum {} \; | sort -k 2 | md5sum | cut -f1 -d ' ' > formulas/"${formula}.md5"
}
` +
		"\ncompact_formula_bin_and_remove_them() {\n" +
		"for bin_dir in `find formulas/\"$formula\" -type d -name \"dist\"`; do\n" +
		"for binary in `ls -1 $bin_dir`; do\n" +
		"cd  ${bin_dir}/${binary}\n" +
		"zip -r \"${binary}.zip\" \"bin\"\n" +
		"mv \"${binary}\".zip ../../\n" +
		`cd - || exit
    done;
    rm -rf "${bin_dir}"
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
	GoMod = `module {{nameModule}}

go 1.14

require github.com/fatih/color v1.9.0`

	Main = `package main

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

	Makefile = `# Go parameters
BINARY_NAME={{name}}
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
CMD_PATH=./main.go
DIST=../dist
DIST_MAC_DIR=$(DIST)/darwin/bin
BIN_MAC=$(BINARY_NAME)-darwin
DIST_LINUX_DIR=$(DIST)/linux/bin
BIN_LINUX=$(BINARY_NAME)-linux
DIST_WIN_DIR=$(DIST)/windows/bin
BIN_WIN=$(BINARY_NAME)-windows.exe

build:
	mkdir -p $(DIST_MAC_DIR) $(DIST_LINUX_DIR) $(DIST_WIN_DIR)
	export MODULE=$(GO111MODULE=on go list -m)
	#LINUX
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags release -ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' -o '$(DIST_LINUX_DIR)/$(BIN_LINUX)' $(CMD_PATH) && cp -r . $(DIST_LINUX_DIR)
	#MAC
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -tags release -ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' -o '$(DIST_MAC_DIR)/$(BIN_MAC)' $(CMD_PATH) && cp -r . $(DIST_MAC_DIR)
	#WINDOWS 64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -tags release -ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' -o '$(DIST_WIN_DIR)/$(BIN_WIN)' $(CMD_PATH) && cp -r . $(DIST_WIN_DIR)

test:
	$(GOTEST) -short ` + "`go list ./... | grep -v vendor/`"

	Dockerfile = `
FROM golang:alpine AS builder

ADD . /app
WORKDIR /app
RUN go build -o main -v main.go

FROM alpine:latest


COPY --from=builder /app/main main
COPY --from=builder /app/set_umask.sh set_umask.sh
RUN chmod +x main
RUN chmod +x set_umask.sh

WORKDIR /app
ENTRYPOINT ["/set_umask.sh"]
CMD ["/main"]`

	MakefileMain = `#Makefiles
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
	Pkg = `package {{nameModule}}

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

	UnzipBinConfigs = `#!/bin/sh
find formulas -name "*.zip" | while read filename; do unzip -o -d "` + "`dirname \"$filename\"`\" \"$filename\"; rm -f \"$filename\"; done;"
)
