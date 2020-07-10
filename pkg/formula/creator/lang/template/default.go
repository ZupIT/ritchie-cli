package template

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
	Umask = `#!/bin/sh
umask 0011
$1`
	ReadMe = `
## {{FormulaCmd}}

### command
` + "```" + `bash
$ rit {{FormulaCmd}}
` + "```" + `

### description
 //TODO explain how to use this command
`
	Help      = `//TODO add some help msg`
	GitIgnore = `
# Created by https://www.gitignore.io

### Go ###
# Binaries for programs and plugins
*.dll
*.so
*.dylib

/bin/
**/bin/*
/dist/
**/dist/*
/test/tests.*
/test/coverage.*

# Test binary, built with " go test -c "
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

### Vim ###
# Swap
[._]*.s[a-v][a-z]
[._]*.sw[a-p]
[._]s[a-rt-v][a-z]
[._]ss[a-gi-z]
[._]sw[a-p]

# Session
Session.vim
Sessionx.vim

# Temporary
.netrwhist
*~
# Auto-generated tag files
# Persistent undo
[._]*.un~

### VisualStudioCode ###
.vscode/*

### VisualStudioCode Patch ###
# Ignore all local history of files
.history

### macOS ###
# General
.DS_Store
.AppleDouble
.LSOverride

# Icon must end with two \r
Icon

# Thumbnails
._*

# Files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

# End of https://www.gitignore.io/api/macos

# End of https://www.gitignore.io/api/macos
# Intellij project files
*.iml
*.ipr
*.iws
.idea/
`
	MainReadMe = `
[Contribute to the Ritchie community](https://github.com/ZupIT/ritchie-formulas/blob/master/CONTRIBUTING.md)

## Documentation

This repository contains rit formulas which can be executed by the [ritchie-cli](https://github.com/ZupIT/ritchie-cli).

- [Gitbook](https://docs.ritchiecli.io)

## Build and test formulas locally

` + "```" + `bash
$ rit build formula
` + "```" + `

## Contribute to the repository with your formulas

1. Fork the repository
2. Create a branch: ` + "`" + ` git checkout -b <branch_name>` + "`" + `
3. Check the step by step of [how to create formulas on Ritchie](https://docs.ritchiecli.io/getting-started/creating-formulas)
4. Add your formulas to the repository and commit your implementation: ` + "`" + `git commit -m '<commit_message>'` + "`" + `
5. Push your branch: ` + "`" + `git push origin <project_name>/<location>` + "`" + `
6. Open a pull request on the repository for analysis.
`
)
