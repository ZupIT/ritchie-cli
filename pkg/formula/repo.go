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

package formula

import "github.com/ZupIT/ritchie-cli/pkg/git"

const RepoCommonsName = RepoName("commons")

type Repo struct {
	Provider RepoProvider `json:"provider"`
	Name     RepoName     `json:"name"`
	Version  RepoVersion  `json:"version"`
	Url      string       `json:"url"`
	Token    string       `json:"token,omitempty"`
	Priority int          `json:"priority"`
}

type Repos []Repo

func (r Repos) Len() int {
	return len(r)
}

func (r Repos) Less(i, j int) bool {
	return r[i].Priority < r[j].Priority
}

func (r Repos) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type RepoName string

func (r RepoName) String() string {
	return string(r)
}

type RepoVersion string

func (r RepoVersion) String() string {
	return string(r)
}

type RepoProvider string

func (r RepoProvider) String() string {
	return string(r)
}

type Git struct {
	Repos       git.Repositories
	NewRepoInfo func(url string, token string) git.RepoInfo
}

type RepoProviders map[RepoProvider]Git

func NewRepoProviders() RepoProviders {
	return RepoProviders{}
}

func (re RepoProviders) Add(provider RepoProvider, git Git) {
	re[provider] = git
}

func (re RepoProviders) Resolve(provider RepoProvider) Git {
	return re[provider]
}

func (re RepoProviders) List() []string {
	var providers []string
	for provider := range re {
		providers = append(providers, provider.String())
	}

	return providers
}

type RepositoryAdder interface {
	Add(repo Repo) error
}

type RepositoryLister interface {
	List() (Repos, error)
}

type RepositoryListerLocal interface {
	ListLocal() (RepoName, error)
}

type RepositoryUpdater interface {
	Update(name RepoName, version RepoVersion) error
}

type RepositoryDeleter interface {
	Delete(name RepoName) error
}

type RepositoryLocalDeleter interface {
	DeleteLocal() error
}

type RepositoryPrioritySetter interface {
	SetPriority(name RepoName, priority int) error
}

type RepositoryCreator interface {
	Create(repo Repo) error
}

type RepositoryAddLister interface {
	RepositoryAdder
	RepositoryLister
}

type RepositoryDelLister interface {
	RepositoryDeleter
	RepositoryLister
}

type RepositoryListCreator interface {
	RepositoryLister
	RepositoryCreator
}

type RepositoryListUpdater interface {
	RepositoryLister
	RepositoryUpdater
}
