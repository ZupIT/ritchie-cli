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

var ErrMakefileNotFound = errors.New("makefile not found")
var ErrTreeJsonNotFound = errors.New("tree.json not found")

type CreateManager struct {
	FormPath    string
	treeManager TreeManager
}

func NewCreator(homePath string, tm TreeManager) CreateManager {
	return CreateManager{FormPath: fmt.Sprintf(FormCreatePathPattern, homePath), treeManager: tm}
}

func (c CreateManager) Create(fCmd, lang, localRepoDir string) (CreateManager, error) {
	_ = fileutil.CreateDirIfNotExists(c.FormPath, os.ModePerm)

	if localRepoDir != "" {

		if !existsTreeJson(localRepoDir) && existsMakefile(localRepoDir){
			return CreateManager{}, ErrTreeJsonNotFound
		}
		if !existsMakefile(localRepoDir) && existsTreeJson(localRepoDir) {
			return CreateManager{}, ErrMakefileNotFound
		}

		c.FormPath = localRepoDir
	}

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

	if existsMakefile(c.FormPath) && existsTreeJson(c.FormPath) {
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

func existsTreeJson(formPath string) bool {
	treePath := fmt.Sprintf(TreeCreatePathPattern, formPath)
		return fileutil.Exists(treePath)
}

func existsMakefile(formPath string) bool {
	makefilePath := fmt.Sprintf(MakefileCreatePathPattern, formPath, Makefile)
		return fileutil.Exists(makefilePath)
}

func generateFormulaFiles(formPath, fCmd, lang string, new bool) error {
	d := strings.Split(fCmd, " ")
	dirForm := strings.Join(d[1:], "/")
	formulaName := fmt.Sprintf("%s_%s", d[len(d)-2], d[len(d)-1])

	var dir string
	if new {
		dir = fmt.Sprintf("%s/%s", formPath, dirForm)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		err = createMakefileMain(formPath, dirForm, formulaName)
		if err != nil {
			return err
		}

	} else {
		dir = fmt.Sprintf("%s/%s", formPath, dirForm)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		err = changeMakefileMain(formPath, fCmd, formulaName)
		if err != nil {
			return err
		}
	}
	err := createConfigFile(dir)
	if err != nil {
		return err
	}
	err = createSrcFiles(dir, formulaName, lang)
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

	if len(s) <= 2 {
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

func changeMakefileMain(formPath, fCmd, fName string) error {
	d := strings.Split(fCmd, " ")
	dir := fmt.Sprintf("%s/%s", formPath, Makefile)
	tplFile, err := fileutil.ReadFile(dir)
	if err != nil {
		return err
	}
	variable := strings.ToUpper(fName) + "=" + strings.Join(d[1:], "/")
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
	tplFile := tpl_go.MakefileMain

	tplFile = strings.ReplaceAll(tplFile, "{{formName}}", strings.ToUpper(name))
	tplFile = strings.ReplaceAll(tplFile, "{{formPath}}", dirForm)

	err := createScripts(dir)
	if err != nil {
		return err
	}
	return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tplFile))
}

func createScripts(dir string) error {
	tplFile := tpl_go.CopyBinConfig

	err := fileutil.WriteFilePerm(fmt.Sprintf("%s/copy-bin-configs.sh", dir), []byte(tplFile), 0755)
	if err != nil {
		return err
	}

	tplFile = tpl_go.UnzipBinConfigs

	return fileutil.WriteFilePerm(fmt.Sprintf("%s/unzip-bin-configs.sh", dir), []byte(tplFile), 0755)
}

func createSrcFiles(dir, pkg, lang string) error {
	srcDir := fmt.Sprintf("%s/src", dir)
	pkgDir := fmt.Sprintf("%s/%s", srcDir, pkg)
	err := fileutil.CreateDirIfNotExists(srcDir, os.ModePerm)
	if err != nil {
		return err
	}
	switch lang {
	case "Go":
		err = createMainFile(srcDir, pkg, tpl_go.Main, "go", "main", false)
		if err != nil {
			return err
		}
		err = createGoModFile(srcDir, pkg)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, tpl_go.Makefile, true)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_go.Dockerfile)
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
		err = createMainFile(srcDir, pkg, tpl_java.Main, "java", "Main", true)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, tpl_java.Makefile, false)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_java.Dockerfile)
		if err != nil {
			return err
		}
		err = createRunTemplate(srcDir, tpl_java.RunTemplate)
		if err != nil {
			return err
		}
		err = createPkgDir(pkgDir)
		if err != nil {
			return err
		}
		err = createPkgFile(pkgDir, pkg, lang)
		if err != nil {
			return err
		}
	case "Node":
		err = createMainFile(srcDir, pkg, tpl_node.Index, "js", "index", false)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, tpl_node.Makefile, false)
		if err != nil {
			return err
		}
		err = createPackageJson(srcDir, tpl_node.PackageJson)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_node.Dockerfile)
		if err != nil {
			return err
		}
		err = createRunTemplate(srcDir, tpl_node.RunTemplate)
		if err != nil {
			return err
		}
		err = createPkgDir(pkgDir)
		if err != nil {
			return err
		}
		err = createPkgFile(pkgDir, pkg, lang)
		if err != nil {
			return err
		}
	case "Python":
		err = createMainFile(srcDir, pkg, tpl_python.Main, "py", "main", false)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, tpl_python.Makefile, false)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_python.Dockerfile)
		if err != nil {
			return err
		}
		err = createPkgDir(pkgDir)
		if err != nil {
			return err
		}
		err = createPkgFile(pkgDir, pkg, lang)
		if err != nil {
			return err
		}
	default:
		err = createMainFile(srcDir, pkg, tpl_shell.Main, "sh", "main", false)
		if err != nil {
			return err
		}
		err = createMakefileForm(srcDir, pkg, dir, tpl_shell.Makefile, false)
		if err != nil {
			return err
		}
		err = createDockerfile(srcDir, tpl_shell.Dockerfile)
		if err != nil {
			return err
		}
		err = createPkgDir(pkgDir)
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

