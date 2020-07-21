package credential

import (
	"encoding/json"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type Finder struct {
	homePath  string
	ctxFinder rcontext.Finder
}

func NewFinder(homePath string, cf rcontext.Finder) Finder {
	return Finder{
		homePath:  homePath,
		ctxFinder: cf,
	}
}

func (f Finder) Find(provider string) (Detail, error) {
	ctx, err := f.ctxFinder.Find()
	if err != nil {
		return Detail{}, err
	} else if ctx.Current == "" {
		ctx.Current = rcontext.DefaultCtx
	}

	cb, err := fileutil.ReadFile(File(f.homePath, ctx.Current, provider))
	if err != nil {
		return Detail{}, err
	}

	cred := &Detail{}
	if err := json.Unmarshal(cb, cred); err != nil {
		return Detail{}, err
	}

	return *cred, nil
}
