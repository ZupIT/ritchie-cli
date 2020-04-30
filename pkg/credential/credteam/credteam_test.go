package credteam

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"net/http"
	"net/http/httptest"
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
			"token":    "55c8bb45-806d-4e2d-a5c0-c96a9076c859",
		},
		Service: "github",
	}

	ctxFinder = ctxFinderMock{}
	sessManager = sessionMock{}

	e := m.Run()
	os.Exit(e)
}

func mockServer(status int, body []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
		_, err := rw.Write(body)
		if err != nil {
			fmt.Sprintln("Error in Write")
			return
		}
	}))
}
