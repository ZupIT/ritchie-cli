package sessteam

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const startSession = "please, you need to start a session"

var (
	ErrInvalidToken    = fmt.Errorf("the access token is invalid. %s", startSession)
	ErrExpiredToken    = fmt.Errorf("the access token has expired. %s", startSession)
	ErrDecodeToken     = fmt.Errorf("unable to decode access token. %s", startSession)
	ErrConvertToStruct = fmt.Errorf("couldn't convert access token into the struct. %s", startSession)
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
		return ErrDecodeToken
	}

	type Token struct {
		Exp int64 `json:"exp"`
	}
	var token Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return ErrConvertToStruct
	}

	tokenTime := time.Unix(token.Exp, 0)
	if time.Since(tokenTime).Seconds() > 0 {
		return ErrExpiredToken
	}

	return nil
}
