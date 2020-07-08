package creator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

	"github.com/thoas/go-funk"

	"github.com/ZupIT/ritchie-cli/pkg/api"
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
	formulaName := cf.FormulaName()

	if err := c.generateFormulaFiles(cf.FormulaPath, pkgName, cf.Lang); err != nil {
		return err
	}

	if c.isNew(cf.WorkspacePath) {

		if err := c.createScript(cf.WorkspacePath); err != nil {
			return err
		}

		if err := c.createMakefileMain(cf.WorkspacePath, cf.FormulaPath, formulaName); err != nil {
			return err
		}

	} else {
		if err := c.changeMakefileMain(cf.WorkspacePath, cf.FormulaCmd, formulaName); err != nil {
			return err
		}
	}

	// Add the command to tree.json only when all other steps are successful
	if err := c.generateTreeJsonFile(cf.WorkspacePath, cf.FormulaCmd, cf.Lang); err != nil {
		return err
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

func (c CreateManager) generateTreeJsonFile(formPath, fCmd, lang string) error {
	treeCommands := formula.Tree{Commands: api.Commands{}}
	treePath := path.Join(formPath, formula.TreePath)
	if !c.file.Exists(treePath) {
		if err := c.dir.Create(filepath.Dir(treePath)); err != nil {
			return err
		}
	} else {
		jsonFile, err := c.file.Read(treePath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(jsonFile, &treeCommands); err != nil {
			return err
		}
	}

	treeCommands, err := updateTree(fCmd, treeCommands, lang, 0)
	if err != nil {
		return err
	}
	treeJsonFile, _ := json.Marshal(&treeCommands)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, treeJsonFile, "", "\t"); err != nil {
		return err
	}

	return c.file.Write(treePath, prettyJSON.Bytes())
}

func (c CreateManager) createMakefileMain(dir, dirForm, name string) error {
	tplFile := template.MakefileMain

	tplFile = strings.ReplaceAll(tplFile, "{{formName}}", strings.ToUpper(name))
	tplFile = strings.ReplaceAll(tplFile, "{{formPath}}", dirForm)

	return c.file.Write(path.Join(dir, formula.MakefilePath), []byte(tplFile))
}

func (c CreateManager) generateFormulaFiles(formulaPath, pkgName, lang string) error {
	if err := c.dir.Create(formulaPath); err != nil {
		return err
	}

	if err := createConfigFile(formulaPath); err != nil {
		return err
	}

	if err := c.createSrcFiles(formulaPath, pkgName, lang); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) changeMakefileMain(formPath, fCmd, fName string) error {
	d := strings.Split(fCmd, " ")
	makefilePath := path.Join(formPath, formula.MakefilePath)
	makeFile, err := c.file.Read(makefilePath)
	if err != nil {
		return err
	}

	variable := strings.ToUpper(fName) + "=" + strings.Join(d[1:], "/")
	makeFile = []byte(strings.ReplaceAll(string(makeFile), "\nFORMULAS=", "\n"+variable+"\nFORMULAS="))
	formulas := formulaValue(makeFile)
	makeFile = []byte(strings.ReplaceAll(string(makeFile), formulas, formulas+" $("+strings.ToUpper(fName)+")"))

	if err = c.file.Write(makefilePath, makeFile); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) createScript(dir string) error {
	tplFile := template.CopyBinConfig

	filePath := path.Join(dir, "/copy-bin-configs.sh")
	if err := c.file.Write(filePath, []byte(tplFile)); err != nil {
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

func updateTree(fCmd string, t formula.Tree, language string, i int) (formula.Tree, error) {
	fc := splitFormulaCommand(fCmd)
	parent := generateParent(fc, i)

	command := funk.Filter(t.Commands, func(command api.Command) bool {
		return command.Usage == fc[i] && command.Parent == parent
	}).([]api.Command)

	if len(fc)-1 == i {
		if len(command) == 0 {
			fn := fc[len(fc)-1]
			var commands []api.Command
			commands = append(t.Commands, api.Command{
				Usage: fn,
				Help:  fmt.Sprintf("%s %s", fc[i-1], fc[i]),
				Formula: true,
				Parent: parent,
			})
			t.Commands = commands
			return t, nil
		} else {
			return formula.Tree{}, ErrRepeatedCommand
		}

	} else {
		if len(command) == 0 {
			commands := append(t.Commands, api.Command{
				Usage:  fc[i],
				Help:   generateCommandHelp(parent, fc, i),
				Parent: parent,
			})
			t.Commands = commands
		}
	}

	return updateTree(fCmd, t, language, i+1)
}

func generateCommandHelp(parent string, fc []string, i int) string {
	var help string
	if parent != "root" {
		help = fc[i-1] + " " + fc[i]
	} else {
		help = fc[i] + " commands"
	}
	return help
}

func splitFormulaCommand(formulaCommand string) []string {
	return funk.Filter(strings.Split(formulaCommand, " "), func(input string) bool {
		return input != "" && input != "rit"
	}).([]string)
}

func generateParent(fc []string, index int) string {
	if index > 0 {
		return "root_" + strings.Join(fc[0:index], "_")
	} else {
		return "root"
	}
}

func formulaValue(file []byte) string {
	fileStr := string(file)
	return strings.Split(strings.Split(fileStr, "FORMULAS=")[1], "\n")[0]
}

func (c CreateManager) existsTreeJson(workspacePath string) bool {
	treePath := path.Join(workspacePath, formula.TreePath)
	if !c.file.Exists(treePath) {
		return false
	}

	read, err := c.file.Read(treePath)
	if err != nil {
		return false
	}

	return len(read) > 0
}

func (c CreateManager) existsMakefile(workspacePath string) bool {
	makefilePath := path.Join(workspacePath, formula.MakefilePath)
	if !c.file.Exists(makefilePath) {
		return false
	}

	read, err := c.file.Read(makefilePath)
	if err != nil {
		return false
	}

	return len(read) > 0
}

func (c CreateManager) isNew(workspacePath string) bool {
	return !c.existsMakefile(workspacePath) || !c.existsTreeJson(workspacePath)
}
