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
	"errors"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const repositoryDoNotExistError = "there is no repositories yet"

type SetPriorityManager struct {
	repo formula.RepositoryListWriter
}

func NewPrioritySetter(repo formula.RepositoryListWriter) SetPriorityManager {
	return SetPriorityManager{repo: repo}
}

func (sm SetPriorityManager) SetPriority(repoName formula.RepoName, priority int) error {
	repos, err := sm.repo.List()
	if err != nil {
		return err
	}

	if repos.Len() <= 0 {
		return errors.New(repositoryDoNotExistError)
	}

	repos = movePosition(repos, repoName, priority)

	if err := sm.repo.Write(repos); err != nil {
		return err
	}

	return nil
}

func movePosition(repos formula.Repos, repoName formula.RepoName, priority int) formula.Repos {
	priority = isValidPriority(priority, repos)
	index := 0
	var repo formula.Repo
	for i := range repos {
		if repoName == repos[i].Name {
			repo = repos[i]
			index = i
			break
		}
	}

	var i int
	for i = index; i > priority; i-- { // Move repos to tail
		r := repos[i-1]
		r.Priority = i
		repos[i] = r
	}

	for i = index; i < priority; i++ { // Move repos to head
		r := repos[i+1]
		r.Priority = i
		repos[i] = r
	}

	repo.Priority = priority
	repos[priority] = repo

	return repos
}

func isValidPriority(priority int, repos formula.Repos) int {
	if priority >= repos.Len() {
		return repos.Len() - 1
	}

	if priority < 0 {
		return 0
	}

	return priority
}
