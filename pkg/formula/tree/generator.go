/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tree

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	root        = "root"
	rootPattern = "root_%s"
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
		formulaPath := filepath.Join(repoPath, dir)
		helpFilePath := filepath.Join(formulaPath, template.HelpFileName)
		if !ge.file.Exists(helpFilePath) { // Ignore folders without help.txt
			continue
		}

		helpFile, err := ge.file.Read(helpFilePath)
		if err != nil {
			return formula.Tree{}, err
		}
		help := formula.Help{}
		err = json.Unmarshal(helpFile, &help)
		if err != nil {
			return formula.Tree{}, err
		}

		cmd := api.Command{
			Id:       fmt.Sprintf(rootPattern, dir),
			Parent:   root,
			Usage:    dir,
			Help:     help.Short,
			LongHelp: help.Long,
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

		formulaPath := filepath.Join(dirPath, dir)
		helpFilePath := filepath.Join(formulaPath, template.HelpFileName)
		help := formula.Help{}
		if ge.file.Exists(helpFilePath) { // Check if help.txt exist
			helpFile, err := ge.file.Read(helpFilePath)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(helpFile, &help)
			if err != nil {
				return nil, err
			}

		}

		cmd := api.Command{
			Id:       fmt.Sprintf("%s_%s", cmd.Id, dir),
			Parent:   cmd.Id,
			Usage:    dir,
			Help:     help.Short,
			LongHelp: help.Long,
		}

		configFilePath := filepath.Join(formulaPath, configFile)
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
