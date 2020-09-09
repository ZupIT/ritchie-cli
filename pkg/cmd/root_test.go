package cmd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type stableVersionCacheMock struct {
	Stable    string `json:"stableVersion"`
	ExpiresAt int64  `json:"expiresAt"`
}

func buildStableBodyMock(expiresAt int64) []byte {
	cache := stableVersionCacheMock{
		Stable:    "2.0.4",
		ExpiresAt: expiresAt,
	}
	b, _ := json.Marshal(cache)
	return b
}

func Test_rootCmd(t *testing.T) {
	type in struct {
		dir stream.DirCreateChecker
		vm  version.Manager
	}

	notExpiredCache := time.Now().Add(time.Hour).Unix()
	versionManager := version.NewManager(
		"any value",
		sMocks.FileWriteReadExisterCustomMock{
			ExistsMock: func(path string) bool {
				return true
			},
			ReadMock: func(path string) ([]byte, error) {
				return buildStableBodyMock(notExpiredCache), nil
			},
		},
	)

	var tests = []struct {
		name    string
		wantErr bool
		in      in
	}{
		{
			name:    "Run with success",
			wantErr: false,
			in: in{
				dir: DirManagerCustomMock{
					exists: func(dir string) bool {
						return true
					},
					create: func(dir string) error {
						return nil
					},
				},
				vm: versionManager,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := os.TempDir()
			rootCmd := NewRootCmd(tmpDir, tt.in.dir, TutorialFinderMock{}, tt.in.vm)

			if err := rootCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("root error = %v | error wanted: %v", err, tt.wantErr)
			}
		})
	}
}
