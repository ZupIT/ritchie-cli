package templates

const (
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
`

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
)
