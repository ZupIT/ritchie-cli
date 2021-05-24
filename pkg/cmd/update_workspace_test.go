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
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateWorkspaceRun(t *testing.T) {

	workspaceTest := &formula.Workspace{
		Name: "Test",
		Dir:  "/Users/dennis/workspaces",
	}

	workspaceTest2 := &formula.Workspace{
		Name: "Test2",
		Dir:  "/Users/dennis/workspaces",
	}

	var tests = []struct {
		name            string
		args            []string
		wspaceList      formula.Workspaces
		inputList       string
		wspaceListErr   error
		wspaceUpdateErr error
		inputListErr    error
		want            error
	}{
		{
			name: "update workspace success",
			args: []string{},
			wspaceList: formula.Workspaces{
				workspaceTest.Name:  workspaceTest.Dir,
				workspaceTest2.Name: workspaceTest2.Dir,
			},
			inputList: convertToInputFormat(workspaceTest),
		},
		{
			name:          "update workspace error (listing workspace)",
			args:          []string{},
			wspaceListErr: errors.New("error to list workspace"),
			want:          errors.New("error to list workspace"),
		},
		{
			name:       "update workspace error (empty workspace)",
			args:       []string{},
			wspaceList: formula.Workspaces{},
			want:       ErrEmptyWorkspace,
		},
		{
			name: "update workspace error (input list)",
			args: []string{},
			wspaceList: formula.Workspaces{
				workspaceTest.Name:  workspaceTest.Dir,
				workspaceTest2.Name: workspaceTest2.Dir,
			},
			inputListErr: errors.New("error to input list"),
			want:         errors.New("error to input list"),
		},
		{
			name: "update workspace error (update workspace)",
			args: []string{},
			wspaceList: formula.Workspaces{
				workspaceTest.Name:  workspaceTest.Dir,
				workspaceTest2.Name: workspaceTest2.Dir,
			},
			inputList:       convertToInputFormat(workspaceTest),
			wspaceUpdateErr: errors.New("error to update workspace"),
			want:            errors.New("error to update workspace"),
		},
		{
			name: "update workspace error (empty flag name)",
			args: []string{"--name="},
			want: errors.New("please provide a value for 'name'"),
		},
		{
			name: "update workspace error (wrong workspace name)",
			args: []string{"--name=Unexpected"},
			wspaceList: formula.Workspaces{
				workspaceTest.Name:  workspaceTest.Dir,
				workspaceTest2.Name: workspaceTest2.Dir,
			},
			want: errors.New("no workspace found with this name"),
		},
		{
			name: "update workspace success (input flag)",
			args: []string{"--name=Test"},
			wspaceList: formula.Workspaces{
				workspaceTest.Name:  workspaceTest.Dir,
				workspaceTest2.Name: workspaceTest2.Dir,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceMock := new(mocks.WorkspaceForm)
			workspaceMock.On("List").Return(tt.wspaceList, tt.wspaceListErr)
			workspaceMock.On("Update", mock.Anything).Return(tt.wspaceUpdateErr)
			inListMock := new(mocks.InputListMock)
			inListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.inputList, tt.inputListErr)

			cmd := NewUpdateWorkspaceCmd(
				workspaceMock,
				inListMock,
			)
			cmd.SetArgs(tt.args)

			got := cmd.Execute()
			assert.Equal(t, tt.want, got)
		})
	}
}

func convertToInputFormat(workspace *formula.Workspace) string {
	return workspace.Name + " (" + workspace.Dir + ")"
}
