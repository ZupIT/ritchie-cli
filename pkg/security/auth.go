package security

import "errors"

var (
	ErrUnknownProvider = errors.New("unknown provider")

	OAuthProvider = AuthProvider("oauth")
)

// Passcode represents a provider authenticator like oauth, etc.
type AuthProvider string

// ChanResponse represents the channel between ritchie CLI and server that provides login and logout operations.
// Only used in Team Edition.
type ChanResponse struct {
	Token    string
	Username string
	Error    error
}
