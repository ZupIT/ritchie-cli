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

package autocomplete

import (
	"io"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
)

type repoListerMock struct{}

func (repoListerMock) List() (formula.Repos, error) {
	return formula.Repos{}, nil
}

type FileReadExisterMock struct{}

func (m FileReadExisterMock) Read(path string) ([]byte, error) {
	return []byte("some data"), nil
}

func (m FileReadExisterMock) Exists(path string) bool {
	return false
}

type GitRepositoryMock struct {
	zipball   func(info git.RepoInfo, version string) (io.ReadCloser, error)
	tags      func(info git.RepoInfo) (git.Tags, error)
	latestTag func(info git.RepoInfo) (git.Tag, error)
}

func (m GitRepositoryMock) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	return m.zipball(info, version)
}

func (m GitRepositoryMock) Tags(info git.RepoInfo) (git.Tags, error) {
	return m.tags(info)
}

func (m GitRepositoryMock) LatestTag(info git.RepoInfo) (git.Tag, error) {
	return m.latestTag(info)
}

func TestGenerate(t *testing.T) {
	type in struct {
		shell ShellName
	}

	type out struct {
		err error
	}

	var defaultGitRepositoryMock = GitRepositoryMock{
		latestTag: func(info git.RepoInfo) (git.Tag, error) {
			return git.Tag{}, nil
		},
		tags: func(info git.RepoInfo) (git.Tags, error) {
			return git.Tags{git.Tag{Name: "1.0.0"}}, nil
		},
		zipball: func(info git.RepoInfo, version string) (io.ReadCloser, error) {
			return nil, nil
		},
	}
	repoProviders := formula.NewRepoProviders()
	repoProviders.Add("Github", formula.Git{Repos: defaultGitRepositoryMock, NewRepoInfo: github.NewRepoInfo})
	isRootCommand := false

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.Commands{}, FileReadExisterMock{}, repoProviders, isRootCommand)
	autocomplete := NewGenerator(treeMan)

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "autocomplete bash",
			in: &in{
				shell: bash,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete zsh",
			in: &in{
				shell: zsh,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete fish",
			in: &in{
				shell: fish,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete powerShell",
			in: &in{
				shell: powerShell,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete error",
			in: &in{
				shell: "err",
			},
			out: &out{
				err: ErrNotSupported,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := autocomplete.Generate(tt.in.shell, &cobra.Command{})

			if err != tt.out.err {
				t.Errorf("Generator(%s) got %v, want %v", tt.name, err, tt.out.err)
			}

			if tt.out.err == nil && got == "" {
				t.Errorf("Generator(%s) autocomplete is empty", tt.name)
			}
		})
	}
}
