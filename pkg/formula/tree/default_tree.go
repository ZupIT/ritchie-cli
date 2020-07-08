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
	repoLister  formula.RepoLister
	coreCmds    []api.Command
}

func NewTreeManager(ritchieHome string, rl formula.RepoLister, coreCmds []api.Command) Manager {
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
		trees[v.Name] = treeRepo
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
				c.Repo = r.Name
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

func (d Manager) treeByRepo(repoName string) (formula.Tree, error) {
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
