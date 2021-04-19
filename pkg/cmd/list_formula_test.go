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
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewListFormula(t *testing.T) {
	type in struct {
		args            []string
		repoList        formula.Repos
		repoListErr     error
		inputListString string
		inputListErr    error
		tree            formula.Tree
		treeErr         error
		tutorial        rtutorial.Finder
		listFormulaErr  error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				args: []string{},
				repoList: formula.Repos{
					{
						Name: "repoName",
					},
				},
				inputListString: "repoName",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputListString, tt.in.inputListErr)
			repoManagerMock := new(mocks.RepoManager)
			repoManagerMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.repoList, tt.in.repoListErr)
			treeManagerMock := new(mocks.TreeManager)
			treeManagerMock.On("TreeByRepo").Return(tt.in.tree, tt.in.treeErr)

			cmd := NewListFormulaCmd(
				repoManagerMock,
				inputListMock,
				treeManagerMock,
				tt.in.tutorial,
			)
			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()
			assert.Equal(t, tt.want, got)
		})
	}
}
