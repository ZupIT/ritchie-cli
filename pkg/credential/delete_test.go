package credential

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type ctxFinderCustomMock struct {
	findMock func() (rcontext.ContextHolder, error)
}

func (cfcm ctxFinderCustomMock) Find() (rcontext.ContextHolder, error) {
	return cfcm.findMock()
}

type fileRemoverErrorMock struct{}

func (fileRemoverErrorMock) Remove(path string) error {
	return errors.New("some error")
}

func TestCredDelete(t *testing.T) {
	type args struct {
		homePath string
		cf       rcontext.Finder
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
				cf: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
				cf: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, errors.New("ReadCredentialsValue error")
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
				cf: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
					},
				},
				fm: fileRemoverErrorMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCredDelete(tt.fields.homePath, tt.fields.cf, tt.fields.fm)
			if err := got.Delete(tt.fields.service); (err != nil) != tt.wantErr {
				t.Errorf("Delete(%s) got %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
