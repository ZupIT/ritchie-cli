package formula

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_go"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_java"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_node"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_python"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tpl/tpl_shell"
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

	err = generateTreeJsonFile(c.FormPath, fCmd, lang)
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

func generateTreeJsonFile(formPath, fCmd, lang string) error {
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

	tree, err = updateTree(fCmd, tree, lang, 0)
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
	return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tplFile))
}

func createScripts(dir string) error {
	tplFile := tpl_go.TemplateCopyBinConfig

	err := fileutil.WriteFilePerm(fmt.Sprintf("%s/copy-bin-configs.sh", dir), []byte(tplFile), 0755)
	if err != nil {
		return err
	}

	tplFile = tpl_go.TemplateUnzipBinConfigs

	return fileutil.WriteFilePerm(fmt.Sprintf("%s/unzip-bin-configs.sh", dir), []byte(tplFile), 0755)
}

func createSrcFiles(dir, pkg, lang string) error {
	srcDir := fmt.Sprintf("%s/src", dir)
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
		err = createDockerfile(srcDir, tpl_go.TemplateDockerfile)
		if err != nil {
			return err
		}
		pkgDir := fmt.Sprintf("%s/pkg/%s", srcDir, pkg)
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
		err = createDockerfile(srcDir, tpl_java.TemplateDockerfile)
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
		err = createMainFile(srcDir, pkg, lang)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, lang)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_node.TemplateDockerfile)
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
	case "Python":
		err = createMainFile(srcDir, pkg, lang)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, lang)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_python.TemplateDockerfile)
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
	default:
		err = createMainFile(srcDir, pkg, lang)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, lang)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_shell.TemplateDockerfile)
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
	}

	return nil
}

func createPkgFile(dir, pkg, lang string) error {
	switch lang {
	case "Go":
		tfgo := tpl_go.TemplatePkg
		tfgo = strings.ReplaceAll(tfgo, nameModule, pkg)

		return fileutil.WriteFile(fmt.Sprintf("%s/%s.go", dir, pkg), []byte(tfgo))
	case "Java":
		tfj := tpl_java.TemplateFileJava
		tfj = strings.ReplaceAll(tfj, nameBin, pkg)
		fu := strings.Title(strings.ToLower(pkg))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.java", dir, fu), []byte(tfj))
	case "Node":
		tfn := tpl_node.TemplateFileNode
		tfn = strings.ReplaceAll(tfn, nameBin, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.js", dir, pkg), []byte(tfn))
	case "Python":
		tfp := tpl_python.TemplateFilePython
		tfp = strings.ReplaceAll(tfp, nameBinFirstUpper, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.py", dir, pkg), []byte(tfp))
	default:
		tfs := tpl_shell.TemplateFileShell
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.sh", dir, pkg), []byte(tfs))
	}
}

func createRunTemplate(dir, lang string) error {
	switch lang {
	case "Go":
		return nil
	case "Java":
		tj := tpl_java.TemplateRunTemplate
		return fileutil.WriteFilePerm(fmt.Sprintf("%s/run_template", dir), []byte(tj), 0777)
	case "Node":
		tn := tpl_node.TemplateRunTemplate
		return fileutil.WriteFilePerm(fmt.Sprintf("%s/run_template", dir), []byte(tn), 0777)
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
		return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tplFile))
	case "Java":
		tfj := tpl_java.TemplateMakefile
		tfj = strings.ReplaceAll(tfj, nameBin, name)
		fu := strings.Title(strings.ToLower(name))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)
		return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tfj))
	case "Node":
		tfn := tpl_node.TemplateMakefile
		tfn = strings.ReplaceAll(tfn, nameBin, name)
		err := fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tfn))
		if err != nil {
			return err
		}
		tfpj := tpl_node.TemplatePackageJson
		return fileutil.WriteFile(fmt.Sprintf("%s/package.json", dir), []byte(tfpj))
	case "Python":
		tfp := tpl_python.TemplateMakefile
		tfp = strings.ReplaceAll(tfp, nameBin, name)
		fu := strings.Title(strings.ToLower(name))
		tfp = strings.ReplaceAll(tfp, nameBinFirstUpper, fu)
		return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tfp))
	default:
		tfs := tpl_shell.TemplateMakefile
		tfs = strings.ReplaceAll(tfs, nameBin, name)
		return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tfs))
	}
}

