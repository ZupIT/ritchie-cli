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

package repo

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/git/github"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewCreateWriteListDetailDeleter(t *testing.T) {
	ritHome := os.TempDir()
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	repoList := NewLister(ritHome, fileManager)
	repoProviders := formula.RepoProviders{
		"Github": formula.Git{
			Repos:       github.NewRepoManager(http.DefaultClient),
			NewRepoInfo: github.NewRepoInfo,
		},
	}
	repoCreator := NewCreator(ritHome, repoProviders, dirManager, fileManager)
	repoWrite := NewWriter(ritHome, fileManager)
	repoDetail := NewDetail(repoProviders)
	repoListWriter := NewListWriter(repoList, repoWrite)
	repoDeleter := NewDeleter(ritHome, repoListWriter, dirManager)

	want := CreateWriteListDetailDeleter{
		RepositoryLister:  repoList,
		RepositoryCreator: repoCreator,
		RepositoryWriter:  repoWrite,
		RepositoryDetail:  repoDetail,
		RepositoryDeleter: repoDeleter,
	}

	got := NewCreateWriteListDetailDeleter(repoList, repoCreator, repoWrite, repoDetail, repoDeleter)
	assert.Equal(t, want, got)

}
