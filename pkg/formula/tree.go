package formula

import "github.com/ZupIT/ritchie-cli/pkg/api"

type Tree struct {
	Commands api.Commands `json:"commands"`
}

type TreeManager interface {
	Tree() (map[string]Tree, error)
	MergedTree(core bool) Tree
}
