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

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
)

func TestNewSingleInitCmd(t *testing.T) {
	cmd := NewInitCmd(defaultRepoAdderMock, defaultGitRepositoryMock, TutorialFinderMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewInitCmd got %v", cmd)
		return
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func Test_initCmd_runPrompt(t *testing.T) {
	type fields struct {
		repo formula.RepositoryAdder
		git  git.Repositories
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run With Success",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
			},
			wantErr: false,
		},
		{
			name: "Fail when call git.LatestTag",
			fields: fields{
				repo: defaultRepoAdderMock,
				git: GitRepositoryMock{
					latestTag: func(info git.RepoInfo) (git.Tag, error) {
						return git.Tag{}, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fail when call repo.Add",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return errors.New("some error")
					},
				},
				git: defaultGitRepositoryMock,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewInitCmd(tt.fields.repo, tt.fields.git, TutorialFinderMock{})
			o.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := o.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
