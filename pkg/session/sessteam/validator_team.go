package sessteam

import (
	"fmt"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const startSession = "please, you need to start a session"

var (
	ErrExpiredToken    = fmt.Errorf("the access token has expired. %s", startSession)
)

type Validator struct {
	manager session.Manager
}

func NewValidator(m session.Manager) Validator {
	return Validator{m}
}

func (t Validator) Validate() error {
	sess, err := t.manager.Current()
	if err != nil {
		return err
	}

	tokenTime := time.Unix(sess.TTL, 0)

	if time.Since(tokenTime).Seconds() > 0 {
		return ErrExpiredToken
	}
	return nil
}
