package credential

import (
	renv "github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type DeleteManager struct {
	homePath string
	env      renv.Finder
	file     stream.FileRemover
}

func NewCredDelete(homePath string, env renv.Finder, fm stream.FileRemover) DeleteManager {
	return DeleteManager{
		homePath: homePath,
		env:      env,
		file:     fm,
	}
}

func (d DeleteManager) Delete(service string) error {
	env, err := d.env.Find()
	if err != nil {
		return err
	}

	if env.Current == "" {
		env.Current = renv.Default
	}

	if err := d.file.Remove(File(d.homePath, env.Current, service)); err != nil {
		return err
	}
	return nil
}
