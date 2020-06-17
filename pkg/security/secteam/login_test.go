package secteam

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"errors"
	"net/http"
	"testing"
)

func TestLoginManager_Login(t *testing.T) {
	type fields struct {
		serverFinder   server.Finder
		httpClient     *http.Client
		sessionManager session.Manager
	}
	type args struct {
		user security.User
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		outErr bool
	}{
		{
			name:    "login success",
			fields:  fields{
				serverFinder:   finderMock{
					sc:  server.Config{
						Organization: "test",
						URL:          "http://localhost:8882",
					},
					err: nil,
				},
				httpClient:     http.DefaultClient,
				sessionManager: sessionManagerMock{},
			},
			in:     args{
				user: security.User{
					Username: "test",
					Password: "test",
				},
			},
			outErr: false,
		},
		{
			name:    "server finder error",
			fields:  fields{
				serverFinder:   finderMock{
					sc:  server.Config{},
					err: errors.New("error"),
				},
				httpClient:     http.DefaultClient,
				sessionManager: sessionManagerMock{},
			},
			in:     args{
				user: security.User{
					Username: "test",
					Password: "test",
				},
			},
			outErr: true,
		},
		{
			name:    "login failed not found",
			fields:  fields{
				serverFinder:   finderMock{
					sc:  server.Config{
						Organization: "test",
						URL:          "http://localhost:8882/notfound",
					},
					err: nil,
				},
				httpClient:     http.DefaultClient,
				sessionManager: sessionManagerMock{},
			},
			in:     args{
				user: security.User{
					Username: "test",
					Password: "test",
				},
			},
			outErr: true,
		},
		{
			name:    "login failed unauthorized",
			fields:  fields{
				serverFinder:   finderMock{
					sc:  server.Config{
						Organization: "test",
						URL:          "http://localhost:8882/unauthorized",
					},
					err: nil,
				},
				httpClient:     http.DefaultClient,
				sessionManager: sessionManagerMock{},
			},
			in:     args{
				user: security.User{
					Username: "test",
					Password: "test",
				},
			},
			outErr: true,
		},
		{
			name:    "session created error",
			fields:  fields{
				serverFinder:   finderMock{
					sc:  server.Config{
						Organization: "test",
						URL:          "http://localhost:8882",
					},
					err: nil,
				},
				httpClient:     http.DefaultClient,
				sessionManager: sessionManagerMock{err: errors.New("error")},
			},
			in:     args{
				user: security.User{
					Username: "test",
					Password: "test",
				},
			},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoginManager(tt.fields.serverFinder, tt.fields.httpClient, tt.fields.sessionManager)
			if err := l.Login(tt.in.user); (err != nil) != tt.outErr {
				t.Errorf("Login() error = %v, outErr %v", err, tt.outErr)
			}
		})
	}
}

type finderMock struct {
	sc server.Config
	err error
}

func (f finderMock) Find() (server.Config, error) {
	return f.sc, f.err
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