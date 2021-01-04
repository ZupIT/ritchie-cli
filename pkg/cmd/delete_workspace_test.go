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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestDeleteWorkspaces(t *testing.T) {
	type in struct {
		wspaceList      formula.Workspaces
		wspaceListErr   error
		wspaceDeleteErr error
		dirExist        bool
		inputList       string
		inputListErr    error
		inputBool       bool
		inputBoolErr    error
		repoDeleteErr   error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				wspaceList: formula.Workspaces{
					"local-commons": "/home/user/commons",
				},
				dirExist:  true,
				inputList: "local-commons",
				inputBool: true,
			},
		},
		{
			name: "error to list workspace",
			in: in{
				wspaceListErr: errors.New("error to list workspace"),
			},
			want: errors.New("error to list workspace"),
		},
		{
			name: "error empty workspace",
			in: in{
				wspaceList: formula.Workspaces{},
			},
			want: ErrEmptyWorkspaces,
		},
		{
			name: "error to input list",
			in: in{
				wspaceList: formula.Workspaces{
					"local-commons": "/home/user/commons",
				},
				dirExist:     true,
				inputListErr: errors.New("error to input list"),
				inputBool:    true,
			},
			want: errors.New("error to input list"),
		},
		{
			name: "error to accept to delete selected workspace",
			in: in{
				wspaceList: formula.Workspaces{
					"local-commons": "/home/user/commons",
				},
				dirExist:     true,
				inputList:    "local-commons",
				inputBoolErr: errors.New("error to accept"),
			},
			want: errors.New("error to accept"),
		},
		{
			name: "not accept to delete selected workspace",
			in: in{
				wspaceList: formula.Workspaces{
					"local-commons": "/home/user/commons",
				},
				dirExist:  true,
				inputList: "local-commons",
				inputBool: false,
			},
		},
		{
			name: "error to delete repo",
			in: in{
				wspaceList: formula.Workspaces{
					"local-commons": "/home/user/commons",
				},
				dirExist:      true,
				inputList:     "local-commons",
				inputBool:     true,
				repoDeleteErr: errors.New("error to delete repo"),
			},
			want: errors.New("error to delete repo"),
		},
		{
			name: "error to delete workspace",
			in: in{
				wspaceList: formula.Workspaces{
					"local-commons": "/home/user/commons",
				},
				dirExist:        true,
				inputList:       "local-commons",
				inputBool:       true,
				wspaceDeleteErr: errors.New("error to delete workspace"),
			},
			want: errors.New("error to delete workspace"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceMock := new(mocks.WorkspaceForm)
			workspaceMock.On("List").Return(tt.in.wspaceList, tt.in.wspaceListErr)
			workspaceMock.On("Delete", mock.Anything).Return(tt.in.wspaceDeleteErr)
			dirMock := new(mocks.DirManager)
			dirMock.On("Exists", mock.Anything).Return(tt.in.dirExist)
			inListMock := new(mocks.InputListMock)
			inListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputList, tt.in.inputListErr)
			inBoolMock := new(mocks.InputBoolMock)
			inBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputBool, tt.in.inputBoolErr)
			repoManagerMock := new(mocks.RepoManager)
			repoManagerMock.On("Delete", mock.Anything).Return(tt.in.repoDeleteErr)

			cmd := NewDeleteWorkspaceCmd(
				os.TempDir(),
				workspaceMock,
				repoManagerMock,
				dirMock,
				inListMock,
				inBoolMock,
			)

			got := cmd.Execute()
			assert.Equal(t, tt.want, got)
		})
	}
}
