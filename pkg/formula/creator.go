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

	"github.com/thoas/go-funk"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_go"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_java"
)

type CreateManager struct {
	FormPath    string
	treeManager TreeManager
}

func NewCreator(homePath string, tm TreeManager) CreateManager {
	return CreateManager{FormPath: fmt.Sprintf(FormCreatePathPattern, homePath), treeManager: tm}
}

func (c CreateManager) Create(fCmd, lang string) (CreateManager, error) {
	_ = fileutil.CreateDirIfNotExists(c.FormPath, os.ModePerm)
	trees, err := c.treeManager.Tree()
	if err != nil {
		return CreateManager{}, err
	}

	err = verifyCommand(fCmd, trees)
	if err != nil {
		return CreateManager{}, err
	}

	err = generateTreeJsonFile(c.FormPath, fCmd)
	if err != nil {
		return CreateManager{}, err
	}

	if fileutil.Exists(fmt.Sprintf(TreeCreatePathPattern, c.FormPath)) && (fileutil.Exists(fmt.Sprintf("%s/%s", c.FormPath, Makefile))) {
		err = generateFormulaFiles(c.FormPath, fCmd, lang, false)
		if err != nil {
			return CreateManager{}, err
		}
	} else {
		err = generateFormulaFiles(c.FormPath, fCmd, lang, true)
		if err != nil {
			return CreateManager{}, err
		}
	}

	return c, nil
}

func generateFormulaFiles(formPath, fCmd, lang string, new bool) error {
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
	err = createSrcFiles(dir, d[len(d)-1], lang)
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

	tree, err = updateTree(fCmd, tree, 0)
	if err != nil {
		return err
	}
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
	tplFile := tpl_go.TemplateMakefileMain

	tplFile = strings.ReplaceAll(tplFile, "{{formName}}", strings.ToUpper(name))
	tplFile = strings.ReplaceAll(string(tplFile), "{{formPath}}", dirForm)

	err := createScripts(dir)
	if err != nil {
		return err
	}

	return fileutil.WriteFile(dir+"/Makefile", []byte(tplFile))
}

func createScripts(dir string) error {
	tplFile := tpl_go.TemplateCopyBinConfig

	err := fileutil.WriteFilePerm(dir+"/copy-bin-configs.sh", []byte(tplFile), 0755)
	if err != nil {
		return err
	}

	tplFile = tpl_go.TemplateUnzipBinConfigs

	return fileutil.WriteFilePerm(dir+"/unzip-bin-configs.sh", []byte(tplFile), 0755)
}

func createSrcFiles(dir, pkg, lang string) error {
	srcDir := dir + "/src"
	err := fileutil.CreateDirIfNotExists(srcDir, os.ModePerm)
	if err != nil {
		return err
	}
	switch lang {
	case "Go":
		err = createMainFile(srcDir, pkg, lang)
		if err != nil {
			return err
		}
		err = createGoModFile(srcDir, pkg)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, lang)
		if err != nil {
			return err
		}
		pkgDir := srcDir + "/pkg/" + pkg
		err = fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = createPkgFile(pkgDir, pkg, lang)
		if err != nil {
			return err
		}
	case "Java":
		err = createMainFile(srcDir, pkg, lang)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, lang)
		if err != nil {
			return err
		}
		err = createRunTemplate(srcDir, lang)
		if err != nil {
			return err
		}
		pkgDir := fmt.Sprintf("%s/%s", srcDir, pkg)
		err = fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = createPkgFile(pkgDir, pkg, lang)
		if err != nil {
			return err
		}
	case "Node":
		log.Println("Formula in Node")
	case "Python":
		log.Println("Formula in Python")
	default:
		log.Println("Formula in Shell")
	}

	return nil
}

func createPkgFile(dir, pkg, lang string) error {
	switch lang {
	case "Go":
		tfgo := tpl_go.TemplatePkg
		tfgo = strings.ReplaceAll(tfgo, nameModule, pkg)

		return fileutil.WriteFile(dir+"/"+pkg+".go", []byte(tfgo))
	case "Java":
		tfj := tpl_java.TemplateFileJava
		tfj = strings.ReplaceAll(tfj, nameBin, pkg)
		fu := strings.Title(strings.ToLower(pkg))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)
		return fileutil.WriteFile(dir+"/"+fu+".java", []byte(tfj))
	case "Node":
	case "Python":
	default:

	}
	return nil
}

func createRunTemplate(dir, lang string) error {
	switch lang {
	case "Go":
		return nil
	case "Java":
		tplFile := tpl_java.TemplateRunTemplate
		return fileutil.WriteFilePerm(dir+"/run_template", []byte(tplFile), 0777)
	case "Node":
	case "Python":
	default:

	}
	return nil
}

func createMakefileForm(dir string, name, pathName, lang string) error {
	switch lang {
	case "Go":
		tplFile := tpl_go.TemplateMakefile
		tplFile = strings.ReplaceAll(tplFile, "{{name}}", name)
		tplFile = strings.ReplaceAll(tplFile, "{{form-path}}", pathName)

		return fileutil.WriteFile(dir+"/Makefile", []byte(tplFile))
	case "Java":
		tfj := tpl_java.TemplateMakefile
		tfj = strings.ReplaceAll(tfj, nameBin, name)
		fu := strings.Title(strings.ToLower(name))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)

		return fileutil.WriteFile(dir+"/Makefile", []byte(tfj))
	case "Node":
	case "Python":
	default:

	}
	return nil
}

func createGoModFile(dir, pkg string) error {
	tplFile := tpl_go.TemplateGoMod
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return fileutil.WriteFile(dir+"/go.mod", []byte(tplFile))
}

func createMainFile(dir, pkg, lang string) error {
	switch lang {
	case "Go":
		tfgo := tpl_go.TemplateMain
		tfgo = strings.ReplaceAll(tfgo, nameModule, pkg)
		return fileutil.WriteFile(dir+"/main.go", []byte(tfgo))
	case "Java":
		tfj := tpl_java.TemplateMain
		tfj = strings.ReplaceAll(tfj, nameBin, pkg)
		fu := strings.Title(strings.ToLower(pkg))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)

		return fileutil.WriteFile(dir+"/Main.java", []byte(tfj))
	case "Node":
	case "Python":
	default:

	}
	return nil
}

func createConfigFile(dir string) error {
	tplFile := tpl_go.TemplateConfig
	return fileutil.WriteFile(dir+"/config.json", []byte(tplFile))
}

func updateTree(fCmd string, t Tree, i int) (Tree, error) {
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
					Bin:    fn + ".sh",
					LBin:   fn + ".sh",
					MBin:   fn + ".sh",
					WBin:   fn + ".bat",
					Bundle: "${so}.zip",
					Config: "config.json",
				},
				Parent: parent,
			})
			t.Commands = commands
			return t, nil
		} else {
			return Tree{}, errors.New("Command already exist ")
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
