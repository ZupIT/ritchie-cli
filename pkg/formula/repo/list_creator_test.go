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
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/git/github"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewListCreator(t *testing.T) {

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

	type in struct {
		repoList   formula.RepositoryLister
		repoCreate formula.RepositoryCreator
	}
	tests := []struct {
		name string
		in   in
		want formula.RepositoryListCreator
	}{
		{
			name: "Build with success",
			in: in{
				repoList:   repoList,
				repoCreate: repoCreator,
			},
			want: ListCreateManager{
				RepositoryLister:  repoList,
				RepositoryCreator: repoCreator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListCreator(tt.in.repoList, tt.in.repoCreate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListCreator() = %v, want %v", got, tt.want)
			}
		})
	}
}
