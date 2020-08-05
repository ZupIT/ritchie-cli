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
	"fmt"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type UpdateManager struct {
	ritHome string
	repo    formula.RepositoryListCreator
	tree    formula.TreeGenerator
	file    stream.FileWriter
}

func NewUpdater(ritHome string,
	repo formula.RepositoryListCreator,
	tree formula.TreeGenerator,
	file stream.FileWriter,
) UpdateManager {
	return UpdateManager{
		ritHome: ritHome,
		repo:    repo,
		tree:    tree,
		file:    file,
	}
}

func (up UpdateManager) Update(name formula.RepoName, version formula.RepoVersion) error {
	repos, err := up.repo.List()
	if err != nil {
		return err
	}

	var repo *formula.Repo
	for i := range repos {
		if name == repos[i].Name {
			repo = &repos[i]
			break
		}
	}

	if repo == nil {
		return fmt.Errorf("repository name %q was not found", name)
	}

	repo.Version = version

	if err := up.repo.Create(*repo); err != nil {
		return err
	}

	repoFilePath := filepath.Join(up.ritHome, reposDirName, reposFileName)
	if err := up.saveRepo(repoFilePath, repos); err != nil {
		return err
	}

	repoPath := filepath.Join(up.ritHome, reposDirName, name.String())
	tree, err := up.tree.Generate(repoPath)
	if err != nil {
		return err
	}

	treeFilePath := filepath.Join(repoPath, "tree.json")
	bytes, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := up.file.Write(treeFilePath, bytes); err != nil {
		return err
	}

	return nil
}

func (up UpdateManager) saveRepo(repoPath string, repos formula.Repos) error {
	bytes, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err := up.file.Write(repoPath, bytes); err != nil {
		return err
	}

	return nil
}
