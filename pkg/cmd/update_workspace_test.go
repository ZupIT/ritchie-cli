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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateWorkspaceRun(t *testing.T) {
	// someError := errors.New("some error")

	workspaceTest := &formula.Workspace{
		Name: "Teste",
		Dir:  "/Users/dennis/workspaces",
	}

	workspaceTest2 := &formula.Workspace{
		Name: "Teste2",
		Dir:  "/Users/dennis/workspaces",
	}

	type in struct {
		args            []string
		wspaceList      formula.Workspaces
		inputList       string
		wspaceListErr   error
		wspaceUpdateErr error
		inputListErr    error
	}
	var tests = []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				args: []string{},
				wspaceList: formula.Workspaces{
					workspaceTest.Name:  workspaceTest.Dir,
					workspaceTest2.Name: workspaceTest2.Dir,
				},
				inputList: workspaceTest.Name,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceMock := new(mocks.WorkspaceForm)
			workspaceMock.On("List").Return(tt.in.wspaceList, tt.in.wspaceListErr)
			workspaceMock.On("Update", mock.Anything).Return(tt.in.wspaceUpdateErr)
			inListMock := new(mocks.InputListMock)
			inListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.wspaceList, tt.in.inputListErr)

			cmd := NewUpdateWorkspaceCmd(
				workspaceMock,
				inListMock,
			)
			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()
			assert.Equal(t, tt.want, got)
		})
	}
}
