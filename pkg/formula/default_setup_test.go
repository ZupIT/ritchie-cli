package formula

import (
	"net/http"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

func TestDefaultSetup_Setup(t *testing.T) {
	def := Definition{
		Path:    "mock/test",
		Bin:     "test-${so}",
		LBin:    "test-${so}",
		MBin:    "test-${so}",
		WBin:    "test-${so}.exe",
		Bundle:  "${so}.zip",
		Config:  "config.json",
		RepoURL: "http://localhost:8882/formulas",
	}

	home := os.TempDir()

	type in struct {
		sess   session.Manager
		config string
		bundle string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				sess: sessManagerMock{sess: session.Session{
					AccessToken:  "1234",
					Organization: "my-org",
					Username:     "fakename",
					Secret:       "1234",
				}},
			},
			want: nil,
		},
		{
			name: "config not found",
			in: in{
				sess: sessManagerMock{sess: session.Session{
					AccessToken:  "1234",
					Organization: "my-org",
					Username:     "fakename",
					Secret:       "1234",
				}},
				config: "config-not-found",
			},
			want: ErrConfigFileNotFound,
		},
		{
			name: "bundle not found",
			in: in{
				sess: sessManagerMock{sess: session.Session{
					AccessToken:  "1234",
					Organization: "my-org",
					Username:     "fakename",
					Secret:       "1234",
				}},
				bundle: "bundle-not-found",
			},
			want: ErrFormulaBinNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			defTest := def
			if in.config != "" {
				defTest.Config = in.config
			}
			if in.bundle != "" {
				defTest.Bundle = in.bundle
			}

			_ = fileutil.RemoveDir(home + "/formulas")
			setup := NewDefaultTeamSetup(home, http.DefaultClient, in.sess)
			_, got := setup.Setup(defTest)

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Setup(%s) got %v, want %v", tt.name, got, tt.want)
			}

		})
	}
}

type sessManagerMock struct {
	sess  session.Session
	error error
}

func (s sessManagerMock) Create(session.Session) error {
	return s.error
}
func (s sessManagerMock) Current() (session.Session, error) {
	return s.sess, s.error
}

func (s sessManagerMock) Destroy() error {
	return s.error
}
