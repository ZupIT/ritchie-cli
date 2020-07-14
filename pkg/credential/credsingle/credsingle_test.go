package credsingle

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

var (
	githubCred credential.Detail
	ctxFinder  rcontext.Finder
)

type ctxFinderMock struct{}

func (ctxFinderMock) Find() (holder rcontext.ContextHolder, err error) {
	return rcontext.ContextHolder{}, nil
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

	e := m.Run()
	os.Exit(e)
}
