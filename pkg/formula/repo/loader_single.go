package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

const commons = "commons"

type SingleLoader struct {
	treePath string
	formula.Adder
}

func NewSingleLoader(treePath string, adder formula.Adder) SingleLoader {
	return SingleLoader{Adder: adder, treePath: treePath}
}

func (m SingleLoader) Load() error {
	r := formula.Repository{
		Priority: 0,
		Name:     commons,
		TreePath: m.treePath,
	}

	if err := m.Add(r); err != nil {
		return err
	}

	return nil
}
