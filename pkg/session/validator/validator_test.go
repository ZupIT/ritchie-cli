package validator

import (
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"os"
	"testing"
)

var (
	sessionManager session.Manager
	validator      session.Validator
)

func TestMain(m *testing.M) {
	homePath := os.TempDir()
	sessionManager = session.NewManager(homePath)
	validator = NewValidator(sessionManager)
	os.Exit(m.Run())
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		in   session.Session
		out  error
	}{
		{
			name: "single session",
			in: session.Session{
				Secret: "s3cr3t",
			},
			out: nil,
		},
		{
			name: "no single session",
			in:   session.Session{},
			out:  session.ErrNoSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = sessionManager.Destroy()

			if tt.in.Secret != "" {
				err := sessionManager.Create(tt.in)
				if err != nil {
					t.Errorf("Create(%s) got %v, want %v", tt.name, err, tt.out)
				}
			}

			out := tt.out
			got := validator.Validate()
			if got != nil && got.Error() != out.Error() {
				t.Errorf("Validate(%s) got %v, want %v", tt.name, got, out)
			}

		})
	}
}
