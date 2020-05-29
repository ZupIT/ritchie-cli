package secsingle

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

func TestLogin(t *testing.T) {
	homePath := os.TempDir()
	sm := session.NewManager(homePath)
	manager := NewPassphraseManager(sm)

	tests := []struct {
		name string
		in   security.Passphrase
		out  error
	}{
		{
			name: "new passphrase",
			in:   security.Passphrase("s3cr3t"),
			out:  nil,
		},
		{
			name: "empty passphrase",
			in:   "",
			out:  security.ErrPassphraseIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out

			got := manager.Save(tt.in)
			if got != out {
				t.Errorf("Save(%s) got %v, want %v", tt.name, got, out)
			}

		})
	}

}
