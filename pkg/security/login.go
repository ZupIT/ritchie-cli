package security

import "errors"

var (
	ErrPasscodeIsRequired = errors.New("passcode is required")
)

// Passcode represents a security code of the user.
// Single Edition: a passphrase defined by the user.
// Team Edition: the organization slug.
type Passcode string

func (p Passcode) String() string {
	return string(p)
}

// LoginManager perform user login by passcode
type LoginManager interface {
	Login(p Passcode) error
}
