package credteam

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
)

func TestFinder(t *testing.T) {
	type out struct {
		err    error
		status int
		want   credential.Detail
	}
	tests := []struct {
		name string
		in   string
		out  out
	}{
		{
			name: "github",
			in:   "github",
			out: out{
				status: 200,
				want:   githubCred,
			},
		},
		{
			name: "not found",
			in:   "aws",
			out: out{
				status: 404,
				err:    ErrNotFoundCredential,
				want:   credential.Detail{},
			},
		},
		{
			name: "server error",
			in:   "gcp",
			out: out{
				status: 500,
				err:    errors.New("internal server error"),
				want:   credential.Detail{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out

			var body []byte
			if out.want.Service != "" {
				body, _ = json.Marshal(&out.want)
			} else {
				body = []byte(out.err.Error())
			}

			srv := mockServer(out.status, body)
			defer srv.Close()

			finder := NewFinder(serverFinderMock{srvURL: srv.URL}, srv.Client(), sessManager, ctxFinder)

			got, err := finder.Find(tt.in)
			fmt.Println("err: ", err)

			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out.err)
			}

			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
