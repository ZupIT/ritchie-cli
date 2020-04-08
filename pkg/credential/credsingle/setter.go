package credsingle

import (
	"encoding/json"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type Setter struct {
	homePath       string
	ctxFinder      rcontext.Finder
	sessionManager session.Manager
}

func NewSetter(homePath string, cf rcontext.Finder, sm session.Manager) Setter {
	return Setter{
		homePath:       homePath,
		ctxFinder:      cf,
		sessionManager: sm,
	}
}

func (s Setter) Set(cred credential.Detail) error {
	ctx, err := s.ctxFinder.Find()
	if err != nil {
		return err
	} else if ctx.Current == "" {
		ctx.Current = rcontext.DefaultCtx
	}

	session, err := s.sessionManager.Current()
	if err != nil {
		return err
	}

	cb, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	hash, err := cryptoutil.SumHashMachine(session.Secret)
	if err != nil {
		return err
	}

	cipher := cryptoutil.Encrypt(hash, string(cb))

	dir := Dir(s.homePath, ctx.Current)
	if err := fileutil.CreateDirIfNotExists(dir, 0700); err != nil {
		return err
	}

	credFile := File(s.homePath, ctx.Current, cred.Service)
	if err := fileutil.WriteFilePerm(credFile, []byte(cipher), 0600); err != nil {
		return err
	}

	return nil

}
