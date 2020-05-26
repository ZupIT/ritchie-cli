package credteam

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/server"
)

func TestSet(t *testing.T) {
	type out struct {
		status int
		err    error
	}

	tests := []struct {
		name string
		in   credential.Detail
		out  out
	}{
		{
			name: "github",
			in:   githubCred,
			out: out{
				status: 201,
			},
		},
		{
			name: "server error",
			out: out{
				status: 500,
				err:    errors.New("internal server error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out

			var body []byte
			if in.Service != "" {
				body, _ = json.Marshal(&in)
			} else {
				body = []byte(out.err.Error())
			}

			srv := mockServer(out.status, body)
			defer srv.Close()

			srvFinder := serverFinderMock{Config: server.Config{URL: srv.URL}}
			setter := NewSetter(srvFinder, srv.Client(), sessManager, ctxFinder)

			err := setter.Set(in)
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Set(%s) got %v, want %v", tt.name, err, out.err)
			}
		})
	}
}
