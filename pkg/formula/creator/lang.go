package creator

import (
	"fmt"
	"os"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_go"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_java"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_node"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_php"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_python"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_ruby"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/templates/template_shell"
)

const (
	main       = "main"
	Main       = "Main"
	index      = "index"
	GoLang     = "Go"
	JavaLang   = "Java"
	NodeLang   = "Node"
	PhpLang    = "Php"
	PythonLang = "Python"
	RubyLang   = "Ruby"
	ShellLang  = "Shell"
)

var Languages = []string{GoLang, JavaLang, NodeLang, PhpLang, PythonLang, RubyLang, ShellLang}

type LangCreator interface {
	Create(srcDir, pkg, pkgDir, dir string) error
}

type Lang struct {
	CreateManager
	FileFormat   string
	StartFile    string
	Main         string
	Makefile     string
	WindowsBuild string
	Run          string
	Dockerfile   string
	PackageJson  string
	File         string
	Pkg          string
	Compiled     bool
	UpperCase    bool
}

type Python struct {
	Lang
}

func NewPython(c CreateManager) Python {
	return Python{Lang{
		CreateManager: c,
		FileFormat:    fileextensions.Python,
		StartFile:     main,
		Main:          template_python.Main,
		Makefile:      template_python.Makefile,
		Dockerfile:    template_python.Dockerfile,
		File:          template_python.File,
		Compiled:      false,
		UpperCase:     false,
	}}
}

func (p Python) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := p.createGenericFiles(srcDir, pkg, dir, p.Lang); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, p.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(p.File)); err != nil {
		return err
	}

	return nil
}

type Java struct {
	Lang
}

func NewJava(c CreateManager) Java {
	return Java{Lang{
		CreateManager: c,
		FileFormat:    fileextensions.Java,
		StartFile:     Main,
		Main:          template_java.Main,
		Makefile:      template_java.Makefile,
		Run:           template_java.Run,
		Dockerfile:    template_java.Dockerfile,
		File:          template_java.File,
		Compiled:      false,
		UpperCase:     true,
	}}
}

func (j Java) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := j.createGenericFiles(srcDir, pkg, dir, j.Lang); err != nil {
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
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, firstUpper, j.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateFileJava)); err != nil {
		return err
	}

	return nil
}

type Go struct {
	Lang
}

func NewGo(c CreateManager) Go {
	return Go{Lang{
		CreateManager: c,
		FileFormat:    fileextensions.Go,
		StartFile:     main,
		Main:          template_go.Main,
		Makefile:      template_go.Makefile,
		Dockerfile:    template_go.Dockerfile,
		Pkg:           template_go.Pkg,
		Compiled:      true,
		UpperCase:     false,
	}}
}

func (g Go) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := g.createGenericFiles(srcDir, pkg, dir, g.Lang); err != nil {
		return err
	}

	if err := createGoModFile(srcDir, pkg); err != nil {
		return err
	}

	if err := fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm); err != nil {
		return err
	}

	templateGo := strings.ReplaceAll(g.Pkg, nameModule, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, g.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateGo)); err != nil {
		return err
	}
	return nil
}

type Node struct {
	Lang
}

func NewNode(c CreateManager) Node {
	return Node{Lang{
		CreateManager: c,
		FileFormat:    fileextensions.JavaScript,
		StartFile:     index,
		Main:          template_node.Index,
		Makefile:      template_node.Makefile,
		Run:           template_node.Run,
		Dockerfile:    template_node.Dockerfile,
		PackageJson:   template_node.PackageJson,
		File:          template_node.File,
		Compiled:      false,
		UpperCase:     false,
	}}
}

func (n Node) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := n.createGenericFiles(srcDir, pkg, dir, n.Lang); err != nil {
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
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, n.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateNode)); err != nil {
		return err
	}

	return nil
}

type Ruby struct {
	Lang
	Gemfile string
}

func NewRuby(c CreateManager) Ruby {
	return Ruby{
		Lang: Lang{
			CreateManager: c,
			FileFormat:    fileextensions.Ruby,
			StartFile:     index,
			Main:          template_ruby.Index,
			Makefile:      template_ruby.Makefile,
			Run:           template_ruby.Run,
			Dockerfile:    template_ruby.Dockerfile,
			File:          template_ruby.File,
			Compiled:      false,
			UpperCase:     false,
		},
		Gemfile:      	   template_ruby.Gemfile,
	}
}

func (n Ruby) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := n.createGenericFiles(srcDir, pkg, dir, n.Lang); err != nil {
		return err
	}

	if err := createRunTemplate(srcDir, n.Run); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	if err := createGemfile(srcDir, n.Gemfile); err != nil {
		return err
	}

	templateRuby := strings.ReplaceAll(n.File, nameBin, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, n.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templateRuby)); err != nil {
		return err
	}

	return nil
}

type Shell struct {
	Lang
}

func NewShell(c CreateManager) Shell {
	return Shell{Lang{
		CreateManager: c,
		FileFormat:    fileextensions.Shell,
		StartFile:     main,
		Main:          template_shell.Main,
		Makefile:      template_shell.Makefile,
		Dockerfile:    template_shell.Dockerfile,
		File:          template_shell.File,
		Compiled:      false,
		UpperCase:     false,
	}}
}

func (s Shell) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := s.createGenericFiles(srcDir, pkg, dir, s.Lang); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, s.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(s.File)); err != nil {
		return err
	}

	return nil
}

type Php struct {
	Lang
}

func NewPhp(c CreateManager) Php {
	return Php{Lang{
		CreateManager: c,
		FileFormat:    fileextensions.Php,
		StartFile:     index,
		Main:          template_php.Index,
		Makefile:      template_php.Makefile,
		Run:           template_php.Run,
		Dockerfile:    template_php.Dockerfile,
		File:          template_php.File,
		Compiled:      false,
		UpperCase:     false,
	}}
}

func (p Php) Create(srcDir, pkg, pkgDir, dir string) error {
	if err := p.createGenericFiles(srcDir, pkg, dir, p.Lang); err != nil {
		return err
	}

	if err := createRunTemplate(srcDir, p.Run); err != nil {
		return err
	}

	if err := createPkgDir(pkgDir); err != nil {
		return err
	}

	templatePHP := strings.ReplaceAll(p.File, nameBin, pkg)
	pkgFile := fmt.Sprintf("%s/%s%s", pkgDir, pkg, p.FileFormat)
	if err := fileutil.WriteFile(pkgFile, []byte(templatePHP)); err != nil {
		return err
	}

	return nil
}
