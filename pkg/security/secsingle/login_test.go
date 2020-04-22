package secsingle

import (
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/session"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	homePath := os.TempDir()
	fileWriter := stream.NewFileWriter()
	fileReader := stream.NewFileReader()
	fileExister := stream.NewFileExister()
	fileRemover := stream.NewFileRemover(fileExister)
	fileManager := stream.NewFileManager(fileWriter, fileReader, fileExister, fileRemover)
	sm := session.NewManager(homePath, fileManager)
	manager := NewLoginManager(sm)

	tests := []struct {
		name string
		in   security.Passcode
		out  error
	}{
		{
			name: "single",
			in:   security.Passcode("s3cr3t"),
			out:  nil,
		},
		{
			name: "passcode is required",
			in:   "",
			out:  security.ErrPasscodeIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out

			got := manager.Login(tt.in)
			if got != out {
				t.Errorf("Login(%s) got %v, want %v", tt.name, got, out)
			}

		})
	}

}
