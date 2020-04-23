package sessteam

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	sessionManager session.Manager
	validator      session.Validator
)

func TestMain(m *testing.M) {
	homePath := os.TempDir()
	fileManager := stream.NewFileManager()
	sessionManager = session.NewManager(homePath, fileManager)
	validator = NewValidator(sessionManager)
	os.Exit(m.Run())
}

func TestValidate(t *testing.T) {

	type in struct {
		session session.Session
		exp     int64
	}

	tests := []struct {
		name string
		in   in
		out  error
	}{
		{
			name: "team session",
			in: in{
				session: session.Session{
					Organization: "zup",
					Username:     "dennis.ritchie",
				},
				exp: time.Now().Add(time.Minute * 15).Unix(),
			},
			out: nil,
		},
		{
			name: "no team session",
			in:   in{},
			out:  session.ErrNoSession,
		},
		{
			name: "expired session token",
			in: in{
				session: session.Session{
					Organization: "zup",
					Username:     "dennis.ritchie",
				},
				exp: time.Now().Add(time.Minute*15).Unix() - 1500,
			},
			out: ErrExpiredToken,
		},
		{
			name: "invalid access token",
			in: in{
				session: session.Session{
					AccessToken:  "dasdasdasdas",
					Organization: "zup",
					Username:     "dennis.ritchie",
				},
			},
			out: ErrInvalidToken,
		},
		{
			name: "decode base64 error",
			in: in{
				session: session.Session{
					AccessToken:  "ds.ds##$F/[",
					Organization: "zup",
					Username:     "dennis.ritchie",
				},
			},
			out: errors.New("illegal base64 data at input byte 2"),
		},
		{
			name: "unmarshall error token",
			in: in{
				session: session.Session{
					AccessToken:  "eyJhbGciOiJIUzI1N.eyJzdWIilIiwiaWF0IjoxNTEIyfQ.SflKxw_adQssw5c",
					Organization: "zup",
					Username:     "dennis.ritchie",
				},
			},
			out: errors.New("invalid character '\\u0094' after object key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = sessionManager.Destroy()

			if tt.in.session.Organization != "" {

				if tt.in.session.AccessToken == "" {
					atClaims := jwt.MapClaims{}
					atClaims["exp"] = tt.in.exp
					at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
					token, _ := at.SignedString([]byte("Test"))
					tt.in.session.AccessToken = token
				}

				err := sessionManager.Create(tt.in.session)
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
