package tree

import (
	"fmt"
	"path"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	root        = "root"
	rootPattern = "root_%s"
	helpFile    = "help.txt"
	configFile  = "config.json"
)

type GeneratorManager struct {
	dir  stream.DirLister
	file stream.FileReadExister
}

func NewGenerator(dir stream.DirLister, file stream.FileReadExister) GeneratorManager {
	return GeneratorManager{dir: dir, file: file}
}

func (ge GeneratorManager) Generate(repoPath string) (formula.Tree, error) {
	dirs, err := ge.dir.List(repoPath, false)
	if err != nil {
		return formula.Tree{}, err
	}

	commands := api.Commands{}
	for _, dir := range dirs { // Generate root commands
		formulaPath := path.Join(repoPath, dir)
		helpFilePath := path.Join(formulaPath, helpFile)
		if !ge.file.Exists(helpFilePath) { // Ignore folders without help.txt
			continue
		}

		helpFile, err := ge.file.Read(helpFilePath)
		if err != nil {
			return formula.Tree{}, err
		}

		cmd := api.Command{
			Id:     fmt.Sprintf(rootPattern, dir),
			Parent: root,
			Usage:  dir,
			Help:   strings.TrimSuffix(string(helpFile), "\n"),
		}

		commands = append(commands, cmd)

		commands, err = ge.subCommands(formulaPath, cmd, commands)
		if err != nil {
			return formula.Tree{}, err
		}
	}

	return formula.Tree{Commands: commands}, nil
}

// subCommands generates the sub-commands for the tree.
// if success returns an api.Commands and a nil error
// if error returns an empty api.Commands and an error not empty
func (ge GeneratorManager) subCommands(dirPath string, cmd api.Command, cmds api.Commands) (api.Commands, error) {
	dirs, err := ge.dir.List(dirPath, false)
	if err != nil {
		return cmds, err
	}

	for _, dir := range dirs {
		if dir == "src" {
			return cmds, nil
		}

		if dir == "bin" { // Ignore /bin directory
			continue
		}

		formulaPath := path.Join(dirPath, dir)
		helpFilePath := path.Join(formulaPath, helpFile)
		var helpFile []byte
		if ge.file.Exists(helpFilePath) { // Check if help.txt exist
			helpFile, err = ge.file.Read(helpFilePath)
			if err != nil {
				return nil, err
			}
		}

		cmd := api.Command{
			Id:     fmt.Sprintf("%s_%s", cmd.Id, dir),
			Parent: cmd.Id,
			Usage:  dir,
			Help:  strings.TrimSuffix(string(helpFile), "\n"),
		}

		configFilePath := path.Join(formulaPath, configFile)
		if ge.file.Exists(configFilePath) { // Case config.json exists, set cmd.Formula as true
			cmd.Formula = true
		}

		cmds = append(cmds, cmd)

		cmds, err = ge.subCommands(formulaPath, cmd, cmds)
		if err != nil {
			return nil, err
		}
	}

	return cmds, nil
}
