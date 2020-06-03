package server

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type FindManager struct {
	serverFile string
}

func NewFinder(ritchieHomeDir string) FindManager {
	return FindManager{serverFile: fmt.Sprintf(serverFilePattern, ritchieHomeDir)}
}

func (f FindManager) Find() (Config, error) {
	cfg := Config{}

	if !fileutil.Exists(f.serverFile) {
		return cfg, nil
	}

	b, err := fileutil.ReadFile(f.serverFile)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
