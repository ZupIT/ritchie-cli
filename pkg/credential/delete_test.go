package credential

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type envFinderCustomMock struct {
	find func() (env.Holder, error)
}

func (e envFinderCustomMock) Find() (env.Holder, error) {
	return e.find()
}

type fileRemoverErrorMock struct{}

func (fileRemoverErrorMock) Remove(path string) error {
	return errors.New("some error")
}

func TestCredDelete(t *testing.T) {
	tmp := os.TempDir()
	defer os.RemoveAll(tmp)

	type args struct {
		homePath string
		env      env.Finder
		fm       stream.FileRemover
		service  string
	}
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "run with success",
		},
		{
			name: "error on env finder",
			err:  errors.New("ReadCredentialsValue error"),
		},
		{
			name: "error on file remover",
			err:  errors.New("ReadCredentialsValue error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteCredential := NewCredDelete(tt.fields.homePath, tt.fields.env, fileManager)

			err := deleteCredential.Delete(tt.fields.service)
			assert.Equal(t, tt.err, err)
			if err != nil {

			}
		})
	}
}
