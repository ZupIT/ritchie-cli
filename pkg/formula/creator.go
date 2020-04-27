package formula

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tplgo"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/thoas/go-funk"
)

type CreateManager struct {
	formPath    string
	treeManager TreeManager
	dir         stream.DirCreater
	file        stream.FileWriteReadExistRemover
}

func NewCreator(
	homePath string,
	tm TreeManager,
	dir stream.DirCreater,
	file stream.FileWriteReadExistRemover) CreateManager {
	return CreateManager{
		formPath:    fmt.Sprintf(FormCreatePathPattern, homePath),
		treeManager: tm,
		dir:         dir,
		file:        file,
	}
}

func (c CreateManager) Create(fCmd string) error {
	_ = c.dir.Create(c.formPath)
	trees, err := c.treeManager.Tree()
	if err != nil {
		return err
	}

	if err := verifyCommand(fCmd, trees); err != nil {
		return err
	}

	if c.file.Exists(fmt.Sprintf(TreeCreatePathPattern, c.formPath)) && (c.file.Exists(fmt.Sprintf("%s/%s", c.formPath, Makefile))) {
		if err := c.generateFormulaFiles(c.formPath, fCmd, false); err != nil {
			return err
		}
	} else {
		if err := c.generateFormulaFiles(c.formPath, fCmd, true); err != nil {
			return err
		}
	}
	err = c.generateTreeJsonFile(c.formPath, fCmd)
	if err != nil {
		return err
	}
	log.Println("Formula successfully created!")
	log.Printf("Your formula is in %s.", c.formPath)
	return nil
}

func (c CreateManager) generateFormulaFiles(formPath, fCmd string, new bool) error {
	d := strings.Split(fCmd, " ")
	dirForm := strings.Join(d[1:], "/")

	var dir string
	if new {
		dir = fmt.Sprintf("%s/%s", formPath, dirForm)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		err = c.createMakefileMain(formPath, dirForm, d[len(d)-1])
		if err != nil {
			return err
		}

	} else {
		dir = fmt.Sprintf("%s/%s", formPath, dirForm)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		err = c.changeMakefileMain(formPath, fCmd, d[len(d)-1])
		if err != nil {
			return err
		}
	}
	err := c.createConfigFile(dir)
	if err != nil {
		return err
	}
	err = c.createSrcFiles(dir, d[len(d)-1])
	if err != nil {
		return err
	}
	return nil
}

func (c CreateManager) generateTreeJsonFile(formPath, fCmd string) error {
	tree := Tree{Commands: []api.Command{}}
	dir := fmt.Sprintf(localTreeFile, formPath)
	jsonFile, err := c.file.Read(dir)
	if err != nil {
		if err := c.dir.Create(filepath.Dir(dir)); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(jsonFile, &tree); err != nil {
			return err
		}
	}

	tree = updateTree(fCmd, tree, 0)
	treeJsonFile, _ := json.Marshal(&tree)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, treeJsonFile, "", "\t"); err != nil {
		return err
	}

	return c.file.Write(dir, prettyJSON.Bytes())
}

func verifyCommand(fCmd string, trees map[string]Tree) error {
	s := strings.Split(fCmd, " ")

	if s[0] != "rit" {
		return errors.New("the formula's command needs to start with \"rit\" [ex.: rit group verb <noun>]")
	}

	if len(s) == 1 || len(s) == 2 {
		return errors.New("the formula's command needs at least 2 words following \"rit\" [ex.: rit group verb <noun>]")
	}
	cp := fmt.Sprintf("root_%s", strings.Join(s[1:len(s)-1], "_"))
	u := s[len(s)-1]
	for _, v := range trees {
		for _, j := range v.Commands {
			if j.Parent == cp && j.Usage == u {
				return errors.New("this command already exists")
			}
		}
	}

	return nil
}

func (c CreateManager) changeMakefileMain(formPath string, fCmd, fName string) error {
	d := strings.Split(fCmd, " ")
	dir := fmt.Sprintf("%s/%s", formPath, Makefile)
	tplFile, err := c.file.Read(dir)
	if err != nil {
		return err
	}
	variable := strings.ToUpper(d[len(d)-1]) + "=" + strings.Join(d[1:], "/")
	tplFile = []byte(strings.ReplaceAll(string(tplFile), "\nFORMULAS=", "\n"+variable+"\nFORMULAS="))
	formulas := formulaValue(tplFile)
	tplFile = []byte(strings.ReplaceAll(string(tplFile), formulas, formulas+" $("+strings.ToUpper(fName)+")"))

	if err = c.file.Write(dir, tplFile); err != nil {
		return err
	}

	return nil
}

