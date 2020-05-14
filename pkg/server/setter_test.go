package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

var (
	srvListener       = "127.0.0.1:57469"
	srvURL            = fmt.Sprintf("http://%s", srvListener)
	errNoSuchHost     = fmt.Errorf("lookup %s: no such host", srvListener)
	errNoSuchHostLong = fmt.Errorf("Get \"%s\": %s", srvURL, errNoSuchHost)
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errNoSuchHost
}

func newClientErrNoSuchHost() *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{}
		}),
	}
}

func TestSet(t *testing.T) {
	type in struct {
		URL string
		hc  *http.Client
	}

	type out struct {
		status int
		err    error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "empty serverURL",
			in:   in{URL: ""},
			out: out{
				err: validator.ErrInvalidURL,
			},
		},
		{
			name: "invalid serverURL",
			in:   in{URL: "invalid.server.URL"},
			out: out{
				err: validator.ErrInvalidURL,
			},
		},
		{
			name: "valid serverURL",
			in:   in{URL: srvURL, hc: http.DefaultClient},
			out: out{
				status: 200,
			},
		},
		{
			name: "no such host error",
			in:   in{URL: srvURL, hc: newClientErrNoSuchHost()},
			out: out{
				err: errNoSuchHostLong,
			},
		},
		{
			name: "server error",
			in:   in{URL: srvURL, hc: http.DefaultClient},
			out: out{
				status: 500,
				err:    fmt.Errorf(ServerErrPattern, srvURL, "500 Internal Server Error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out
			s := NewSetter(os.TempDir(), in.hc)

			if in.URL != "" {
				srv := mockServer(out.status)
				defer srv.Close()
			}

			got := s.Set(in.URL)
			if got != nil && got.Error() != out.err.Error() {
				t.Errorf("Set(%s) got %v, want %v", in.URL, got, out.err)
			}
		})
	}
}

func mockServer(status int) *httptest.Server {
	l, err := net.Listen("tcp", srvListener)
	if err != nil {
		log.Fatal(err)
	}

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
	}))
	srv.Listener.Close()
	srv.Listener = l
	srv.Start()

	return srv
}
