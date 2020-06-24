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
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const (
	// ServerErrPattern error message pattern
	ServerErrPattern = "Server (%s) returned %s"
	otpUrlPattern    = "%s/otp"
)

var (
	// ErrOrgIsRequired error message for org
	ErrOrgIsRequired = errors.New("Organization is required")
)

type otpResponse struct {
	Otp bool `json:"otp"`
}

type SetterManager struct {
	serverFile string
	httpClient *http.Client
	insecureSSL bool
}

func NewSetter(ritHomeDir string, hc *http.Client, i bool) Setter {
	return SetterManager{
		serverFile: fmt.Sprintf(serverFilePattern, ritHomeDir),
		httpClient: hc,
		insecureSSL: i,
	}
}

func (sm SetterManager) Set(cfg *Config) error {
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
	resp, err := sm.httpClient.Do(req)
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

	cfg.PinningKey, cfg.PinningAddr, err = sm.sslCertificationBase64(cfg.URL)
	if err != nil {
		return fmt.Errorf("error pinning SSL server, verify your server url(%s)", cfg.URL)
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := fileutil.WriteFile(sm.serverFile, b); err != nil {
		return err
	}
	return nil
}

func (sm SetterManager)sslCertificationBase64(url string) (cert, addr string, err error) {
	if !strings.HasPrefix(url, "https") {
		return "", "", nil
	}
	u := strings.Replace(url, "https://", "", 1)

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
		InsecureSkipVerify: sm.insecureSSL,
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
