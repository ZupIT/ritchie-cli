package formula

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_go"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_java"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_node"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_python"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_shell"
)

type Lang struct {
	FileFormat string
	StartFile  string
	Main       string
	Makefile   string
	Run        string
	Dockerfile string
	PackageJson string
	File string
	Compiled   bool
	UpperCase  bool
}

var Python = Lang{
	FileFormat: "py",
	StartFile:  "main",
	Main:       tpl_python.Main,
	Makefile:   tpl_python.Makefile,
	Dockerfile: tpl_python.Dockerfile,
	File: tpl_python.File,
	Compiled:   false,
	UpperCase:  false,
}

var Java = Lang{
	FileFormat: "java",
	StartFile:  "Main",
	Main:       tpl_java.Main,
	Makefile:   tpl_java.Makefile,
	Run:        tpl_java.Run,
	Dockerfile: tpl_java.Dockerfile,
	File: tpl_java.File,
	Compiled:   false,
	UpperCase:  true,
}

var Go = Lang{
	FileFormat: "go",
	StartFile:  "main",
	Main:       tpl_go.Main,
	Makefile:   tpl_go.Makefile,
	Dockerfile: tpl_go.Dockerfile,
	Compiled:   false,
	UpperCase:  true,
}

var Node = Lang{
	FileFormat: "js",
	StartFile:  "index",
	Main:       tpl_node.Index,
	Makefile:   tpl_node.Makefile,
	Run:        tpl_node.Run,
	Dockerfile: tpl_node.Dockerfile,
	PackageJson: tpl_node.PackageJson,
	File: tpl_node.File,
	Compiled:   false,
	UpperCase:  true,
}

var Shell = Lang{
	FileFormat: "java",
	StartFile:  "Main",
	Main:       tpl_shell.Main,
	Makefile:   tpl_shell.Makefile,
	Dockerfile: tpl_shell.Dockerfile,
	File: tpl_shell.File,
	Compiled:   false,
	UpperCase:  true,
}