package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const (
	serverFilePattern = "%s/server"
	serverDown = "please, check your server. It doesn't seem to be UP"
)

type SetterManager struct {
	serverFile string
}

func NewSetter(ritchieHomeDir string) Setter {
	return SetterManager{
		serverFile: fmt.Sprintf(serverFilePattern, ritchieHomeDir),
	}
}

func (s SetterManager) Set(url string) error {
	if err := validator.IsValidURL(url); err != nil {
		return err
	}
	resp, err := http.Get(url)
	if (err != nil) || (resp.StatusCode != http.StatusOK) {
		return fmt.Errorf(
			"%v: %w",
			"HttpStatus returned: " + resp.Status + " for URL: " + url,
			errors.New(serverDown))
	}
	if err := fileutil.WriteFile(s.serverFile, []byte(url)); err != nil {
		return err
	}
	return nil
}
