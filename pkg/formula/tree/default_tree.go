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

package tree

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const core = "CORE"

type Manager struct {
	ritchieHome string
	repo        formula.RepositoryListDetailWriter
	coreCmds    api.Commands
}

func NewTreeManager(
	ritchieHome string,
	repo formula.RepositoryListDetailWriter,
	coreCmds api.Commands,
) Manager {
	return Manager{
		ritchieHome: ritchieHome,
		repo:        repo,
		coreCmds:    coreCmds,
	}
}

func (d Manager) Tree() (map[formula.RepoName]formula.Tree, error) {
	trees := make(map[formula.RepoName]formula.Tree)
	trees[core] = formula.Tree{Commands: d.coreCmds}

	rr, err := d.repo.List()
	if err != nil {
		return nil, err
	}

	for _, v := range rr {
		treeRepo, err := d.TreeByRepo(v.Name)
		if err != nil {
			return nil, err
		}
		trees[v.Name] = treeRepo
	}

	return trees, nil
}

func (d Manager) MergedTree(core bool) formula.Tree {
	mergedCommands := make(api.Commands)
	rr, _ := d.repo.List()

	var hasUpdate bool
	wg := &sync.WaitGroup{}
	for i := range rr {
		wg.Add(1)
		go d.updateCache(wg, &hasUpdate, &rr[i])
	}
	wg.Wait()

	if hasUpdate {
		_ = d.repo.Write(rr)
	}

	for i := rr.Len() - 1; i >= 0; i-- {
		tree, err := d.TreeByRepo(rr[i].Name)
		if err != nil {
			continue
		}

		for k, v := range tree.Commands {
			v.Repo = rr[i].Name.String()
			if rr[i].Version != rr[i].LatestVersion && v.Parent == "root" {
				v.RepoNewVersion = rr[i].LatestVersion.String()
			}
			mergedCommands[k] = v
		}
	}

	ids := make([]api.CommandID, 0)
	for id := range mergedCommands {
		ids = append(ids, id)
	}

	sort.Sort(api.ByLen(ids))

	if core {
		for k, v := range d.coreCmds {
			mergedCommands[k] = v
		}
	}

	return formula.Tree{
		Version:    Version,
		Commands:   mergedCommands,
		CommandsID: ids,
	}
}

func (d Manager) TreeByRepo(repoName formula.RepoName) (formula.Tree, error) {
	treeFilePath := filepath.Join(d.ritchieHome, "repos", repoName.String(), FileName)
	treeFile, err := ioutil.ReadFile(treeFilePath)
	if _, err := os.Stat(treeFilePath); os.IsNotExist(err) {
		return formula.Tree{}, nil
	}

	if err != nil {
		return formula.Tree{}, err
	}

	var tree formula.Tree
	if err = json.Unmarshal(treeFile, &tree); err != nil {
		return formula.Tree{}, err
	}

	return tree, nil
}

func (d Manager) updateCache(wg *sync.WaitGroup, hasUpdate *bool, repo *formula.Repo) {
	defer wg.Done()
	if repo.IsLocal || !repo.CacheExpired() {
		return
	}

	if tag := d.repo.LatestTag(*repo); tag != "" {
		repo.LatestVersion = formula.RepoVersion(tag)
		repo.UpdateCache()
		*hasUpdate = true
	}
}
