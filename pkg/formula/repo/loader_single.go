package repo

const commons = "commons"

type SingleLoader struct {
	treePath string
	Adder
}

func NewSingleLoader(treePath string, adder Adder) SingleLoader {
	return SingleLoader{Adder: adder, treePath: treePath}
}

func (m SingleLoader) Load() error {
	r := Repository{
		Priority: 0,
		Name:     commons,
		TreePath: m.treePath,
	}

	if err := m.Add(r); err != nil {
		return err
	}

	return nil
}
