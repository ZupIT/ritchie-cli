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

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func Test_listWorkspaceCmd_runFunc(t *testing.T) {
	finderTutorial := rtutorial.NewFinder(os.TempDir(), stream.NewFileManager())
	type in struct {
		WorkspaceLister formula.WorkspaceLister
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				WorkspaceLister: WorkspaceListerCustomMock{
					list: func() (formula.Workspaces, error) {
						return formula.Workspaces{
							"workspace1": "path/to/workspace1",
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run with success with more than 1 workspace",
			in: in{
				WorkspaceLister: WorkspaceListerCustomMock{
					list: func() (formula.Workspaces, error) {
						return formula.Workspaces{
							"workspace1": "path/to/workspace1",
							"workspace2": "path/to/workspace2",
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Return err when list fail",
			in: in{
				WorkspaceLister: WorkspaceListerCustomMock{
					list: func() (formula.Workspaces, error) {
						return nil, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lr := NewListWorkspaceCmd(tt.in.WorkspaceLister, finderTutorial)
			lr.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := lr.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("listWorkspaceCmd_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type WorkspaceListerCustomMock struct {
	list func() (formula.Workspaces, error)
}

func (m WorkspaceListerCustomMock) List() (formula.Workspaces, error) {
	return m.list()
}
