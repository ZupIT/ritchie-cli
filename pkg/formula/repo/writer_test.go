package repo

import (
	"errors"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestWrite(t *testing.T) {

	type in struct {
		ritHome string
		file    stream.FileWriter
		repos   formula.Repos
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				ritHome: os.TempDir(),
				file: fileWriteMock{
					write: func(path string, byte []byte) error {
						return nil
					},
				},
				repos: formula.Repos{
					{
						Provider: "Github",
						Name:     "commons",
						Version:  "2.13.0",
						Priority: 0,
						IsLocal:  false,
					},
				},
			},
		},
		{
			name: "error",
			in: in{
				ritHome: os.TempDir(),
				file: fileWriteMock{
					write: func(path string, byte []byte) error {
						return errors.New("error to write file")
					},
				},
				repos: formula.Repos{
					{
						Provider: "Github",
						Name:     "commons",
						Version:  "2.13.0",
						Priority: 0,
						IsLocal:  false,
					},
				},
			},
			want: errors.New("error to write file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewWriter(tt.in.ritHome, tt.in.file)
			got := repo.Write(tt.in.repos)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Write(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}

}

type fileWriteMock struct {
	write func(path string, byte []byte) error
}

func (f fileWriteMock) Write(path string, byte []byte) error {
	return f.write(path, byte)
}
