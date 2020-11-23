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

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	reposDirName = "repos"
	treeFileName = "tree.json"
	core         = "CORE"
)

type Manager struct {
	ritchieHome   string
	repoLister    formula.RepositoryLister
	coreCmds      []api.Command
	file          stream.FileReadExister
	repoProviders formula.RepoProviders
	isRootCommand bool
}

func NewTreeManager(
	ritchieHome string,
	rl formula.RepositoryLister,
	coreCmds []api.Command,
	file stream.FileReadExister,
	rp formula.RepoProviders,
	isRootCommand bool,
) Manager {
	return Manager{
		ritchieHome:   ritchieHome,
		repoLister:    rl,
		coreCmds:      coreCmds,
		file:          file,
		repoProviders: rp,
		isRootCommand: isRootCommand,
	}
}

func (d Manager) Tree() (map[string]formula.Tree, error) {
	trees := make(map[string]formula.Tree)
	trees[core] = formula.Tree{Commands: d.coreCmds}

	rr, err := d.repoLister.List()
	if err != nil {
		return nil, err
	}
	for _, v := range rr {
		treeRepo, err := d.treeByRepo(v.Name)
		if err != nil {
			return nil, err
		}
		trees[v.Name.String()] = treeRepo
	}

	return trees, nil
}

func (d Manager) MergedTree(core bool) formula.Tree {
	trees := make(map[string]api.Command)
	treeMain := formula.Tree{Commands: []api.Command{}}
	if core {
		treeMain = formula.Tree{Commands: d.coreCmds}
	}

	for _, v := range treeMain.Commands {
		key := v.Parent + "_" + v.Usage
		trees[key] = v
	}

	rr, _ := d.repoLister.List()
	for _, r := range rr {
		treeRepo, err := d.treeByRepo(r.Name)
		if err != nil {
			continue
		}

		var cc []api.Command
		for _, c := range treeRepo.Commands {
			key := c.Parent + "_" + c.Usage
			if trees[key].Usage == "" {
				c.Repo = r.Name.String()
				trees[key] = c
				cc = append(cc, c)
			}
		}
		treeMain.Commands = append(treeMain.Commands, cc...)
	}

	return treeMain
}

// nolint
func (d Manager) getLatestTag(repo formula.Repo) string {
	formulaGit := d.repoProviders.Resolve(repo.Provider)

	repoInfo := formulaGit.NewRepoInfo(repo.URL, repo.Token)
	tag, err := formulaGit.Repos.LatestTag(repoInfo)
	if err != nil {
		return ""
	}

	return tag.Name
}

func (d Manager) treeByRepo(repoName formula.RepoName) (formula.Tree, error) {
	treeCmdFile := filepath.Join(d.ritchieHome, reposDirName, repoName.String(), treeFileName)
	return d.loadTree(treeCmdFile)
}

func (d Manager) loadTree(treeCmdFile string) (formula.Tree, error) {
	tree := formula.Tree{}

	if !d.file.Exists(treeCmdFile) {
		return tree, nil
	}

	treeFile, err := d.file.Read(treeCmdFile)
	if err != nil {
		return tree, err
	}

	if err = json.Unmarshal(treeFile, &tree); err != nil {
		return tree, err
	}

	return tree, nil
}
