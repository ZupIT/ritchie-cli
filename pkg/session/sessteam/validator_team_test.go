package sessteam

import (
	"testing"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/session"
)

func TestValidator_Validate(t1 *testing.T) {
	type fields struct {
		manager session.Manager
	}
	tests := []struct {
		name   string
		fields fields
		out error
	}{
		{
			name:   "success",
			fields: fields{
				manager: sessionManagerMock{
					se: session.Session{
						AccessToken:  "token",
						Organization: "test",
						Username:     "test",
						Secret:       "test",
						TTL:          time.Now().Unix() + 3000,
					},
					err: nil,
				},
			},
			out: nil,
		},
		{
			name:   "expired",
			fields: fields{
				manager: sessionManagerMock{
					se: session.Session{
						AccessToken:  "token",
						Organization: "test",
						Username:     "test",
						Secret:       "test",
						TTL:          time.Now().Unix() -10,
					},
					err: nil,
				},
			},
			out: ErrExpiredToken,
		},

	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			v := NewValidator(tt.fields.manager)
			got := v.Validate()
			if got != tt.out {
				t.Errorf("Save(%s) got %v, want %v", tt.name, got, tt.out)
			}
		})
	}
}

type sessionManagerMock struct {
	se session.Session
	err error
}

func (s sessionManagerMock) Create(se session.Session) error {
	return s.err
}
func (s sessionManagerMock) Current() (session.Session, error) {
	return s.se, nil
}
func (s sessionManagerMock) Destroy() error {
	return s.err
}