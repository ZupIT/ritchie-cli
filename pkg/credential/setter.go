package credential

import (
	"encoding/json"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type Setter struct {
	homePath  string
	ctxFinder rcontext.Finder
}

func NewSetter(homePath string, cf rcontext.Finder) Setter {
	return Setter{
		homePath:  homePath,
		ctxFinder: cf,
	}
}

func (s Setter) Set(cred Detail) error {
	ctx, err := s.ctxFinder.Find()
	if err != nil {
		return err
	} else if ctx.Current == "" {
		ctx.Current = rcontext.DefaultCtx
	}

	cb, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	dir := Dir(s.homePath, ctx.Current)
	if err := fileutil.CreateDirIfNotExists(dir, 0700); err != nil {
		return err
	}

	credFile := File(s.homePath, ctx.Current, cred.Service)
	if err := fileutil.WriteFilePerm(credFile, cb, 0600); err != nil {
		return err
	}

	return nil

}
