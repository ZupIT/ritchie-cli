package credential

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const errNotFoundTemplate = `
Fail to resolve credential for provider %s.
Try again after use:
	∙ rit set credential
`

type Finder struct {
	homePath  string
	ctxFinder rcontext.Finder
	file      stream.FileReader
}

func NewFinder(homePath string, cf rcontext.Finder, file stream.FileReader) Finder {
	return Finder{
		homePath:  homePath,
		ctxFinder: cf,
		file:      file,
	}
}

func (f Finder) Find(provider string) (Detail, error) {
	ctx, err := f.ctxFinder.Find()

	if err != nil {
		return Detail{}, err
	}
	if ctx.Current == "" {
		ctx.Current = rcontext.DefaultCtx
	}

	cb, err := f.file.Read(File(f.homePath, ctx.Current, provider))
	if err != nil {
		errMsg := fmt.Sprintf(errNotFoundTemplate, provider)
		return Detail{}, errors.New(prompt.Red(errMsg))
	}

	cred := &Detail{}
	if err := json.Unmarshal(cb, cred); err != nil {
		return Detail{}, err
	}
	return *cred, nil

}
