package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

func TestSet(t *testing.T) {

	s := NewSetter(os.TempDir())

	type out struct {
		status int
		err    error
	}

	tests := []struct {
		name string
		in   string
		out  out
	}{
		{
			name: "empty serverURL",
			in:   "",
			out: out{
				err: validator.ErrInvalidServerURL,
			},
		},
		{
			name: "existing serverURL",
			in:   "http://localhost/mocked",
			out: out{
				status: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out

			var body []byte
			var got error
			if in != "" {
				body, _ = json.Marshal(&in)
				server := mockServer(out.status, body)
				err := s.Set(server.URL)
				if err != nil {
					fmt.Sprintln("Error in set")
					return
				}
				defer server.Close()
				got = s.Set(server.URL)
			} else {
				got = s.Set(in)
			}

			if got != nil && errors.Unwrap(got) != out.err {
				t.Errorf("Set(%s) got %v, want %v", in, got, out)
			}

			if got == nil && out.status != 200 {
				t.Errorf("Set(%s) got %v, want %v", in, out.status, 200)
			}
		})
	}
}

func mockServer(status int, body []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
		_, err := rw.Write(body)
		if err != nil {
			fmt.Sprintln("Error in Write")
			return
		}
	}))
}
