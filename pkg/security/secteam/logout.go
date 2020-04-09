package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type LogoutManager struct {
	provider       security.AuthProvider
	sessionManager session.Manager
	serverURL string
}

func NewLogoutManager(p security.AuthProvider, sm session.Manager, serverURL string) LogoutManager {
	return LogoutManager{provider: p, sessionManager: sm, serverURL: serverURL}
}

func (l LogoutManager) Logout() error {
	session, err := l.sessionManager.Current()
	if err != nil {
		return err
	}

	cr, err := logoutChannelProvider(l.provider, session.Organization, l.serverURL)
	if err != nil {
		return err
	}
	resp := <-cr
	if resp.Error != nil {
		return resp.Error
	}

	return l.sessionManager.Destroy()
}
