package security

import (
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	// ErrPassphraseIsRequired error for required passphrase
	ErrPassphraseIsRequired = prompt.Error("passphrase is required")
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
