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
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type DeleteLocalManager struct {
	ritHome string
	dir     stream.DirRemover
}

func NewLocalDeleter(ritHome string, dir stream.DirRemover) DeleteLocalManager {
	return DeleteLocalManager{
		ritHome: ritHome,
		dir:     dir,
	}
}

func (dm DeleteLocalManager) DeleteLocal() error {
	if err := dm.deleteRepoDir("local"); err != nil {
		return err
	}
	return nil
}

func (dm DeleteLocalManager) deleteRepoDir(repoName formula.RepoName) error {
	repoPath := filepath.Join(dm.ritHome, reposDirName, repoName.String())
	if err := dm.dir.Remove(repoPath); err != nil {
		return err
	}
	return nil
}
