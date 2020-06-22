package server

import (
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

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := fileutil.WriteFile(s.serverFile, b); err != nil {
		return err
	}
	return nil
}
