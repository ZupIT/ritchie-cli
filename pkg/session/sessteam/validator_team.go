package sessteam

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"strings"
	"time"
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
		return fmt.Errorf("oidc: malformed jwt, expected 3 parts got %d", size)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("oidc: malformed jwt payload: %v", err)
	}

	type Token struct {
		Exp int64 `json:"exp"`
	}
	var token Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return fmt.Errorf("error unmarshal token: %v", err)
	}

	tokenTime := time.Unix(token.Exp, 0)
	if time.Since(tokenTime).Seconds() > 0 {
		return fmt.Errorf("token expired, time token: %v", token.Exp)
	}

	return nil
}
