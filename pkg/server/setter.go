package server

import (
	"fmt"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const (
	serverFilePattern = "%s/server"
	// ServerErrPattern error message pattern
	ServerErrPattern = "Server (%s) returned %s"
)

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

func (s SetterManager) Set(url string) error {
	if err := validator.IsValidURL(url); err != nil {
		return err
	}
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(ServerErrPattern, url, resp.Status)
	}
	if err := fileutil.WriteFile(s.serverFile, []byte(url)); err != nil {
		return err
	}
	return nil
}
