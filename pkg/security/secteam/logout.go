package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type LogoutManager struct {
	sessionManager session.Manager
}

func NewLogoutManager(sm session.Manager) LogoutManager {
	return LogoutManager{
		sessionManager: sm,
	}
}

func (l LogoutManager) Logout() error {
	return l.sessionManager.Destroy()
}
