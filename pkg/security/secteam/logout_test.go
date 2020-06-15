package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"errors"
	"testing"
)

func TestLogoutManager_Logout(t *testing.T) {
	type fields struct {
		sessionManager session.Manager
	}
	tests := []struct {
		name   string
		fields fields
		outErr bool
	}{
		{
			name:   "logout success",
			fields: fields{
				sessionManager: sessionManagerMock{},
			},
			outErr: false,
		},
		{
			name:   "logout failed",
			fields: fields{
				sessionManager: sessionManagerMock{err: errors.New("error")},
			},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLogoutManager(tt.fields.sessionManager)
			if err := l.Logout(); (err != nil) != tt.outErr {
				t.Errorf("Logout() error = %v, outErr %v", err, tt.outErr)
			}
		})
	}
}