package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

const commons = "commons"

type SingleLoader struct {
	treePath string
	formula.RepoAdder
}

func NewSingleLoader(treePath string, adder formula.RepoAdder) SingleLoader {
	return SingleLoader{RepoAdder: adder, treePath: treePath}
}

func (m SingleLoader) Load() error {
	r := formula.Repository{
		Priority: 0,
		Name:     commons,
		ZipUrl: m.treePath,
	}

	if err := m.Add(r); err != nil {
		return err
	}

	return nil
}
