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
	"path"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/modifier"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
	filesToModify      = []string{"README.md", "metadata.json"}
)

type CreateManager struct {
	dir  stream.DirCreateCopyRemoveChecker
	file stream.FileWriteReadExister
}

func NewCreator(
	dir stream.DirCreateCopyRemoveChecker,
	file stream.FileWriteReadExister,
) CreateManager {
	return CreateManager{dir: dir, file: file}
}

func (c CreateManager) Create(cf formula.Create) error {
	if err := c.isValidCmd(cf.FormulaPath); err != nil {
		return err
	}

	if err := c.dir.Create(cf.WorkspacePath); err != nil {
		return err
	}

	fCmdName := cf.FormulaCmdName()

	modifiers := modifier.NewModifiers(cf)
	if err := c.generateFormulaFiles(cf.FormulaPath, cf.Lang, fCmdName, cf.WorkspacePath, modifiers); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) isValidCmd(formulaPath string) error {
	if c.dir.Exists(formulaPath) {
		return ErrRepeatedCommand
	}

	return nil
}

func (c CreateManager) generateFormulaFiles(formulaPath, lang, fCmdName, workSpcPath string, modifiers []modifier.Modifier) error {
	if err := c.dir.Create(formulaPath); err != nil {
		return err
	}

	if err := c.createHelpFiles(fCmdName, workSpcPath); err != nil {
		return err
	}

	if err := c.applyLangTemplate(lang, formulaPath); err != nil {
		return err
	}

	for _, file := range filesToModify {
		if err := c.modifyTemplateFile(formulaPath, file, modifiers); err != nil {
			return err
		}
	}

	return nil
}

func (c CreateManager) applyLangTemplate(lang, formulaPath string) error {
	filePath := path.Join("templates", "languages", lang)
	if err := template.RestoreAssets(formulaPath, filePath); err != nil {
		return err
	}

	oldPath := filepath.Join(formulaPath, filePath)
	if err := c.dir.Copy(oldPath, formulaPath); err != nil {
		return err
	}

	if err := c.dir.Remove(oldPath); err != nil {
		return err
	}

	return nil
}

func (c CreateManager) createHelpFiles(formulaCmdName, workSpacePath string) error {
	dirs := strings.Split(formulaCmdName, " ")
	for i := 0; i < len(dirs); i++ {
		d := dirs[0 : i+1]
		tPath := filepath.Join(workSpacePath, filepath.Join(d...))
		helpPath := filepath.Join(tPath, template.HelpFileName)
		if !c.file.Exists(helpPath) {
			folderName := filepath.Base(tPath)
			tpl := strings.ReplaceAll(template.HelpJson, "{{folderName}}", folderName)
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

func (c CreateManager) modifyTemplateFile(formulaPath, fileName string, modifiers []modifier.Modifier) error {
	filePath := filepath.Join(formulaPath, fileName)
	data, err := c.file.Read(filePath)
	if err != nil {
		return err
	}

	data = modifier.Modify(data, modifiers)

	if err := c.file.Write(filePath, data); err != nil {
		return err
	}

	return nil
}
