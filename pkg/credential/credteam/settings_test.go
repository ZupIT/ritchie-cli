package credteam

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFields(t *testing.T) {
	tmp := os.TempDir()
	fileManager := stream.NewFileManager()
	serverSetter := server.NewSetter(tmp, fileManager)
	serverFinder := server.NewFinder(tmp, fileManager)

	type out struct {
		err    error
		status int
		want   credential.Fields
	}

	tests := []struct {
		name string
		out  out
	}{
		{
			name: "github, aws",
			out: out{
				status: 200,
				want: credential.Fields{
					"github": []credential.Field{
						{
							Name: "username",
							Type: "text",
						},
						{
							Name: "token",
							Type: "password",
						},
					},
					"aws": []credential.Field{
						{
							Name: "accessKeyId",
							Type: "text",
						},
						{
							Name: "secretAccessKey",
							Type: "password",
						},
					},
				},
			},
		},
		{
			name: "not found",
			out: out{
				err:    ErrFieldsNotFound,
				status: 404,
			},
		},
		{
			name: "server error",
			out: out{
				err:    errors.New("internal server error"),
				status: 500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out

			var body []byte
			if out.want != nil {
				body, _ = json.Marshal(&out.want)
			} else {
				body = []byte(out.err.Error())
			}

			s := mockServer(out.status, body)
			_ = serverSetter.Set(s.URL)
			defer s.Close()
			settings := NewSettings(serverFinder, s.Client(), sessManager, ctxFinder)

			got, err := settings.Fields()
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Fields(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Fields(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
