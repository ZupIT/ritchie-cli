package secsingle

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type LoginManager struct {
	session session.Manager
}

func NewLoginManager(s session.Manager) LoginManager {
	return LoginManager{session: s}
}

func (s LoginManager) Login(p security.Passcode) error {
	if p == "" {
		return security.ErrPasscodeIsRequired
	}
	sess := session.Session{Secret: p.String()}
	return s.session.Create(sess)
}
