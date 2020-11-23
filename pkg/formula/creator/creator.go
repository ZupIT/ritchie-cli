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

package creator

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/modifier"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
)

type CreateManager struct {
	treeManager tree.Manager
	dir         stream.DirCreateChecker
	file        stream.FileWriteReadExister
	tplM        template.Manager
}

func NewCreator(
	tm tree.Manager,
	dir stream.DirCreateChecker,
	file stream.FileWriteReadExister,
	tplM template.Manager,
) CreateManager {
	return CreateManager{treeManager: tm, dir: dir, file: file, tplM: tplM}
}

func (c CreateManager) Create(cf formula.Create) error {
	if err := c.isValidCmd(cf.FormulaPath); err != nil {
		return err
	}

	if err := c.tplM.Validate(); err != nil {
		return err
	}

	if err := c.dir.Create(cf.Workspace.Dir); err != nil {
		return err
	}

	fCmdName := cf.FormulaCmdName()

	modifiers := modifier.NewModifiers(cf)
	if err := c.generateFormulaFiles(cf.FormulaPath, cf.Lang, fCmdName, cf.Workspace.Dir, modifiers); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) isValidCmd(fPath string) error {
	if c.dir.Exists(fPath) {
		return ErrRepeatedCommand
	}

	return nil
}

func (c CreateManager) generateFormulaFiles(
	fPath,
	lang,
	fCmdName,
	workSpcPath string,
	modifiers []modifier.Modifier,
) error {

	if err := c.dir.Create(fPath); err != nil {
		return err
	}

	if err := c.createHelpFiles(fCmdName, workSpcPath); err != nil {
		return err
	}

	if err := c.applyLangTemplate(lang, fPath, workSpcPath, modifiers); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) applyLangTemplate(lang, formulaPath, workspacePath string, modifiers []modifier.Modifier) error {

	tFiles, err := c.tplM.LangTemplateFiles(lang)
	if err != nil {
		return err
	}

	for _, f := range tFiles {
		if f.IsDir {
			newPath, err := c.tplM.ResolverNewPath(f.Path, formulaPath, lang, workspacePath)
			if err != nil {
				return err
			}
			err = c.dir.Create(newPath)
			if err != nil {
				return err
			}
		} else {
			newPath, err := c.tplM.ResolverNewPath(f.Path, formulaPath, lang, workspacePath)
			if err != nil {
				return err
			}
			if c.file.Exists(newPath) {
				continue
			}
			tpl, err := c.file.Read(f.Path)
			if err != nil {
				return err
			}
			newDir, _ := filepath.Split(newPath)
			err = c.dir.Create(newDir)
			if err != nil {
				return err
			}
			tpl = modifier.Modify(tpl, modifiers)
			err = c.file.Write(newPath, tpl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c CreateManager) createHelpFiles(formulaCmdName, workSpacePath string) error {
	dirs := strings.Split(formulaCmdName, " ")
	var commands string
	for i := 0; i < len(dirs); i++ {
		d := dirs[0 : i+1]
		tPath := filepath.Join(workSpacePath, filepath.Join(d...))
		helpPath := filepath.Join(tPath, template.HelpFileName)
		if !c.file.Exists(helpPath) {
			folderName := filepath.Base(tPath)
			commands += folderName + " "
			tpl := strings.ReplaceAll(template.HelpJson, "{{command}}", commands)
			help := formula.Help{}

			err := json.Unmarshal([]byte(tpl), &help)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(help, "", "\t")
			if err != nil {
				return err
			}
			err = c.file.Write(helpPath, b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
