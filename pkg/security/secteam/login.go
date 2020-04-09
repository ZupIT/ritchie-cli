package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"net/http"
)

type LoginManager struct {
	homePath, serverURL string
	provider            security.AuthProvider
	httpClient          *http.Client
	sessionManager      session.Manager
}

func NewLoginManager(
	homePath,
	serverURL string,
	provider security.AuthProvider,
	hc *http.Client,
	sm session.Manager) LoginManager {
	return LoginManager{
		homePath:       homePath,
		serverURL:      serverURL,
		provider:       provider,
		httpClient:     hc,
		sessionManager: sm,
	}
}

func (l LoginManager) Login(p security.Passcode) error {
	org := p.String()
	cr, err := loginChannelProvider(l.provider, org, l.serverURL)
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