func createDockerfile(dir string, tpl string) error {
	return fileutil.WriteFile(fmt.Sprintf("%s/Dockerfile", dir), []byte(tpl))

}

func createGoModFile(dir, pkg string) error {
	tplFile := tpl_go.TemplateGoMod
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return fileutil.WriteFile(fmt.Sprintf("%s/go.mod", dir), []byte(tplFile))
}

func createMainFile(dir, pkg, lang string) error {
	switch lang {
	case "Go":
		tfgo := tpl_go.TemplateMain
		tfgo = strings.ReplaceAll(tfgo, nameModule, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/main.go", dir), []byte(tfgo))
	case "Java":
		tfj := tpl_java.TemplateMain
		tfj = strings.ReplaceAll(tfj, nameBin, pkg)
		fu := strings.Title(strings.ToLower(pkg))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)
		return fileutil.WriteFile(fmt.Sprintf("%s/Main.java", dir), []byte(tfj))
	case "Node":
		tfn := tpl_node.TemplateIndex
		tfn = strings.ReplaceAll(tfn, nameBin, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/index.js", dir), []byte(tfn))
	case "Python":
		tfp := tpl_python.TemplateMain
		tfp = strings.ReplaceAll(tfp, nameBin, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/main.py", dir), []byte(tfp))
	default:
		tfs := tpl_shell.TemplateMain
		tfs = strings.ReplaceAll(tfs, nameBin, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/main.sh", dir), []byte(tfs))
	}
}

func createConfigFile(dir string) error {
	tplFile := tpl_go.TemplateConfig
	return fileutil.WriteFile(fmt.Sprintf("%s/config.json", dir), []byte(tplFile))
}

func updateTree(fCmd string, t Tree, lang string, i int) (Tree, error) {
	fc := splitFormulaCommand(fCmd)
	parent := generateParent(fc, i)

	command := funk.Filter(t.Commands, func(command api.Command) bool {
		return command.Usage == fc[i] && command.Parent == parent
	}).([]api.Command)

	if len(fc)-1 == i {
		if len(command) == 0 {
			pathValue := strings.Join(fc, "/")
			fn := fc[len(fc)-1]
			var commands []api.Command
			if lang == "Python" {
				commands = append(t.Commands, api.Command{
					Usage: fn,
					Help:  fmt.Sprintf("%s %s", fc[i-1], fc[i]),
					Formula: api.Formula{
						Path:   pathValue,
						Bin:    "main.py",
						LBin:   "main.py",
						MBin:   "main.py",
						WBin:   fmt.Sprintf("%s.bat", fn),
						Bundle: "${so}.zip",
						Config: "config.json",
					},
					Parent: parent,
				})
			} else if lang == "Go" {
				commands = append(t.Commands, api.Command{
					Usage: fn,
					Help:  fmt.Sprintf("%s %s", fc[i-1], fc[i]),
					Formula: api.Formula{
						Path:   pathValue,
						Bin:    fmt.Sprintf("%s-${so}", fn),
						LBin:   fmt.Sprintf("%s-${so}", fn),
						MBin:   fmt.Sprintf("%s-${so}", fn),
						WBin:   fmt.Sprintf("%s-${so}.exe", fn),
						Bundle: "${so}.zip",
						Config: "config.json",
					},
					Parent: parent,
				})
			} else {
				commands = append(t.Commands, api.Command{
					Usage: fn,
					Help:  fmt.Sprintf("%s %s", fc[i-1], fc[i]),
					Formula: api.Formula{
						Path:   pathValue,
						Bin:    fmt.Sprintf("%s.sh", fn),
						LBin:   fmt.Sprintf("%s.sh", fn),
						MBin:   fmt.Sprintf("%s.sh", fn),
						WBin:   fmt.Sprintf("%s.bat", fn),
						Bundle: "${so}.zip",
						Config: "config.json",
					},
					Parent: parent,
				})
			}
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

	return updateTree(fCmd, t, lang, i+1)
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
