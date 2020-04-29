package server

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const serverFilePattern = "%s/server"

type SetterManager struct {
	serverFile string
	finder     Finder
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
	if err := fileutil.WriteFile(s.serverFile, []byte(url)); err != nil {
		return err
	}
	return nil
}
