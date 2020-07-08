package creator

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/golang"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/java"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/node"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/php"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/python"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/shell"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/lang/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
)

type CreateManager struct {
	treeManager tree.Manager
	dir         stream.DirCreater
	file        stream.FileWriteReadExister
}

func NewCreator(tm tree.Manager, dir stream.DirCreater, file stream.FileWriteReadExister) CreateManager {
	return CreateManager{treeManager: tm, dir: dir, file: file}
}

func (c CreateManager) Create(cf formula.Create) error {
	if err := c.isValidCmd(cf.FormulaCmd); err != nil {
		return err
	}

	if err := c.dir.Create(cf.WorkspacePath); err != nil {
		return err
	}

	pkgName := cf.PkgName()
	fCmdName := cf.FormulaCmdName()

	if err := c.generateFormulaFiles(cf.FormulaPath, pkgName, cf.Lang, fCmdName, cf.WorkspacePath); err != nil {
		return err
	}

	if c.isNew(cf.WorkspacePath) {
		if err := createGitIgnoreFile(cf.WorkspacePath); err != nil {
			return err
		}
		if err := createMainReadMe(cf.WorkspacePath); err != nil {
			return err
		}
	}

	return nil
}

func (c CreateManager) isValidCmd(fCmd string) error {
	trees, err := c.treeManager.Tree()
	if err != nil {
		return err
	}

	s := strings.Split(fCmd, " ")
	cp := fmt.Sprintf("root_%s", strings.Join(s[1:len(s)-1], "_"))
	u := s[len(s)-1]
	for _, v := range trees {
		for _, j := range v.Commands {
			if j.Parent == cp && j.Usage == u {
				return ErrRepeatedCommand

			}
		}
	}
	return nil
}

func (c CreateManager) generateFormulaFiles(fPath, pkgName, lang, fCmdName, workSpcPath string) error {

	if err := c.dir.Create(fPath); err != nil {
		return err
	}

	if err := createHelpFiles(fCmdName, workSpcPath); err != nil {
		return err
	}

	if err := createReadMeFile(fCmdName, fPath); err != nil {
		return err
	}

	if err := createConfigFile(fPath); err != nil {
		return err
	}

	if err := c.createSrcFiles(fPath, pkgName, lang); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) createSrcFiles(dir, pkg, language string) error {
	srcDir := fmt.Sprintf("%s/src", dir)
	pkgDir := fmt.Sprintf("%s/%s", srcDir, pkg)
	if err := fileutil.CreateDirIfNotExists(srcDir, os.ModePerm); err != nil {
		return err
	}
	switch language {
	case formula.GoLang:
		pkgDir := fmt.Sprintf("%s/pkg/%s", srcDir, pkg)
		goCreator := golang.New(c, c.createGenericFiles)
		if err := goCreator.Create(srcDir, pkg, pkgDir, dir); err != nil {
			return err
		}
	case formula.JavaLang:
		javaCreator := java.New(c, c.createGenericFiles)
		if err := javaCreator.Create(srcDir, pkg, pkgDir, dir); err != nil {
			return err
		}
	case formula.NodeLang:
		nodeCreator := node.New(c, c.createGenericFiles)
		if err := nodeCreator.Create(srcDir, pkg, pkgDir, dir); err != nil {
			return err
		}
	case formula.PhpLang:
		phpCreator := php.New(c, c.createGenericFiles)
		if err := phpCreator.Create(srcDir, pkg, pkgDir, dir); err != nil {
			return err
		}
	case formula.PythonLang:
		pythonCreator := python.New(c, c.createGenericFiles)
		if err := pythonCreator.Create(srcDir, pkg, pkgDir, dir); err != nil {
			return err
		}
	case formula.ShellLang:
		shellCreator := shell.New(c, c.createGenericFiles)
		if err := shellCreator.Create(srcDir, pkg, pkgDir, dir); err != nil {
			return nil
		}
	}
	return nil
}

func (c CreateManager) createGenericFiles(srcDir, pkg, dir string, l formula.Lang) error {
	if err := createMainFile(srcDir, pkg, l.Main, l.FileFormat, l.StartFile, l.UpperCase); err != nil {
		return err
	}

	if err := c.createMakefileForm(srcDir, pkg, dir, l.Makefile, l.Compiled); err != nil {
		return err
	}

	if err := c.createWindowsBuild(srcDir, pkg, l.WindowsBuild); err != nil {
		return err
	}

	if err := createDockerfile(pkg, srcDir, l.Dockerfile); err != nil {
		return err
	}

	if err := createUmask(srcDir); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) createWindowsBuild(dir, name, tpl string) error {
	if tpl == "" {
		return nil
	}

	tpl = strings.ReplaceAll(tpl, formula.NameBin, name)

	buildFile := path.Join(dir, "/build.bat")
	return c.file.Write(buildFile, []byte(tpl))
}

