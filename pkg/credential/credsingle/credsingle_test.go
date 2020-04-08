package credsingle

import (
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"os"
	"testing"
)

var (
	githubCred  credential.Detail
	ctxFinder   rcontext.Finder
	sessManager session.Manager
)

type ctxFinderMock struct{}

func (ctxFinderMock) Find() (holder rcontext.ContextHolder, err error) {
	return rcontext.ContextHolder{}, nil
}

type sessionMock struct{}

func (sessionMock) Create(s session.Session) error {
	return nil
}

func (sessionMock) Current() (s session.Session, err error) {
	return session.Session{
		Secret:   "s3cr3t",
		Username: "dennis.ritchie",
	}, nil
}

func (sessionMock) Destroy() error {
	return nil
}

func TestMain(m *testing.M) {
	githubCred = credential.Detail{
		Username: "dennis.ritchie",
		Credential: credential.Credential{
			"username": "dennis.ritchie",
			"password": "unix@clang",
		},
		Service: "github",
	}

	ctxFinder = ctxFinderMock{}
	sessManager = sessionMock{}

	e := m.Run()
	os.Exit(e)
}
