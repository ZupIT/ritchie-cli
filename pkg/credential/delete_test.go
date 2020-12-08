package credential

import (
	"errors"
	"testing"

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
	type args struct {
		homePath string
		env      env.Finder
		fm       stream.FileRemover
		service  string
	}
	tests := []struct {
		name    string
		wantErr bool
		fields  args
	}{
		{
			name:    "Run with success",
			wantErr: false,
			fields: args{
				homePath: "",
				service:  "",
				env: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				fm: fileManager,
			},
		},
		{
			name:    "error",
			wantErr: true,
			fields: args{
				homePath: "",
				service:  "",
				env: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, errors.New("ReadCredentialsValue error")
					},
				},
				fm: fileManager,
			},
		},
		{
			name:    "error",
			wantErr: true,
			fields: args{
				homePath: "",
				service:  "",
				env: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				fm: fileRemoverErrorMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCredDelete(tt.fields.homePath, tt.fields.env, tt.fields.fm)
			if err := got.Delete(tt.fields.service); (err != nil) != tt.wantErr {
				t.Errorf("Delete(%s) got %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
