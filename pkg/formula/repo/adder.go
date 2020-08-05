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
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	reposDirName  = "repos"
	reposFileName = "repositories.json"
)

type AddManager struct {
	ritHome string
	repo    formula.RepositoryCreator
	tree    formula.TreeGenerator
	dir     stream.DirCreateListCopyRemover
	file    stream.FileWriteCreatorReadExistRemover
}

func NewAdder(
	ritHome string,
	repo formula.RepositoryCreator,
	tree formula.TreeGenerator,
	dir stream.DirCreateListCopyRemover,
	file stream.FileWriteCreatorReadExistRemover,
) AddManager {
	return AddManager{
		ritHome: ritHome,
		repo:    repo,
		tree:    tree,
		dir:     dir,
		file:    file,
	}
}

func (ad AddManager) Add(repo formula.Repo) error {
	if err := ad.repo.Create(repo); err != nil {
		return err
	}

	repos := formula.Repos{}
	repoPath := filepath.Join(ad.ritHome, reposDirName, reposFileName)
	if ad.file.Exists(repoPath) {
		read, err := ad.file.Read(repoPath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(read, &repos); err != nil {
			return err
		}
	}

	repos = setPriority(repo, repos)

	if err := ad.saveRepo(repoPath, repos); err != nil {
		return err
	}

	newRepoPath := filepath.Join(ad.ritHome, reposDirName, repo.Name.String())

	tree, err := ad.tree.Generate(newRepoPath)
	if err != nil {
		return err
	}

	treeFilePath := filepath.Join(newRepoPath, "tree.json")
	bytes, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := ad.file.Write(treeFilePath, bytes); err != nil {
		return err
	}

	return nil
}

func (ad AddManager) saveRepo(repoPath string, repos formula.Repos) error {
	bytes, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(repoPath)
	if err := ad.dir.Create(dirPath); err != nil {
		return err
	}

	if err := ad.file.Write(repoPath, bytes); err != nil {
		return err
	}

	return nil
}

func setPriority(repo formula.Repo, repos formula.Repos) formula.Repos {
	exist := func() bool {
		for i := range repos {
			r := repos[i]
			if repo.Name == r.Name {
				repos[i].Priority = repo.Priority
				return true
			}
		}
		return false
	}

	if !exist() {
		repos = append(repos, repo)
	}

	for i := range repos {
		r := repos[i]
		if repo.Name != r.Name && r.Priority >= repo.Priority {
			repos[i].Priority++
		}
	}

	sort.Sort(repos)

	return repos
}
