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
	"path/filepath"
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const core = "CORE"

type Manager struct {
	ritchieHome   string
	repo          formula.RepositoryLister
	coreCmds      api.Commands
	file          stream.FileReadExister
	repoProviders formula.RepoProviders
}

func NewTreeManager(
	ritchieHome string,
	rl formula.RepositoryLister,
	coreCmds api.Commands,
	file stream.FileReadExister,
	rp formula.RepoProviders,
) Manager {
	return Manager{
		ritchieHome:   ritchieHome,
		repo:          rl,
		coreCmds:      coreCmds,
		file:          file,
		repoProviders: rp,
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
		treeRepo, err := d.treeByRepo(v.Name)
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

	for i := rr.Len() - 1; i >= 0; i-- {
		tree, err := d.treeByRepo(rr[i].Name)
		if err != nil {
			continue
		}

		for k, v := range tree.Commands {
			v.Repo = rr[i].Name.String()
			mergedCommands[k] = v
		}
	}

	var ids []api.CommandID
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
		Version:    treeVersion,
		Commands:   mergedCommands,
		CommandsID: ids,
	}
}

func (d Manager) treeByRepo(repoName formula.RepoName) (formula.Tree, error) {
	tree := formula.Tree{}
	treeFilePath := filepath.Join(d.ritchieHome, "repos", repoName.String(), "tree.json")
	if !d.file.Exists(treeFilePath) {
		return tree, nil
	}

	treeFile, err := d.file.Read(treeFilePath)
	if err != nil {
		return tree, err
	}

	if err = json.Unmarshal(treeFile, &tree); err != nil {
		return tree, err
	}

	return tree, nil
}