func createPkgDir(pkgDir string) error {
	return fileutil.CreateDirIfNotExists(pkgDir, os.ModePerm)
}

func createPkgFile(dir, pkg, lang string) error {
	switch lang {
	case "Go":
		tfgo := strings.ReplaceAll(tpl_go.Pkg, nameModule, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.go", dir, pkg), []byte(tfgo))
	case "Java":
		tfj := strings.ReplaceAll(tpl_java.File, nameBin, pkg)
		fu := strings.Title(strings.ToLower(pkg))
		tfj = strings.ReplaceAll(tfj, nameBinFirstUpper, fu)
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.java", dir, fu), []byte(tfj))
	case "Node":
		tfn := tpl_node.File
		tfn = strings.ReplaceAll(tfn, nameBin, pkg)
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.js", dir, pkg), []byte(tfn))
	case "Python":
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.py", dir, pkg), []byte(tpl_python.File))
	default:
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.sh", dir, pkg), []byte(tpl_shell.File))
	}
}

func createRunTemplate(dir, tpl string) error {
	return fileutil.WriteFilePerm(fmt.Sprintf("%s/run_template", dir), []byte(tpl), 0777)
}

func createMakefileForm(dir, name, pathName, tpl string, compiled bool) error {
	if compiled {
		tpl = strings.ReplaceAll(tpl, "{{name}}", name)
		tpl = strings.ReplaceAll(tpl, "{{form-path}}", pathName)
		return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tpl))
	}
	tpl = strings.ReplaceAll(tpl, nameBin, name)
	return fileutil.WriteFile(fmt.Sprintf("%s/Makefile", dir), []byte(tpl))
}

func createDockerfile(dir, tpl string) error {
	return fileutil.WriteFile(fmt.Sprintf("%s/Dockerfile", dir), []byte(tpl))
}

func createGoModFile(dir, pkg string) error {
	tplFile := tpl_go.GoMod
	tplFile = strings.ReplaceAll(tplFile, nameModule, pkg)
	return fileutil.WriteFile(fmt.Sprintf("%s/go.mod", dir), []byte(tplFile))
}

func createMainFile(dir, pkg, tpl, fileFormat, startFile string, uc bool) error {
	if uc {
		tpl = strings.ReplaceAll(tpl, nameBin, pkg)
		tpl = strings.ReplaceAll(tpl, nameBinFirstUpper, strings.Title(strings.ToLower(pkg)))
		return fileutil.WriteFile(fmt.Sprintf("%s/%s.%s", dir, startFile, fileFormat), []byte(tpl))
	}

	tpl = strings.ReplaceAll(tpl, nameBin, pkg)
	return fileutil.WriteFile(fmt.Sprintf("%s/%s.%s", dir, startFile, fileFormat), []byte(tpl))
}

func createConfigFile(dir string) error {
	tplFile := tpl_go.Config
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

func createPackageJson(dir, tpl string) error {
	return fileutil.WriteFile(fmt.Sprintf("%s/package.json", dir), []byte(tpl))
}
