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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
)

const ZipRemoteProvider = "ZipRemote"

var ErrLocalRepo = errors.New("local repository cannot be updated")

type UpdateManager struct {
	ritHome string
	repo    formula.RepositoryCreateWriteListDetailDeleter
	tree    formula.TreeGenerator
}

func NewUpdater(
	ritHome string,
	repo formula.RepositoryCreateWriteListDetailDeleter,
	tree formula.TreeGenerator,
) UpdateManager {
	return UpdateManager{
		ritHome: ritHome,
		repo:    repo,
		tree:    tree,
	}
}

func (up UpdateManager) Update(name formula.RepoName, version formula.RepoVersion, url string) error {
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

	if repo.IsLocal && repo.Provider != ZipRemoteProvider {
		return ErrLocalRepo
	}

	if repo.Provider != ZipRemoteProvider {
		latestTag := up.repo.LatestTag(*repo)
		repo.LatestVersion = formula.RepoVersion(latestTag)
	} else if repo.Provider == ZipRemoteProvider {
		repo.LatestVersion = version
		repo.Url = url
	}

	repo.UpdateCache()
	repo.Version = version
	repo.TreeVersion = tree.Version

	if err := up.repo.Create(*repo); err != nil {
		return err
	}

	if err := up.repo.Write(repos); err != nil {
		return err
	}

	repoPath := filepath.Join(up.ritHome, reposDirName, name.String())
	treeData, err := up.tree.Generate(repoPath)
	if err != nil {
		return err
	}

	treeFilePath := filepath.Join(repoPath, tree.FileName)
	bytes, err := json.MarshalIndent(treeData, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(treeFilePath, bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}
