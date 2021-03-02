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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
)

var ErrInvalidRepo = errors.New("the selected repository has no formulas")

type AddManager struct {
	ritHome string
	repo    formula.RepositoryCreateWriteListDetailDeleter
	tree    formula.TreeGenerator
}

func NewAdder(
	ritHome string,
	repo formula.RepositoryCreateWriteListDetailDeleter,
	tree formula.TreeGenerator,
) AddManager {
	return AddManager{
		ritHome: ritHome,
		repo:    repo,
		tree:    tree,
	}
}

func (ad AddManager) Add(repo formula.Repo) error {
	if !repo.IsLocal {
		latestTag := ad.repo.LatestTag(repo)
		repo.LatestVersion = formula.RepoVersion(latestTag)
		repo.UpdateCache()

		if err := ad.repo.Create(repo); err != nil {
			return err
		}
	}

	repos, err := ad.repo.List()
	if err != nil {
		return err
	}

	repo.TreeVersion = tree.Version
	repos = setPriority(repo, repos)

	if err := ad.repo.Write(repos); err != nil {
		return err
	}

	if err := ad.treeGenerate(repo); err != nil {
		return err
	}

	return nil
}

func (ad AddManager) treeGenerate(repo formula.Repo) error {
	newRepoPath := filepath.Join(ad.ritHome, reposDirName, repo.Name.String())
	tree, err := ad.tree.Generate(newRepoPath)
	if err != nil {
		return err
	}

	if err := ad.isValidRepo(repo, tree); err != nil {
		return err
	}

	treeFilePath := filepath.Join(newRepoPath, "tree.json")
	bytes, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(treeFilePath, bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (ad AddManager) isValidRepo(repo formula.Repo, tree formula.Tree) error {
	if len(tree.Commands) == 0 && !repo.IsLocal {
		if err := ad.repo.Delete(repo.Name); err != nil {
			return err
		}
		return ErrInvalidRepo
	}
	return nil
}

func setPriority(repo formula.Repo, repos formula.Repos) formula.Repos {
	exist := func() bool {
		for i := range repos {
			r := repos[i]
			if repo.Name == r.Name {
				repos[i] = repo
				return true
			}
		}
		return false
	}

	if !exist() {
		repos = append(repos, repo)
	}

	repos = movePosition(repos, repo.Name, repo.Priority)

	return repos
}
