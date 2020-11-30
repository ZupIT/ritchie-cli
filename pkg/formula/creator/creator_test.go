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

package creator

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	fCmdExists        = "rit add repo"
	fCmdCorrectGo     = "rit scaffold generate test_go"
	fCmdCorrectJava   = "rit scaffold generate test_java"
	fCmdCorrectNode   = "rit scaffold generate test_node"
	fCmdCorrectPython = "rit scaffold generate test_python"
	fCmdCorrectShell  = "rit scaffold generate test_shell"
	langGo            = "go"
	langJava          = "java"
	langNode          = "node"
	langPython        = "python"
	langShell         = "bash shell"
)

func TestCreator(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	api.RitchieHomeDir()

	resultDir := path.Join(os.TempDir(), "/customWorkSpace")
	_ = dirManager.Remove(resultDir)
	_ = dirManager.Create(resultDir)

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

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.CoreCmds, FileReadExisterMock{}, repoProviders, isRootCommand)

	tplM := template.NewManager("../../../testdata", dirManager)

	type in struct {
		formCreate formula.Create
		dir        stream.DirCreateChecker
		file       stream.FileWriteReadExister
		tplM       template.Manager
		helpPath   string
	}

	type out struct {
		err error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "command exists",
			in: in{
				formCreate: formula.Create{
					FormulaCmd: fCmdExists,
					Lang:       langGo,
					Workspace:  formula.Workspace{Dir: resultDir},
					FormulaPath: func() string {
						fp := path.Join(resultDir, "/add/repo")
						_ = dirManager.Remove(fp)
						_ = dirManager.Create(fp)
						return fp
					}(),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: ErrRepeatedCommand,
			},
		},
		{
			name: "command correct-go",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:  fCmdCorrectGo,
					Lang:        langGo,
					Workspace:   formula.Workspace{Dir: resultDir},
					FormulaPath: path.Join(resultDir, "/scaffold/generate/test_go"),
				},
				dir:      dirManager,
				file:     fileManager,
				tplM:     tplM,
				helpPath: "/tmp/customWorkSpace/scaffold/generate/test_go/help.json",
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-java",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:  fCmdCorrectJava,
					Lang:        langJava,
					Workspace:   formula.Workspace{Dir: resultDir},
					FormulaPath: path.Join(resultDir, "/scaffold/generate/test_java"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-node",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:  fCmdCorrectNode,
					Lang:        langNode,
					Workspace:   formula.Workspace{Dir: resultDir},
					FormulaPath: path.Join(resultDir, "/scaffold/generate/test_node"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-python",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:  fCmdCorrectPython,
					Lang:        langPython,
					Workspace:   formula.Workspace{Dir: resultDir},
					FormulaPath: path.Join(resultDir, "/scaffold/generate/test_python"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-shell",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:  fCmdCorrectShell,
					Lang:        langShell,
					Workspace:   formula.Workspace{Dir: resultDir},
					FormulaPath: path.Join(resultDir, "/scaffold/generate/test_shell"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
			},
			out: out{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			creator := NewCreator(treeMan, tt.in.dir, tt.in.file, tt.in.tplM)
			out := tt.out
			got := creator.Create(in.formCreate)
			if out.err == nil {
				require.NoError(t, got)
			} else {
				require.NotNil(t, out.err)
				require.Equal(t, out.err.Error(), out.err.Error())
			}

			if in.helpPath != "" {
				bytes, _ := fileManager.Read(in.helpPath)
				help := formula.Help{}
				_ = json.Unmarshal(bytes, &help)
				errorMsg := "help.json should not be empty"
				expected := "scaffold generate test_go formula"
				require.Equal(t, expected, help.Short, errorMsg)
				require.Equal(t, expected, help.Short, errorMsg)
			}
		})
	}
}

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
