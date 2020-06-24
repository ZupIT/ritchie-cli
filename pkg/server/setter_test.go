package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const (
	urlHttp = "http://localhost:8882"
	urlHttpError = "http://localhost:8882/server/error"
	urlHttps = "https://localhost"
	urlHttpsError = "https://localhost:9999"
)

var (
	errNoSuchHost     = fmt.Errorf("lookup %s: no such host", strings.Replace(urlHttp, "http://", "", 1))
	errNoSuchHostLong = fmt.Errorf("Get \"%s\": %s", urlHttp, errNoSuchHost)
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

func makeHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func TestSet(t *testing.T) {
	type in struct {
		cfg Config
		hc  *http.Client
	}

	tests := []struct {
		name string
		in   in
		outErr error
	}{
		{
			name: "empty organization",
			in:   in{cfg: Config{Organization: ""}},
			outErr: ErrOrgIsRequired,
		},
		{
			name: "empty serverURL",
			in:   in{cfg: Config{Organization: "org", URL: ""}},
			outErr: validator.ErrInvalidURL,
		},
		{
			name: "invalid serverURL",
			in:   in{cfg: Config{Organization: "org", URL: "invalid.server.URL"}},
			outErr: validator.ErrInvalidURL,
		},
		{
			name: "trailing slash on serverURL",
			in:   in{cfg: Config{Organization: "org", URL: fmt.Sprintf("%s/", urlHttp)}, hc: http.DefaultClient},
			outErr: nil,
		},
		{
			name: "valid serverURL http",
			in:   in{cfg: Config{Organization: "org", URL: urlHttp}, hc: http.DefaultClient},
			outErr: nil,
		},
		{
			name: "no such host error",
			in:   in{cfg: Config{Organization: "org", URL: urlHttp}, hc: newClientErrNoSuchHost()},
			outErr:  errNoSuchHostLong,
		},
		{
			name: "server error",
			in:   in{cfg: Config{Organization: "org", URL: urlHttpError}, hc: http.DefaultClient},
			outErr:    fmt.Errorf(ServerErrPattern, urlHttpError, "500 Server Error"),
		},
		{
			name: "pinning server https",
			in:   in{cfg: Config{Organization: "org", URL: urlHttps}, hc: makeHttpClient()},
			outErr: nil,
		},
		{
			name: "pinning server https error",
			in:   in{cfg: Config{Organization: "org", URL: urlHttpsError}, hc: makeHttpClient()},
			outErr: errors.New("dial tcp"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			s := NewSetter(os.TempDir(), in.hc)

			got := s.Set(in.cfg)
			if tt.outErr != nil && got == nil {
				t.Errorf("Set(%s) got %v, want %v", in.cfg, got, tt.outErr)
			}
			if got != nil && got.Error() != tt.outErr.Error() {
				if !strings.Contains(got.Error(), tt.outErr.Error()) {
					t.Errorf("Set(%s) got %v, want %v", in.cfg, got, tt.outErr)
				}
			}
		})
	}
}

