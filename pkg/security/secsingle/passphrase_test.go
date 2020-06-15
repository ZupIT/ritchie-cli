package secsingle

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

func TestPassphraseManager_Save(t *testing.T) {
	type fields struct {
		session session.Manager
	}
	type args struct {
		p security.Passphrase
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    error
	}{
		{
			name:   "new passphrase",
			fields: fields{
				session: sessionManagerMock{},
			},
			in:     args{
				p: security.Passphrase("s3cr3t"),
			},
			out: nil,
		},
		{
			name:   "error create session",
			fields: fields{
				session: sessionManagerMock{err: security.ErrPassphraseIsRequired},
			},
			in:     args{
				p: security.Passphrase(""),
			},
			out: security.ErrPassphraseIsRequired,
		},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := NewPassphraseManager(tt.fields.session)
			got := pm.Save(tt.in.p)
			if got != tt.out {
				t.Errorf("Save(%s) got %v, want %v", tt.name, got, tt.out)
			}
		})
	}
}

type sessionManagerMock struct {
	err error
}

func (s sessionManagerMock) Create(se session.Session) error {
	return s.err
}
func (s sessionManagerMock) Current() (session.Session, error) {
	return session.Session{}, nil
}
func (s sessionManagerMock) Destroy() error {
	return s.err
}