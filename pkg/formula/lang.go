package formula

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_java"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_python"
)

type Lang struct {
	FileFormat string
	StartFile  string
	Main       string
	Makefile   string
	Run        string
	Dockerfile string
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

var Node = Lang{
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

var Shell = Lang{
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