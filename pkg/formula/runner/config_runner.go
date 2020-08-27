package runner

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strconv"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var _ formula.ConfigRunner = ConfigManager{}

const fileName = "default-formula-runner"

var ErrConfigNotFound = errors.New("you must configure your default formula execution method, run \"rit set formula-runner\" to set up")

type ConfigManager struct {
	filePath string
	file     stream.FileWriteReadExister
}

func NewConfigManager(ritHome string, file stream.FileWriteReadExister) ConfigManager {
	return ConfigManager{
		filePath: filepath.Join(ritHome, fileName),
		file:     file,
	}
}

func (c ConfigManager) Create(runType formula.RunnerType) error {
	data, err := json.Marshal(runType)
	if err != nil {
		return err
	}

	if err := c.file.Write(c.filePath, data); err != nil {
		return err
	}

	return nil
}

func (c ConfigManager) Find() (formula.RunnerType, error) {
	if !c.file.Exists(c.filePath) {
		return -1, ErrConfigNotFound
	}

	data, err := c.file.Read(c.filePath)
	if err != nil {
		return -1, err
	}

	runType, err := strconv.Atoi(string(data))
	if err != nil {
		return -1, err
	}

	return formula.RunnerType(runType), nil
}