func formulaValue(file []byte) string {
	fileStr := string(file)
	return strings.Split(strings.Split(fileStr, "FORMULAS=")[1], "\n")[0]
}

func (c CreateManager) createMakefileMain(dir, dirForm, name string) error {
	tplFile := tplgo.TemplateMakefileMain

	tplFile = strings.ReplaceAll(tplFile, "{{formName}}", strings.ToUpper(name))
	tplFile = strings.ReplaceAll(string(tplFile), "{{formPath}}", dirForm)

	err := c.createScripts(dir)
	if err != nil {
		return err
	}

	return c.file.Write(dir+"/Makefile", []byte(tplFile))
}

func (c CreateManager) createScripts(dir string) error {
	tplFile := tplgo.TemplateCopyBinConfig

	err := c.file.Write(dir+"/copy-bin-configs.sh", []byte(tplFile))
	if err != nil {
		return err
	}

	tplFile = tplgo.TemplateUnzipBinConfigs

	return c.file.Write(dir+"/unzip-bin-configs.sh", []byte(tplFile))
}

func (c CreateManager) createSrcFiles(dir, pkg string) error {
	srcDir := dir + "/src"
	err := c.dir.Create(srcDir)
	if err != nil {
		return err
	}
	err = c.createMainFile(srcDir, pkg)
	if err != nil {
		return err
	}
	err = c.createGoModFile(srcDir, pkg)
	if err != nil {
		return err
	}
	err = c.createMakefileForm(srcDir, pkg, dir)
	if err != nil {
		return err
	}
	pkgDir := srcDir + "/pkg/" + pkg
	err = c.dir.Create(pkgDir)
	if err != nil {
		return err
	}
	err = c.createPkgFile(pkgDir, pkg)
	if err != nil {
		return err
	}

	return nil
}

func (c CreateManager) createPkgFile(dir, pkg string) error {
	tplFile := tplgo.TemplatePkg
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)

	return c.file.Write(dir+"/"+pkg+".go", []byte(tplFile))
}

func (c CreateManager) createMakefileForm(dir string, name, pathName string) error {
	tplFile := tplgo.TemplateMakefile
	tplFile = strings.ReplaceAll(tplFile, "{{name}}", name)
	tplFile = strings.ReplaceAll(tplFile, "{{form-path}}", pathName)

	return c.file.Write(dir+"/Makefile", []byte(tplFile))
}

func (c CreateManager) createGoModFile(dir, pkg string) error {
	tplFile := tplgo.TemplateGoMod
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return c.file.Write(dir+"/go.mod", []byte(tplFile))
}

func (c CreateManager) createMainFile(dir, pkg string) error {
	tplFile := tplgo.TemplateMain
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return c.file.Write(dir+"/main.go", []byte(tplFile))
}

func (c CreateManager) createConfigFile(dir string) error {
	tplFile := tplgo.TemplateConfig
	return c.file.Write(dir+"/config.json", []byte(tplFile))
}

func updateTree(fCmd string, t Tree, i int) Tree {
	fc := splitFormulaCommand(fCmd)
	parent := generateParent(fc, i)

	command := funk.Filter(t.Commands, func(command api.Command) bool {
		return command.Usage == fc[i] && command.Parent == parent
	}).([]api.Command)

	if len(fc)-1 == i {
		if len(command) == 0 {
			pathValue := strings.Join(fc, "/")
			fn := fc[len(fc)-1]
			commands := append(t.Commands, api.Command{
				Usage: fn,
				Help:  fmt.Sprintf("%s %s", fc[i-1], fc[i]),
				Formula: api.Formula{
					Path:   pathValue,
					Bin:    fn + "-${so}",
					Config: "config.json",
				},
				Parent: parent,
			})
			t.Commands = commands
			return t
		} else {
			log.Fatal("Command already exist.")
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

	return updateTree(fCmd, t, i+1)
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
