package credsingle

import (
	"encoding/json"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type Setter struct {
	homePath       string
	ctxFinder      rcontext.Finder
	sessionManager session.Manager
	dir            stream.DirCreater
	file           stream.FileWriter
}

func NewSetter(homePath string, cf rcontext.Finder, sm session.Manager, dir stream.DirCreater, file stream.FileWriter) Setter {
	return Setter{
		homePath:       homePath,
		ctxFinder:      cf,
		sessionManager: sm,
		dir:            dir,
		file:           file,
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
	if err := s.dir.Create(dir); err != nil {
		return err
	}

	credFile := File(s.homePath, ctx.Current, cred.Service)
	if err := s.file.Write(credFile, []byte(cipher)); err != nil {
		return err
	}

	return nil

}
