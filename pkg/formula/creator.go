package formula

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tplgo"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CreateManager struct {
	formPath    string
	treeManager TreeManager
}

func NewCreator(homePath string, tm TreeManager) CreateManager {
	return CreateManager{formPath: fmt.Sprintf(FormCreatePathPattern, homePath), treeManager: tm}
}

func (c CreateManager) Create(fCmd string) error {
	_ = fileutil.CreateDirIfNotExists(c.formPath, os.ModePerm)
	trees, err := c.treeManager.Tree()
	if err != nil {
		return err
	}

	err = verifyCommand(fCmd, trees)
	if err != nil {
		return err
	}

	if fileutil.Exists(fmt.Sprintf(TreeCreatePathPattern, c.formPath)) && (fileutil.Exists(fmt.Sprintf("%s/%s", c.formPath, Makefile))) {
		generateFormulaFiles(c.formPath, fCmd, false)

	} else {
		generateFormulaFiles(c.formPath, fCmd, true)
	}
	err = generateTreeJsonFile(c.formPath, fCmd)
	if err != nil {
		return err
	}
	log.Println("Formula successfully created!")
	log.Printf("Your formula is in %s", c.formPath)
	return nil
}

func generateFormulaFiles(formPath, fCmd string, new bool) error {
	d := strings.Split(fCmd, " ")
	dirForm := strings.Join(d[1:], "/")

	var dir string
	if new {
		dir = fmt.Sprintf("%s/%s", formPath, dirForm)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		err = createMakefileMain(formPath, dirForm, d[len(d)-1])
		if err != nil {
			return err
		}

	} else {
		dir = fmt.Sprintf("%s/%s", formPath, dirForm)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		err = changeMakefileMain(formPath, fCmd, d[len(d)-1])
		if err != nil {
			return err
		}
	}
	err := createConfigFile(dir)
	if err != nil {
		return err
	}
	err = createSrcFiles(dir, d[len(d)-1])
	if err != nil {
		return err
	}
	return nil
}

func generateTreeJsonFile(formPath, fCmd string) error {
	tree := Tree{Commands: []api.Command{}}
	dir := fmt.Sprintf(localTreeFile, formPath)
	jsonFile, err := fileutil.ReadFile(dir)
	if err != nil {
		if err := fileutil.CreateDirIfNotExists(filepath.Dir(dir), 0755); err != nil {
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

	return fileutil.WriteFile(dir, prettyJSON.Bytes())
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

func changeMakefileMain(formPath string, fCmd, fName string) error {
	d := strings.Split(fCmd, " ")
	dir := fmt.Sprintf("%s/%s", formPath, Makefile)
	tplFile, err := fileutil.ReadFile(dir)
	if err != nil {
		return err
	}
	variable := strings.ToUpper(d[len(d)-1]) + "=" + strings.Join(d[1:], "/")
	tplFile = []byte(strings.ReplaceAll(string(tplFile), "\nFORMULAS=", "\n"+variable+"\nFORMULAS="))
	formulas := formulaValue(tplFile)
	tplFile = []byte(strings.ReplaceAll(string(tplFile), formulas, formulas+" $("+strings.ToUpper(fName)+")"))

	err = fileutil.WriteFile(dir, tplFile)
	if err != nil {
		return err
	}

	return nil
}

func formulaValue(file []byte) string {
	fileStr := string(file)
	return strings.Split(strings.Split(fileStr, "FORMULAS=")[1], "\n")[0]
}

func createMakefileMain(dir, dirForm, name string) error {
	tplFile := tplgo.TemplateMakefileMain

	tplFile = strings.ReplaceAll(tplFile, "{{formName}}", strings.ToUpper(name))
	tplFile = strings.ReplaceAll(string(tplFile), "{{formPath}}", dirForm)

	err := createScripts(dir)
	if err != nil {
		return err
	}

	return fileutil.WriteFile(dir+"/Makefile", []byte(tplFile))
}

func createScripts(dir string) error {
	tplFile := tplgo.TemplateCopyBinConfig

	err := fileutil.WriteFilePerm(dir+"/copy-bin-configs.sh", []byte(tplFile), 0755)
	if err != nil {
		return err
	}

	tplFile = tplgo.TemplateUnzipBinConfigs

	return fileutil.WriteFilePerm(dir+"/unzip-bin-configs.sh", []byte(tplFile), 0755)
}

func createSrcFiles(dir, pkg string) error {
	srcDir := dir + "/src"
	err := fileutil.CreateDirIfNotExists(srcDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = createMainFile(srcDir, pkg)
	if err != nil {
		return err
	}
	err = createGoModFile(srcDir, pkg)
	if err != nil {
		return err
	}
	err = createMakefileForm(srcDir, pkg, dir)
	if err != nil {
		return err
	}
	pkgDir := srcDir + "/pkg/" + pkg
	err = fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = createPkgFile(pkgDir, pkg)
	if err != nil {
		return err
	}

	return nil
}

func createPkgFile(dir, pkg string) error {
	tplFile := tplgo.TemplatePkg
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)

	return fileutil.WriteFile(dir+"/"+pkg+".go", []byte(tplFile))
}

func createMakefileForm(dir string, name, pathName string) error {
	tplFile := tplgo.TemplateMakefile
	tplFile = strings.ReplaceAll(tplFile, "{{name}}", name)
	tplFile = strings.ReplaceAll(tplFile, "{{form-path}}", pathName)

	return fileutil.WriteFile(dir+"/Makefile", []byte(tplFile))
}

func createGoModFile(dir, pkg string) error {
	tplFile := tplgo.TemplateGoMod
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return fileutil.WriteFile(dir+"/go.mod", []byte(tplFile))
}

func createMainFile(dir, pkg string) error {
	tplFile := tplgo.TemplateMain
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return fileutil.WriteFile(dir+"/main.go", []byte(tplFile))
}

func createConfigFile(dir string) error {
	tplFile := tplgo.TemplateConfig
	return fileutil.WriteFile(dir+"/config.json", []byte(tplFile))
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
