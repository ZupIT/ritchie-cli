package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const (
	// ServerErrPattern error message pattern
	ServerErrPattern = "Server (%s) returned %s"
	otpUrlPattern    = "%s/otp"
)

var (
	// ErrOrgIsRequired error message for org
	ErrOrgIsRequired = prompt.NewError("Organization is required")
)

type otpResponse struct {
	Otp bool `json:"otp"`
}

type SetterManager struct {
	serverFile string
	httpClient *http.Client
}

func NewSetter(ritHomeDir string, hc *http.Client) Setter {
	return SetterManager{
		serverFile: fmt.Sprintf(serverFilePattern, ritHomeDir),
		httpClient: hc,
	}
}

func (s SetterManager) Set(cfg *Config) error {
	if cfg.Organization == "" {
		return ErrOrgIsRequired
	}

	if err := validator.IsValidURL(cfg.URL); err != nil {
		return err
	}
	cfg.URL = strings.TrimRight(cfg.URL, "/")
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(otpUrlPattern, cfg.URL), nil)
	if err != nil {
		return err
	}
	req.Header.Set(headers.XOrg, cfg.Organization)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(ServerErrPattern, cfg.URL, resp.Status)
	}

	var otpR otpResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpR); err != nil {
		return err
	}
	cfg.Otp = otpR.Otp

	cfg.PinningKey, cfg.PinningAddr, err = sslCertificationBase64(cfg.URL)
	if err != nil {
		return fmt.Errorf("error pinning SSL server, verify your server url(%s)", cfg.URL)
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := fileutil.WriteFile(s.serverFile, b); err != nil {
		return err
	}
	return nil
}

func sslCertificationBase64(url string) (cert, addr string, err error) {
	if !strings.HasPrefix(url, "https") {
		return "", "", nil
	}
	u := strings.Replace(url, "https://", "", 1)//#nosec

	s := strings.Split(strings.Split(u, "/")[0], ":")
	addr = s[0]
	switch len(s) {
	case 1:
		addr = fmt.Sprintf("%s:%s", s[0], "443")
	case 2:
		addr = fmt.Sprintf("%s:%s", s[0], s[1])
	default:
		return cert, addr, errors.New("url formatter error")
	}

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return cert, addr, err
	}
	connState := conn.ConnectionState()
	peerCert := connState.PeerCertificates[0]
	der, err := x509.MarshalPKIXPublicKey(peerCert.PublicKey)
	if err != nil {
		return cert, addr, err
	}
	return base64.StdEncoding.EncodeToString(der), addr, nil
}
