package formula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const (
	treeLocalCmdPattern = "%s/repo/local/tree.json"
	treeRepoCmdPattern  = "%s/repo/cache/%s-tree.json"
	core                = "CORE"
	local               = "LOCAL"
)

type TreeManager struct {
	ritchieHome string
	repoLister  Lister
	coreCmds    []api.Command
}

func NewTreeManager(ritchieHome string, rl Lister, coreCmds []api.Command) TreeManager {
	return TreeManager{ritchieHome: ritchieHome, repoLister: rl, coreCmds: coreCmds}
}

func (d TreeManager) Tree() (map[string]Tree, error) {
	trees := make(map[string]Tree)
	trees[core] = Tree{d.coreCmds}

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

func (d TreeManager) MergedTree(core bool) Tree {
	trees := make(map[string]api.Command)
	treeMain := Tree{[]api.Command{}}
	if core {
		treeMain = Tree{d.coreCmds}
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
				trees[key] = v
				cc = append(cc, v)
			}
		}
		treeMain.Commands = append(treeMain.Commands, cc...)
	}

	rr, _ := d.repoLister.List()
	for _, v := range rr {
		treeRepo, err := d.treeByRepo(v.Name)
		if err != nil {
			continue
		}
		var cc []api.Command
		for _, v := range treeRepo.Commands {
			key := v.Parent + "_" + v.Usage
			if trees[key].Usage == "" {
				trees[key] = v
				cc = append(cc, v)
			}
		}
		treeMain.Commands = append(treeMain.Commands, cc...)
	}

	return treeMain
}

func (d TreeManager) localTree() (Tree, error) {
	treeCmdFile := fmt.Sprintf(treeLocalCmdPattern, d.ritchieHome)
	return loadTree(treeCmdFile)
}

func (d TreeManager) treeByRepo(repo string) (Tree, error) {
	treeCmdFile := fmt.Sprintf(treeRepoCmdPattern, d.ritchieHome, repo)
	return loadTree(treeCmdFile)
}

func loadTree(treeCmdFile string) (Tree, error) {
	tree := Tree{}
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
