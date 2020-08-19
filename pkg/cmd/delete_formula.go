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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/spf13/cobra"
)

const (
	localRepo     = "/repos/local"
	localTreeJson = "/repos/local/tree.json"
)

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
		treeManager    formula.TreeManager
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
	treeManager formula.TreeManager,
) *cobra.Command {
	d := deleteFormulaCmd{
		userHomeDir,
		ritchieHomeDir,
		workspace,
		directory,
		inBool,
		inText,
		inList,
		treeManager,
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
		workspaces[formula.DefaultWorkspaceName] = defaultWorkspace

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
			return errors.New("Could not find formula")
		}

		if ans, err := d.inBool.Bool("Are you sure you want to delete the formula: rit "+strings.Join(groups, " "), []string{"no", "yes"}); err != nil {
			return err
		} else if !ans {
			return nil
		}

		if err := d.deleteFormula(wspace.Dir, groups, 0); err != nil {
			return err
		}

		ritchieLocalWorkspace := d.ritchieHomeDir + localRepo
		if err := d.deleteFormula(ritchieLocalWorkspace, groups, 0); err != nil {
			return err
		}

		if err := d.deleteFormulaTreeJson(groups); err != nil {
			return err
		}

		prompt.Success("✔ Formula successfully deleted!")

		return nil
	}
}

func (d deleteFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		deleteStdin := deleteFormulaStdin{}

		if err := stdin.ReadJson(os.Stdin, &deleteStdin); err != nil {
			return err
		}

		if err := d.deleteFormula(deleteStdin.Workspace, deleteStdin.Groups, 0); err != nil {
			return err
		}

		ritchieLocalWorkspace := d.ritchieHomeDir + localRepo
		if err := d.deleteFormula(ritchieLocalWorkspace, deleteStdin.Groups, 0); err != nil {
			return err
		}

		if err := d.deleteFormulaTreeJson(deleteStdin.Groups); err != nil {
			return err
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

func (d deleteFormulaCmd) deleteFormula(workspace string, groups []string, index int) error {
	if index == len(groups) {
		err := os.RemoveAll(workspace)
		if err != nil {
			return err
		}

		return nil
	}

	err := d.deleteFormula(workspace+"/"+groups[index], groups, index+1)
	if err != nil {
		return err
	}

	ok, err := canDelete(workspace)
	if err != nil {
		return err
	}

	if ok {
		err := os.RemoveAll(workspace)
		if err != nil {
			return err
		}
	}

	return nil
}

func canDelete(workspace string) (bool, error) {
	files, err := ioutil.ReadDir(workspace)
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

func (d deleteFormulaCmd) deleteFormulaTreeJson(groups []string) error {
	tree, err := d.treeManager.Tree()
	if err != nil {
		return err
	}

	localTree := tree["LOCAL"]
	rootFormulaName := generateRootFormulaName(groups)
	index := findCommandIdTreeJson(localTree, rootFormulaName)
	if index == -1 {
		return errors.New("Could not find formula in tree.json")
	}

	localTree = deleteCommandTreeJson(localTree, index)

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	if err := ioutil.WriteFile(d.ritchieHomeDir+localTreeJson, jsonString, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func generateRootFormulaName(groups []string) string {
	var name = "root"
	for _, group := range groups {
		name += "_" + group
	}

	return name
}

func findCommandIdTreeJson(tree formula.Tree, id string) int {
	for i, command := range tree.Commands {
		if command.Id == id {
			return i
		}
	}

	return -1
}

func canDeleteFromTreeJson(tree formula.Tree, id string) bool {
	for _, command := range tree.Commands {
		if command.Parent == id {
			return false
		}
	}

	return true
}

func deleteCommandTreeJson(tree formula.Tree, index int) formula.Tree {
	if index == -1 {
		return tree
	}

	if canDeleteFromTreeJson(tree, tree.Commands[index].Id) {
		parent := tree.Commands[index].Parent
		tree.Commands = append(tree.Commands[:index], tree.Commands[index+1:]...)

		parentIndex := findCommandIdTreeJson(tree, parent)
		return deleteCommandTreeJson(tree, parentIndex)
	} else {
		return tree
	}
}
