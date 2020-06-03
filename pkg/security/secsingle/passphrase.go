package secsingle

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

type PassphraseManager struct {
	session session.Manager
}

func NewPassphraseManager(s session.Manager) PassphraseManager {
	return PassphraseManager{session: s}
}

func (pm PassphraseManager) Save(p security.Passphrase) error {
	if p == "" {
		return security.ErrPassphraseIsRequired
	}
	sess := session.Session{Secret: p.String()}
	return pm.session.Create(sess)
}