func (c CreateManager) createMakefileForm(dir, name, pathName, tpl string, compiled bool) error {
	makefilePath := path.Join(dir, formula.MakefilePath)
	if compiled {
		tpl = strings.ReplaceAll(tpl, "{{name}}", name)
		tpl = strings.ReplaceAll(tpl, "{{form-path}}", pathName)
		return c.file.Write(makefilePath, []byte(tpl))
	}
	tpl = strings.ReplaceAll(tpl, formula.NameBin, name)
	return c.file.Write(makefilePath, []byte(tpl))
}

func createDockerfile(pkg, dir, tpl string) error {
	tpl = strings.ReplaceAll(tpl, "{{bin-name}}", pkg)
	return fileutil.WriteFile(fmt.Sprintf("%s/Dockerfile", dir), []byte(tpl))
}

func createUmask(dir string) error {
	uMaskFile := fmt.Sprintf("%s/set_umask.sh", dir)
	return fileutil.WriteFile(uMaskFile, []byte(template.Umask))
}

func createMainFile(dir, pkg, tpl, fileFormat, startFile string, uc bool) error {
	if uc {
		tpl = strings.ReplaceAll(tpl, formula.NameBin, pkg)
		tpl = strings.ReplaceAll(tpl, formula.NameBinFirstUpper, strings.Title(strings.ToLower(pkg)))
		return fileutil.WriteFile(fmt.Sprintf("%s/%s%s", dir, startFile, fileFormat), []byte(tpl))
	}
	tpl = strings.ReplaceAll(tpl, formula.NameModule, pkg)
	tpl = strings.ReplaceAll(tpl, formula.NameBin, pkg)
	return fileutil.WriteFilePerm(fmt.Sprintf("%s/%s%s", dir, startFile, fileFormat), []byte(tpl), 0777)
}

func createConfigFile(dir string) error {
	tplFile := template.Config
	return fileutil.WriteFile(fmt.Sprintf("%s/config.json", dir), []byte(tplFile))
}

func createHelpFiles(formulaCmdName, workSpacePath string) error {
	dirs := strings.Split(formulaCmdName, " ")
	for i := 0; i < len(dirs); i++ {
		d := dirs[0 : i+1]
		tPath := path.Join(workSpacePath, path.Join(d...))
		helpPath := fmt.Sprintf("%s/help.txt", tPath)
		if !fileutil.Exists(helpPath) {
			err := fileutil.WriteFile(helpPath, []byte(template.Help))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createReadMeFile(formulaCmdName, formulaPath string) error {
	tpl := strings.ReplaceAll(template.ReadMe, "{{FormulaCmd}}", formulaCmdName)
	return fileutil.WriteFile(fmt.Sprintf("%s/README.md", formulaPath), []byte(tpl))
}

func createGitIgnoreFile(workspacePath string) error {
	tpl := template.GitIgnore
	return fileutil.WriteFile(fmt.Sprintf("%s/.gitignore", workspacePath), []byte(tpl))
}

func createMainReadMe(workspacePath string) error {
	tpl := template.MainReadMe
	return fileutil.WriteFile(fmt.Sprintf("%s/README.md", workspacePath), []byte(tpl))
}

func (c CreateManager) existsGitIgnore(workspacePath string) bool {
	gitIgnorePath := path.Join(workspacePath, ".gitignore")
	if !c.file.Exists(gitIgnorePath) {
		return false
	}

	read, err := c.file.Read(gitIgnorePath)
	if err != nil {
		return false
	}

	return len(read) > 0
}

func (c CreateManager) existsMainReadMe(workspacePath string) bool {
	mainReadMePath := path.Join(workspacePath, "README.md")
	if !c.file.Exists(mainReadMePath) {
		return false
	}

	read, err := c.file.Read(mainReadMePath)
	if err != nil {
		return false
	}

	return len(read) > 0
}

func (c CreateManager) isNew(workspacePath string) bool {
	return !c.existsGitIgnore(workspacePath) || !c.existsMainReadMe(workspacePath)
}
