package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"net/http"
)

type LoginManager struct {
	homePath 			string
	serverFinder 		server.Finder
	provider            security.AuthProvider
	httpClient          *http.Client
	sessionManager      session.Manager
}

func NewLoginManager(
	homePath string,
	serverFinder server.Finder,
	provider security.AuthProvider,
	hc *http.Client,
	sm session.Manager) LoginManager {
	return LoginManager{
		homePath:       homePath,
		serverFinder:   serverFinder,
		provider:       provider,
		httpClient:     hc,
		sessionManager: sm,
	}
}

func (l LoginManager) Login(p security.Passcode) error {
	org := p.String()
	serverUrl, err := l.serverFinder.Find()
	if err != nil {
		return err
	}

	cr, err := loginChannelProvider(l.provider, org, serverUrl)
	if err != nil {
		return err
	}
	resp := <-cr
	if resp.Error != nil {
		return resp.Error
	}

	sess := session.Session{
		AccessToken:  resp.Token,
		Organization: org,
		Username:     resp.Username,
	}
	err = l.sessionManager.Create(sess)
	if err != nil {
		return err
	}

	return nil
}
