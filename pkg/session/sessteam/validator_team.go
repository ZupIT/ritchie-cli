package sessteam

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("the access token is invalid. please, you need to start a session")
	ErrExpiredToken = errors.New("the access token has expired. please, you need to start a session")
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

	parts := strings.Split(sess.AccessToken, ".")
	size := len(parts)
	if size < 2 {
		return ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	type Token struct {
		Exp int64 `json:"exp"`
	}
	var token Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return err
	}

	tokenTime := time.Unix(token.Exp, 0)
	if time.Since(tokenTime).Seconds() > 0 {
		return ErrExpiredToken
	}

	return nil
}
