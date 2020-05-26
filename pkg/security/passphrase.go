package security

import "errors"

var (
	// ErrPassphraseIsRequired error for required passphrase
	ErrPassphraseIsRequired = errors.New("passphrase is required")
)

// Passphrase represents a security code defined by the user.
type Passphrase string

func (p Passphrase) String() string {
	return string(p)
}

// PassphraseManager manages passphrase lifecyle
type PassphraseManager interface {
	Save(Passphrase) error
}
