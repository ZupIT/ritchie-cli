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
	"fmt"
	"io/ioutil"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const (
	treeLocalCmdPattern = "%s/repos/local/tree.json"
	treeRepoCmdPattern  = "%s/repos/%s/tree.json"
	core                = "CORE"
	local               = "LOCAL"
)

type Manager struct {
	ritchieHome string
	repoLister  formula.RepositoryLister
	coreCmds    []api.Command
}

func NewTreeManager(ritchieHome string, rl formula.RepositoryLister, coreCmds []api.Command) Manager {
	return Manager{ritchieHome: ritchieHome, repoLister: rl, coreCmds: coreCmds}
}

func (d Manager) Tree() (map[string]formula.Tree, error) {
	trees := make(map[string]formula.Tree)
	trees[core] = formula.Tree{Commands: d.coreCmds}

	treeLocal, err := d.localTree()
	if err != nil {
		return nil, err
	}
	trees[local] = treeLocal

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
	treeLocal, err := d.localTree()
	if err == nil {
		var cc []api.Command
		for _, v := range treeLocal.Commands {
			key := v.Parent + "_" + v.Usage
			if trees[key].Usage == "" {
				v.Repo = "local"
				trees[key] = v
				cc = append(cc, v)
			}
		}
		treeMain.Commands = append(treeMain.Commands, cc...)
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

func (d Manager) localTree() (formula.Tree, error) {
	treeCmdFile := fmt.Sprintf(treeLocalCmdPattern, d.ritchieHome)
	return loadTree(treeCmdFile)
}

func (d Manager) treeByRepo(repoName formula.RepoName) (formula.Tree, error) {
	treeCmdFile := fmt.Sprintf(treeRepoCmdPattern, d.ritchieHome, repoName)
	return loadTree(treeCmdFile)
}

func loadTree(treeCmdFile string) (formula.Tree, error) {
	tree := formula.Tree{}
	if !fileutil.Exists(treeCmdFile) {
		return tree, nil
	}

	treeFile, err := ioutil.ReadFile(treeCmdFile)
	if err != nil {
		return tree, err
	}

	if err = json.Unmarshal(treeFile, &tree); err != nil {
		return tree, err
	}

	return tree, nil
}
