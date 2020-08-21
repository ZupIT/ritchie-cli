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

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const msgFormulaNotFound = "Could not find formula"

var ErrCouldNotFindFormula = errors.New(msgFormulaNotFound)

type (
	deleteFormulaStdin struct {
		Workspace string   `json:"workspace"`
		Groups    []string `json:"groups"`
	}

	deleteFormulaCmd struct {
		userHomeDir    string
		ritchieHomeDir string
		workspace      formula.WorkspaceAddListValidator
		directory      stream.DirListChecker
		inBool         prompt.InputBool
		inText         prompt.InputText
		inList         prompt.InputList
		treeGen        formula.TreeGenerator
	}
)

func NewDeleteFormulaCmd(
	userHomeDir string,
	ritchieHomeDir string,
	workspace formula.WorkspaceAddListValidator,
	directory stream.DirListChecker,
	inBool prompt.InputBool,
	inText prompt.InputText,
	inList prompt.InputList,
	treeGen formula.TreeGenerator,
) *cobra.Command {
	d := deleteFormulaCmd{
		userHomeDir,
		ritchieHomeDir,
		workspace,
		directory,
		inBool,
		inText,
		inList,
		treeGen,
	}

	cmd := &cobra.Command{
		Use:     "formula",
		Short:   "Delete specific formula",
		Example: "rit delete formula",
		RunE:    RunFuncE(d.runStdin(), d.runPrompt()),
	}

	return cmd
}

func (d deleteFormulaCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		workspaces, err := d.workspace.List()
		if err != nil {
			return err
		}

		defaultWorkspace := filepath.Join(d.userHomeDir, formula.DefaultWorkspaceDir)
		if d.directory.Exists(defaultWorkspace) {
			workspaces[formula.DefaultWorkspaceName] = defaultWorkspace
		}

		wspace, err := FormulaWorkspaceInput(workspaces, d.inList, d.inText)
		if err != nil {
			return err
		}

		if wspace.Dir != defaultWorkspace {
			if err := d.workspace.Validate(wspace); err != nil {
				return err
			}

			if err := d.workspace.Add(wspace); err != nil {
				return err
			}
		}

		groups, err := d.readFormulas(wspace.Dir)
		if err != nil {
			return err
		}

		if groups == nil {
			return ErrCouldNotFindFormula
		}

		question := fmt.Sprintf("Are you sure you want to delete the formula: rit %s", strings.Join(groups, " "))
		if ans, err := d.inBool.Bool(question, []string{"no", "yes"}); err != nil {
			return err
		} else if !ans {
			return nil
		}

		// Delete formula on user workspace
		if err := d.deleteFormula(wspace.Dir, groups, 0); err != nil {
			return err
		}

		ritchieLocalWorkspace := filepath.Join(d.ritchieHomeDir, "repos", "local")
		if d.formulaExistsInWorkspace(ritchieLocalWorkspace, groups) {
			if err := d.deleteFormula(ritchieLocalWorkspace, groups, 0); err != nil {
				return err
			}

			if err := d.recriateTreeJson(ritchieLocalWorkspace); err != nil {
				return err
			}
		}

		prompt.Success("✔ Formula successfully deleted!")

		return nil
	}
}

func (d deleteFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		deleteStdin := deleteFormulaStdin{}

		if err := stdin.ReadJson(cmd.InOrStdin(), &deleteStdin); err != nil {
			return err
		}

		// Delete formula on user workspace
		if err := d.deleteFormula(deleteStdin.Workspace, deleteStdin.Groups, 0); err != nil {
			return err
		}

		ritchieLocalWorkspace := filepath.Join(d.ritchieHomeDir, "repos", "local")
		if d.formulaExistsInWorkspace(ritchieLocalWorkspace, deleteStdin.Groups) {
			if err := d.deleteFormula(ritchieLocalWorkspace, deleteStdin.Groups, 0); err != nil {
				return err
			}

			if err := d.recriateTreeJson(ritchieLocalWorkspace); err != nil {
				return err
			}
		}

		return nil
	}
}

func (d deleteFormulaCmd) readFormulas(dir string) ([]string, error) {
	dirs, err := d.directory.List(dir, false)
	if err != nil {
		return nil, err
	}

	dirs = sliceutil.Remove(dirs, docsDir)
	var groups []string
	if isFormula(dirs) {
		return groups, nil
	}

	selected, err := d.inList.List("Select a formula or group: ", dirs)
	if err != nil {
		return nil, err
	}

	var aux []string
	aux, err = d.readFormulas(filepath.Join(dir, selected))
	if err != nil {
		return nil, err
	}

	aux = append([]string{selected}, aux...)
	groups = append(groups, aux...)

	return groups, nil
}

func (d deleteFormulaCmd) deleteFormula(path string, groups []string, index int) error {
	if index == len(groups) {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}

		return nil
	}

	newPath := filepath.Join(path, groups[index])
	err := d.deleteFormula(newPath, groups, index+1)
	if err != nil {
		return err
	} else if index == 0 {
		return nil
	}

	ok, err := canDelete(path)
	if err != nil {
		return err
	}

	if ok {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d deleteFormulaCmd) recriateTreeJson(workspace string) error {
	localTree, err := d.treeGen.Generate(workspace)
	if err != nil {
		return err
	}

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	if err := ioutil.WriteFile(filepath.Join(d.ritchieHomeDir, "repos", "local", "tree.json"), jsonString, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (d deleteFormulaCmd) formulaExistsInWorkspace(workspace string, groups []string) bool {
	var path = workspace

	for _, group := range groups {
		path = filepath.Join(path, group)
	}

	return d.directory.Exists(path)
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
