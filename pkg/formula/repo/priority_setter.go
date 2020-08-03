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
	"sort"

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

	for i := range repos {
		if repoName == repos[i].Name {
			repos[i].Priority = priority
		} else if repos[i].Priority >= priority {
			repos[i].Priority++
		}
	}

	sort.Sort(repos)

	bytes, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err := sm.file.Write(repoPath, bytes); err != nil {
		return err
	}

	return nil
}
