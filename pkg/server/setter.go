package server

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

const serverFilePattern = "%s/server"

type SetterManager struct {
	serverFile string
	file       stream.FileWriter
}

func NewSetter(ritchieHomeDir string, file stream.FileWriter) Setter {
	return SetterManager{
		serverFile: fmt.Sprintf(serverFilePattern, ritchieHomeDir),
		file:       file,
	}
}

func (s SetterManager) Set(url string) error {
	if err := validator.IsValidURL(url); err != nil {
		return err
	}
	if err := s.file.Write(s.serverFile, []byte(url)); err != nil {
		return err
	}
	return nil
}
