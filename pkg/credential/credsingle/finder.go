package credsingle

import (
	"encoding/json"

	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type Finder struct {
	homePath       string
	ctxFinder      rcontext.Finder
	sessionManager session.Manager
	file           stream.FileReader
}

func NewFinder(homePath string, cf rcontext.Finder, sm session.Manager, file stream.FileReader) Finder {
	return Finder{
		homePath:       homePath,
		ctxFinder:      cf,
		sessionManager: sm,
		file:           file,
	}
}

func (f Finder) Find(provider string) (credential.Detail, error) {
	ctx, err := f.ctxFinder.Find()
	if err != nil {
		return credential.Detail{}, err
	} else if ctx.Current == "" {
		ctx.Current = rcontext.DefaultCtx
	}

	cb, err := f.file.Read(File(f.homePath, ctx.Current, provider))
	if err != nil {
		return credential.Detail{}, err
	}

	session, err := f.sessionManager.Current()
	if err != nil {
		return credential.Detail{}, err
	}

	hash, err := cryptoutil.SumHashMachine(session.Secret)
	if err != nil {
		return credential.Detail{}, err
	}

	plain := cryptoutil.Decrypt(hash, string(cb))
	cred := &credential.Detail{}
	if err := json.Unmarshal([]byte(plain), cred); err != nil {
		return credential.Detail{}, err
	}

	return *cred, nil
}
