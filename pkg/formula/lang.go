package formula

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_go"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_java"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_node"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_python"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_shell"
)

const (
	main = "main"
	Main = "Main"
	index = "index"
	PythonName = "Python"
	PyFormat = "py"
	JavaName = "Java"
	JavaFormat = "java"
	GoName = "Go"
	GoFormat = "go"
	NodeName = "Node"
	NodeFormat = "js"
	ShellFormat = "sh"
)

type LangCreator interface {
	Create(srcDir, pkg, pkgDir, dir string) error
}

type Lang struct {
	FileFormat  string
	StartFile   string
	Main        string
	Makefile    string
	Run         string
	Dockerfile  string
	PackageJson string
	File        string
	Pkg         string
	Compiled    bool
	UpperCase   bool
}

type Python struct {
	Lang
}

func NewPython() Python {
	return Python{Lang{
		FileFormat: PyFormat,
		StartFile:  main,
		Main:       tpl_python.Main,
		Makefile:   tpl_python.Makefile,
		Dockerfile: tpl_python.Dockerfile,
		File:       tpl_python.File,
		Compiled:   false,
		UpperCase:  false,
	}}
}

func (p Python) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := createGenericFiles(srcDir, pkg, dir, p.Lang); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	pkgFile := fmt.Sprintf("%s/%s.%s", pkgDir, pkg, p.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(p.File)); err != nil {
		return err
	}

	return nil
}

type Java struct {
	Lang
}

func NewJava() Java {
	return Java{Lang{
		FileFormat: JavaFormat,
		StartFile:  Main,
		Main:       tpl_java.Main,
		Makefile:   tpl_java.Makefile,
		Run:        tpl_java.Run,
		Dockerfile: tpl_java.Dockerfile,
		File:       tpl_java.File,
		Compiled:   false,
		UpperCase:  true,
	}}
}

func (j Java) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := createGenericFiles(srcDir, pkg, dir, j.Lang); err != nil {
		return err
	}

	if err := createRunTemplate(srcDir, j.Run); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	templateFileJava := strings.ReplaceAll(j.File, nameBin, pkg)
	firstUpper := strings.Title(strings.ToLower(pkg))
	templateFileJava = strings.ReplaceAll(templateFileJava, nameBinFirstUpper, firstUpper)
	pkgFile := fmt.Sprintf("%s/%s.%s", pkgDir, firstUpper, j.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateFileJava)); err != nil {
		return err
	}

	return nil
}

type Go struct {
	Lang
}

func NewGo() Go {
	return Go{Lang{
		FileFormat: GoFormat,
		StartFile:  main,
		Main:       tpl_go.Main,
		Makefile:   tpl_go.Makefile,
		Dockerfile: tpl_go.Dockerfile,
		Pkg: tpl_go.Pkg,
		Compiled:   false,
		UpperCase:  true,
	}}
}

func (g Go) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := createGenericFiles(srcDir, pkg, dir, g.Lang); err != nil {
		return err
	}

	if err := createGoModFile(srcDir, pkg); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	templateGo := strings.ReplaceAll(g.Pkg, nameModule, pkg)
	pkgFile := fmt.Sprintf("%s/%s.%s", pkgDir,pkg ,g.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateGo)); err != nil {
		return err
	}
	return nil
}

type Node struct {
	Lang
}

func NewNode() Node {
	return Node{Lang{
		FileFormat:  NodeFormat,
		StartFile:   index,
		Main:        tpl_node.Index,
		Makefile:    tpl_node.Makefile,
		Run:         tpl_node.Run,
		Dockerfile:  tpl_node.Dockerfile,
		PackageJson: tpl_node.PackageJson,
		File:        tpl_node.File,
		Compiled:    false,
		UpperCase:   true,
	}}
}

func (n Node) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := createGenericFiles(srcDir, pkg, dir, n.Lang); err != nil {
		return err
	}

	if err := createRunTemplate(srcDir, n.Run); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	if err := createPackageJson(srcDir, n.PackageJson); err != nil {
		return err
	}

	templateNode := strings.ReplaceAll(n.File, nameBin, pkg)
	pkgFile := fmt.Sprintf("%s/%s.%s", pkgDir, pkg, n.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateNode)); err != nil {
		return err
	}

	return nil
}

type Shell struct {
	Lang
}

func NewShell() Shell {
	return Shell{Lang{
		FileFormat: ShellFormat,
		StartFile:  main,
		Main:       tpl_shell.Main,
		Makefile:   tpl_shell.Makefile,
		Dockerfile: tpl_shell.Dockerfile,
		File:       tpl_shell.File,
		Compiled:   false,
		UpperCase:  true,
	}}
}

func (s Shell) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := createGenericFiles(srcDir, pkg, dir, s.Lang); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	pkgFile := fmt.Sprintf("%s/%s.%s", pkgDir, pkg, s.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(s.File)); err != nil {
		return err
	}

	return nil
}
