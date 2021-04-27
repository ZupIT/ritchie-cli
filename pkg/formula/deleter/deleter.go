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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implier.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package deleter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrRepeatedCommand = prompt.NewError("this command already exists")
)

type DeleteManager struct {
	dir            stream.DirCreateCheckerCopy
	file           stream.FileWriteRemover
	treeGen        formula.TreeGenerator
	ritchieHomeDir string
}

func NewDeleter(
	dir stream.DirCreateCheckerCopy,
	file stream.FileWriteRemover,
	treeGen formula.TreeGenerator,
	ritchieHomeDir string,
) DeleteManager {
	return DeleteManager{dir, file, treeGen, ritchieHomeDir}
}

func (d *DeleteManager) Delete(fr formula.Delete) error {
	if err := d.deleteFormula(fr.Workspace.Dir, fr.GroupsFormula, 0); err != nil {
		return err
	}

	ritchieLocalWorkspace := filepath.Join(d.ritchieHomeDir, "repos", "local-default")
	if d.formulaExistsInWorkspace(ritchieLocalWorkspace, fr.GroupsFormula) {
		if err := d.deleteFormula(ritchieLocalWorkspace, fr.GroupsFormula, 0); err != nil {
			return err
		}

		if err := d.recreateTreeJSON(ritchieLocalWorkspace); err != nil {
			return err
		}
	}

	return nil
}

func (d *DeleteManager) deleteFormula(path string, groups []string, index int) error {
	if index == len(groups) {
		nested, err := nestedFormula(path)
		if err != nil {
			return err
		}

		if nested {
			return d.safeRemoveFormula(path)
		}

		return os.RemoveAll(path)
	}

	newPath := filepath.Join(path, groups[index])
	if err := d.deleteFormula(newPath, groups, index+1); err != nil {
		return err
	}

	if index == 0 {
		return nil
	}

	ok, err := canDelete(path)
	if err != nil {
		return err
	}

	if ok {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func (d *DeleteManager) safeRemoveFormula(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() && (file.Name() == "src" || file.Name() == "bin") {
			pathToDelete := filepath.Join(path, file.Name())
			if err := os.RemoveAll(pathToDelete); err != nil {
				return err
			}
		} else if !file.IsDir() {
			pathToDelete := filepath.Join(path, file.Name())
			if err := d.file.Remove(pathToDelete); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *DeleteManager) formulaExistsInWorkspace(path string, groups []string) bool {
	for _, group := range groups {
		path = filepath.Join(path, group)
	}

	return d.dir.Exists(path)
}

func (d *DeleteManager) recreateTreeJSON(workspace string) error {
	localTree, err := d.treeGen.Generate(workspace)
	if err != nil {
		return err
	}

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	pathLocalTreeJSON := filepath.Join(d.ritchieHomeDir, "repos", "local-default", "tree.json")
	if err = d.file.Write(pathLocalTreeJSON, jsonString); err != nil {
		return err
	}

	return nil
}

func nestedFormula(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "src" && file.Name() != "bin" {
			return true, nil
		}
	}

	return false, nil
}

func canDelete(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.IsDir() {
			return false, nil
		}
	}

	return true, nil
}
