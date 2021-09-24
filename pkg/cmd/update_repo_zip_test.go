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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestUpdateRepoZipCmd(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		existentRepo formula.Repo
		repoName     string
		repoURL      string
		repoVersion  string
		iListErr     error
		iULRErr      error
		repoUpdErr   error
		repoListErr  error
		want         error
	}{
		{
			name:         "run with success",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote", Version: "1.0.0"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.1.0",
			want:         nil,
		},
		{
			name:         "fail when the version field is invalid",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote", Version: "1.0.0"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			want:         ErrSameVersion,
		},
		{
			name:         "fail when repo.List return err",
			existentRepo: formula.Repo{},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			repoListErr:  errors.New("could not list repositories"),
			want:         errors.New("could not list repositories"),
		},
		{
			name:     "return error when the version field is empty",
			repoName: "zipremote",
			repoURL:  "https://provider.com/download-repo/repo-1.0.0.zip",
			want:     ErrVersionNotEmpty,
		},
		{
			name:         "return error for the url field",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			repoName:     "zipremote",
			repoURL:      "https://provider.com/download-repo/repo-1.0.0.zip",
			repoVersion:  "1.0.0",
			iULRErr:      errors.New("input url error"),
			want:         errors.New("input url error"),
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
			name:         "fail flags when repo.List return err",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			args:         []string{"--name=zipremote", "--version=1.0.0", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip"},
			repoListErr:  errors.New("could not list repositories"),
			want:         errors.New("could not list repositories"),
		},
		{
			name:         "fail flags with empty name",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			args:         []string{"--version=1.0.0", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip"},
			want:         errors.New(missingFlagText(nameFlagName)),
		},
		{
			name:         "fail flags with empty repo url",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote"},
			args:         []string{"--version=1.0.0", "--name=zipremote"},
			want:         errors.New(missingFlagText(repoUrlFlagName)),
		},
		{
			name:         "fail flags with empty version",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote", Version: "1.0.0"},
			args:         []string{"--name=zipremote", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip"},
			want:         errors.New(missingFlagText(versionFlagName)),
		},
		{
			name:         "fail flags with same version",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote", Version: "1.0.0"},
			args:         []string{"--name=zipremote", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip", "--version=1.0.0"},
			want:         ErrSameVersion,
		},
		{
			name:         "fail flags with the name field is invalid",
			existentRepo: formula.Repo{},
			args:         []string{"--name=zip", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip", "--version=1.0.0"},
			repoUpdErr:   errors.New("repository name zip was not found"),
			want:         errors.New("repository name zip was not found"),
		},
		{
			name:         "success flags",
			existentRepo: formula.Repo{Provider: "ZipRemote", Name: "zipremote", Version: "1.0.0"},
			args:         []string{"--name=zipremote", "--repoUrl=https://provider.com/download-repo/repo-1.0.0.zip", "--version=1.1.0"},
			want:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			questionTypeNewVersion := fmt.Sprintf("Type your new version for %q:", tt.repoName)

			inputURLMock := new(mocks.InputURLMock)
			inputURLMock.On("URL", mock.Anything, mock.Anything).Return(tt.repoURL, tt.iULRErr)
			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.repoName, tt.iListErr)
			inputTextValidatorMock := new(mocks.InputTextValidatorMock)
			inputTextValidatorMock.On("Text", questionSelectARepo, mock.Anything).Return(tt.repoName, nil)
			inputTextValidatorMock.On("Text", questionTypeNewVersion, mock.Anything).Return(tt.repoVersion, nil)
			repoListUpdaterMock := new(mocks.UpdaterMock)
			repoListUpdaterMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(tt.repoUpdErr)
			repoListUpdaterMock.On("List").Return(formula.Repos{tt.existentRepo}, tt.repoListErr)

			cmd := NewUpdateRepoZipCmd(
				repoListUpdaterMock,
				inputTextValidatorMock,
				inputURLMock,
				inputListMock,
			)

			cmd.SetArgs(tt.args)
			got := cmd.Execute()

			assert.Equal(t, tt.want, got)
		})
	}
}
