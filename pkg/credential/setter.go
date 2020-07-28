package credential

import (
	"encoding/json"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type SetManager struct {
	homePath string
	ctx      rcontext.Finder
}

func NewSetter(homePath string, cf rcontext.Finder) SetManager {
	return SetManager{
		homePath: homePath,
		ctx:      cf,
	}
}

func (s SetManager) Set(cred Detail) error {
	ctx, err := s.ctx.Find()
	if err != nil {
		return err
	}
	if ctx.Current == "" {
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
