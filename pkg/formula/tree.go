package formula

import "github.com/ZupIT/ritchie-cli/pkg/api"

type Tree struct {
	Commands []api.Command `json:"commands"`
}

type Manager interface {
	Tree() (map[string]Tree, error)
	MergedTree(core bool) Tree
}
