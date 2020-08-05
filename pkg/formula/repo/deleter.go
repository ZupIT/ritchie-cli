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
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type DeleteManager struct {
	ritHome string
	file    stream.FileWriteReadExister
	dir     stream.DirRemover
}

func NewDeleter(ritHome string, file stream.FileWriteReadExister, dir stream.DirRemover) DeleteManager {
	return DeleteManager{
		ritHome: ritHome,
		file:    file,
		dir:     dir,
	}
}

func (dm DeleteManager) Delete(repoName formula.RepoName) error {
	if err := dm.deleteRepoDir(repoName); err != nil {
		return err
	}
	if err := dm.deleteFromReposFile(repoName); err != nil {
		return err
	}
	return nil
}

func (dm DeleteManager) deleteRepoDir(repoName formula.RepoName) error {
	repoPath := filepath.Join(dm.ritHome, reposDirName, repoName.String())
	if err := dm.dir.Remove(repoPath); err != nil {
		return err
	}
	return nil
}

func (dm DeleteManager) deleteFromReposFile(repoName formula.RepoName) error {
	repos := formula.Repos{}

	repoFilePath := filepath.Join(dm.ritHome, reposDirName, reposFileName)
	file, err := dm.file.Read(repoFilePath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &repos); err != nil {
		return err
	}

	var idx int
	for i := range repos {
		if repos[i].Name == repoName {
			idx = i
			break
		}
	}
	repos = append(repos[:idx], repos[idx+1:]...)

	newFile, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err = dm.file.Write(repoFilePath, newFile); err != nil {
		return err
	}

	return nil
}
