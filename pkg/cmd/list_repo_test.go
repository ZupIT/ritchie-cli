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

func Test_listRepoCmd_runFunc(t *testing.T) {
	finderTutorial := rtutorial.NewFinder(os.TempDir(), stream.NewFileManager())
	type in struct {
		RepositoryLister formula.RepositoryLister
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "someRepo1",
							},
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Run with success with more than 1 repo",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return formula.Repos{
							{
								Name: "someRepo1",
							},
							{
								Name: "someRepo2",
							},
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Return err when list fail",
			in: in{
				RepositoryLister: RepositoryListerCustomMock{
					list: func() (formula.Repos, error) {
						return nil, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lr := NewListRepoCmd(tt.in.RepositoryLister, finderTutorial)
			lr.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := lr.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("setCredentialCmd_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type RepositoryListerCustomMock struct {
	list func() (formula.Repos, error)
}

func (m RepositoryListerCustomMock) List() (formula.Repos, error) {
	return m.list()
}
