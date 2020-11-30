package credential

import (
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type DeleteManager struct {
	homePath    string
	ctxFinder   rcontext.Finder
	fileRemover stream.FileRemover
}

func NewCredDelete(homePath string, cf rcontext.Finder, fm stream.FileRemover) DeleteManager {
	return DeleteManager{
		homePath:    homePath,
		ctxFinder:   cf,
		fileRemover: fm,
	}
}

func (d DeleteManager) Delete(service string) error {
	ctx, err := d.ctxFinder.Find()
	if err != nil {
		return err
	}

	if ctx.Current == "" {
		ctx.Current = rcontext.DefaultCtx
	}

	if err := d.fileRemover.Remove(File(d.homePath, ctx.Current, service)); err != nil {
		return err
	}
	return nil
}
