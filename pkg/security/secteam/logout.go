package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type LogoutManager struct {
	provider       security.AuthProvider
	sessionManager session.Manager
	serverFinder   server.Finder
}

func NewLogoutManager(p security.AuthProvider, sm session.Manager, serverFinder server.Finder) LogoutManager {
	return LogoutManager{
		provider:       p,
		sessionManager: sm,
		serverFinder:   serverFinder,
	}
}

func (l LogoutManager) Logout() error {
	session, err := l.sessionManager.Current()
	if err != nil {
		return err
	}

	serverURL, err := l.serverFinder.Find()
	if err != nil {
		return err
	}

	cr, err := logoutChannelProvider(l.provider, session.Organization, serverURL)
	if err != nil {
		return err
	}
	resp := <-cr
	if resp.Error != nil {
		return resp.Error
	}

	return l.sessionManager.Destroy()
}
