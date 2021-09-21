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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestAddRepoZipCmd(t *testing.T) {
	home := filepath.Join(os.TempDir(), "rit-add-repo-zip")
	ritHome := filepath.Join(home, ".rit")

	defer os.RemoveAll(home)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	reposPath := filepath.Join(ritHome, "repos")
	repoPathLocalDefault := filepath.Join(reposPath, "local-default")
	repoPathWS := filepath.Join(home, "ritchie-formulas-local")

	workspaces := formula.Workspaces{}
	workspaces["Default"] = repoPathWS

	tests := []struct {
		name         string
		args         []string
		existentRepo formula.Repo
		repoName     string
		repoURL      string
		repoVersion  string
		choose       bool
		tutorial     rtutorial.TutorialHolder
		iListErr     error
		iAddErr      error
		iULRErr      error
		iTutorialErr error
		want         error
	}{
		{
			name:        "run with success",
			repoName:    "zipremote",
			repoURL:     "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion: "1.0.0",
			want:        nil,
		},
		{
			name:         "run with success when user add a new commons",
			existentRepo: formula.Repo{Provider: "Github", Name: formula.RepoCommonsName},
			repoName:     "commons",
			repoURL:      "https://github.com/ZupIT/ritchie-formulas/archive/refs/tags/2.16.2.zip",
			repoVersion:  "2.16.2",
			choose:       true,
			want:         nil,
		},
		{
			name:         "run with success when user add a new commons",
			existentRepo: formula.Repo{Provider: "Github", Name: formula.RepoCommonsName},
			repoName:     "commons",
			repoURL:      "https://github.com/ZupIT/ritchie-formulas/archive/refs/tags/2.16.2.zip",
			repoVersion:  "2.16.2",
			want:         nil,
		},
		{
			name:         "run with success when user add a repo existent",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			want:         nil,
		},
		{
			name:         "return nil when user add a new commons incorrectly",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			want:         nil,
		},
		{
			name:         "fail when repo.Add return err",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			iAddErr:      errors.New("could not add repository"),
			want:         errors.New("could not add repository"),
		},
		{
			name:        "return error when the field repository name is empty",
			repoURL:     "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion: "1.0.0",
			want:        ErrRepoNameNotEmpty,
		},
		{
			name:     "return error when the version field is empty",
			repoName: "zipremote",
			repoURL:  "https://provider.com/download-repo/repo-1.0.0.zip",
			want:     ErrVersionNotEmpty,
		},
		{
			name:        "return error for the url field",
			repoName:    "zipremote",
			repoURL:     "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion: "1.0.0",
			iULRErr:     errors.New("input url error"),
			want:        errors.New("input url error"),
		},
		{
			name:         "tutorial status enabled",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			tutorial:     rtutorial.TutorialHolder{Current: ""},
			want:         nil,
		},
		{
			name:         "return error when tutorial.Find fail",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			iTutorialErr: errors.New("tutorial find error"),
			want:         errors.New("tutorial find error"),
		},
		{
			name:         "fail when repo.List return err",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			iListErr:     errors.New("could not list repository"),
			want:         errors.New("could not list repository"),
		},
		{
			name: "fail flags with empty name",
			args: []string{"--version=1.0.0", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip"},
			want: errors.New(missingFlagText(nameFlagName)),
		},
		{
			name: "fail flags with empty repo url",
			args: []string{"--version=1.0.0", "--name=zipremote"},
			want: errors.New(missingFlagText(repoUrlFlagName)),
		},
		{
			name: "fail flags with empty version",
			args: []string{"--name=zipremote", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip"},
			want: errors.New(missingFlagText(versionFlagName)),
		},
		{
			name: "success flags",
			args: []string{"--name=zipremote", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip", "--version=1.0.0"},
			want: nil,
		},
	}
	checkerManager := tree.NewChecker(treeMock{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(home)

			_ = dirManager.Remove(ritHome)

			createSaved := func(path string) {
				_ = dirManager.Remove(path)
				_ = dirManager.Create(path)
			}

			createSaved(repoPathLocalDefault)
			createSaved(repoPathWS)

			zipFile := filepath.Join("..", "..", "testdata", "ritchie-formulas-test.zip")
			zipRepositories := filepath.Join("..", "..", "testdata", "repositories.zip")
			_ = streams.Unzip(zipRepositories, reposPath)
			_ = streams.Unzip(zipFile, repoPathLocalDefault)
			_ = streams.Unzip(zipFile, repoPathWS)

			createTree(repoPathWS, repoPathLocalDefault, treeGen, fileManager)
			setWorkspace(workspaces, ritHome)

			inputURLMock := new(mocks.InputURLMock)
			inputURLMock.On("URL", mock.Anything, mock.Anything).Return(tt.repoURL, tt.iULRErr)
			inputBoolMock := new(mocks.InputBoolMock)
			inputBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.choose, nil)
			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", "Repository name:", mock.Anything).Return(tt.repoName, nil)
			inputTextValidatorMock.On("Text", "Version:", mock.Anything).Return(tt.repoVersion, nil)
			tutorialFindMock := new(mocks.TutorialFindSetterMock)
			tutorialFindMock.On("Find").Return(tt.tutorial, tt.iTutorialErr)
			repoListerAdderMock := new(mocks.RepoManager)
			repoListerAdderMock.On("Add", mock.Anything).Return(tt.iAddErr)
			repoListerAdderMock.On("List").Return(formula.Repos{tt.existentRepo}, tt.iListErr)

			cmd := NewAddRepoZipCmd(
				repoListerAdderMock,
				inputTextValidatorMock,
				inputURLMock,
				inputBoolMock,
				tutorialFindMock,
				checkerManager,
			)

			cmd.SetArgs(tt.args)
			got := cmd.Execute()

			assert.Equal(t, tt.want, got)
		})
	}
}
