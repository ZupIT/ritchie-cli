package server

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	serverFile string
	file       stream.FileReadExister
}

func NewFinder(ritchieHomeDir string, file stream.FileReadExister) FindManager {
	return FindManager{
		serverFile: fmt.Sprintf(serverFilePattern, ritchieHomeDir),
		file:       file,
	}
}

func (f FindManager) Find() (string, error) {
	serverURL := ""

	if !f.file.Exists(f.serverFile) {
		return serverURL, nil
	}

	b, err := f.file.Read(f.serverFile)
	if err != nil {
		return serverURL, err
	}

	serverURL = string(b)

	return serverURL, nil
}
