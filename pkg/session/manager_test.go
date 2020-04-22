package session

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	sessionManager Manager
)

func TestMain(m *testing.M) {
	homePath := os.TempDir()
	fileReader := stream.NewFileReader()
	fileWriter := stream.NewFileWriter()
	fileExister := stream.NewFileExister()
	fileRemover := stream.NewFileRemover(fileExister)
	fileManager := stream.NewFileManager(fileWriter, fileReader, fileExister, fileRemover)
	sessionManager = NewManager(homePath, fileManager)
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
		in   Session
		out  error
	}{
		{
			name: "team session",
			in: Session{
				AccessToken:  "SflKxwRJSM.eKKF2QT4fwpMeJf36.POk6yJV_adQssw5c",
				Organization: "zup",
				Username:     "dennis.ritchie",
			},
			out: nil,
		},
		{
			name: "single session",
			in: Session{
				Secret: "s3cr3t",
			},
			out: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = sessionManager.Destroy()

			err := sessionManager.Create(tt.in)
			if err != nil {
				t.Errorf("Create(%s) got %v, want %v", tt.name, err, tt.out)
			}

		})
	}
}

func TestFind(t *testing.T) {
	type out struct {
		want Session
		err  error
	}

	tests := []struct {
		name string
		out  out
	}{
		{
			name: "team session",
			out: out{
				want: Session{
					AccessToken:  "SflKxwRJSM.eKKF2QT4fwpMeJf36.POk6yJV_adQssw5c",
					Organization: "zup",
					Username:     "dennis.ritchie",
				},
			},
		},
		{
			name: "single session",
			out: out{
				want: Session{
					Secret: "s3cr3t",
				},
			},
		},
		{
			name: "no session",
			out: out{
				err: ErrNoSession,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = sessionManager.Destroy()

			out := tt.out
			if out.want.Secret != "" || out.want.Organization != "" {
				err := sessionManager.Create(out.want)
				if err != nil {
					t.Errorf("Create(%s) got %v, want %v", tt.name, err, tt.out)
				}
			}

			got, err := sessionManager.Current()
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out)
			} else if out.err == nil && got.Secret == "" {
				t.Errorf("Find(%s) got %v, want != nil", tt.name, got)
			}

		})
	}
}

func TestRemove(t *testing.T) {

	tests := []struct {
		name string
		in   Session
		out  error
	}{
		{
			name: "team session",
			in: Session{
				AccessToken:  "SflKxwRJSM.eKKF2QT4fwpMeJf36.POk6yJV_adQssw5c",
				Organization: "zup",
				Username:     "dennis.ritchie",
			},
			out: nil,
		},
		{
			name: "single session",
			in: Session{
				Secret: "s3cr3t",
			},
			out: nil,
		},
		{
			name: "no session",
			in:   Session{},
			out:  ErrNoSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = sessionManager.Destroy()

			if tt.in.Secret != "" || tt.in.Organization != "" {
				err := sessionManager.Create(tt.in)
				if err != nil {
					t.Errorf("Create(%s) got %v, want %v", tt.name, err, tt.out)
				}
			}

			out := tt.out
			err := sessionManager.Destroy()
			if err != nil && err.Error() != out.Error() {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, err, out)
			}
		})
	}
}
