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
	"fmt"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

type CreateManager struct {
	ritHome string
	github  github.Repositories
	dir     stream.DirCreateListCopyRemover
	file    stream.FileWriteCreatorReadExistRemover
}

func NewCreator(
	ritHome string,
	github github.Repositories,
	dir stream.DirCreateListCopyRemover,
	file stream.FileWriteCreatorReadExistRemover,
) CreateManager {
	return CreateManager{
		ritHome: ritHome,
		github:  github,
		dir:     dir,
		file:    file,
	}
}

func (cr CreateManager) Create(repo formula.Repo) error {
	repoInfo := github.NewRepoInfo(repo.Url, repo.Token)
	zipball, err := cr.github.Zipball(repoInfo, repo.Version.String()) // Download zip repository from github
	if err != nil {
		return err
	}

	defer zipball.Close()

	repoPath := path.Join(cr.ritHome, reposDirName, repo.Name.String())
	if err := cr.dir.Remove(repoPath); err != nil { // Remove old repo directory
		return err
	}

	if err := cr.dir.Create(repoPath); err != nil { // Create new repo directory
		return err
	}

	zipFile := path.Join(repoPath, fmt.Sprintf("%s.zip", repo.Name))
	if err := cr.file.Create(zipFile, zipball); err != nil { // Create .zip file inside repo directory
		return err
	}

	if err := streams.Unzip(zipFile, repoPath); err != nil {
		return err
	}

	if err := cr.file.Remove(zipFile); err != nil { // Remove .zip file
		return err
	}

	dirs, err := cr.dir.List(repoPath, true) // Get the directories after unzip the new repo
	if err != nil {
		return err
	}

	src := path.Join(repoPath, dirs[0])                // Get the first directory created by unzip
	if err := cr.dir.Copy(src, repoPath); err != nil { // Copy all formulas inside directory created by unzip to repo path
		return err
	}

	if err := cr.dir.Remove(src); err != nil { // Remove directory created by unzip
		return err
	}

	return nil
}
