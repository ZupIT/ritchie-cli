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
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	repositoryDoNotExistError = "there is no repositories yet"
)

type SetPriorityManager struct {
	ritHome string
	file    stream.FileWriteReadExister
}

func NewPrioritySetter(ritHome string, file stream.FileWriteReadExister) SetPriorityManager {
	return SetPriorityManager{
		ritHome: ritHome,
		file:    file,
	}
}

func (sm SetPriorityManager) SetPriority(repoName formula.RepoName, priority int) error {
	var repos formula.Repos
	repoPath := filepath.Join(sm.ritHome, reposDirName, reposFileName)
	if !sm.file.Exists(repoPath) {
		return errors.New(repositoryDoNotExistError)
	}
	read, err := sm.file.Read(repoPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(read, &repos); err != nil {
		return err
	}

	repos = movePosition(repos, repoName, priority)

	bytes, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err := sm.file.Write(repoPath, bytes); err != nil {
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
	for i = index; i > priority; i-- { // Move repos to back
		r := repos[i-1]
		r.Priority = i
		repos[i] = r
	}

	for i = index; i < priority; i++ { // Move repos to front
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
