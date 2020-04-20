package server

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type FindManager struct {
	serverFile string
}

func NewFinder(ritchieHomeDir string) FindManager {
	return FindManager{serverFile: fmt.Sprintf(serverFilePattern, ritchieHomeDir)}
}

func (f FindManager) Find() (string, error) {
	serverURL := ""

	if !fileutil.Exists(f.serverFile) {
		return serverURL, nil
	}

	file, err := fileutil.ReadFile(f.serverFile)
	if err != nil {
		return serverURL, err
	}

	serverURL = string(file)

	return serverURL, nil
}