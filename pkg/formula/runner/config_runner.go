/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"gopkg.in/yaml.v2"
)

var _ formula.ConfigRunner = ConfigManager{}

const (
	FileName         = "default-formula-runner"
	ConfigJSONFormat = "json"
	ConfigYAMLFormat = "yml"
	LoadConfigErrMsg = `failed to load formula config file
	try running rit update repo
	config file path not found: %s`
)

var ErrConfigNotFound = errors.New("you must configure your default formula execution method, run \"rit set formula-runner\" to set up")

type ConfigManager struct {
	filePath string
}

func NewConfigManager(ritHome string) ConfigManager {
	return ConfigManager{
		filePath: filepath.Join(ritHome, FileName),
	}
}

func (c ConfigManager) Create(runType formula.RunnerType) error {
	data, err := json.Marshal(runType)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.filePath, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (c ConfigManager) Find() (formula.RunnerType, error) {
	data, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return formula.DefaultRun, ErrConfigNotFound
	}

	runType, err := strconv.Atoi(string(data))
	if err != nil {
		return formula.DefaultRun, err
	}

	return formula.RunnerType(runType), nil
}

func LoadConfigs(f stream.FileReadExister, formulaPath string, def formula.Definition) (formula.Config, error) {
	configPath := def.ConfigYAMLPath(formulaPath)
	configFormat := ConfigYAMLFormat

	if !f.Exists(configPath) { // formula.yml
		configPath = def.ConfigPath(formulaPath)
		configFormat = ConfigJSONFormat
		if !f.Exists(configPath) { // config.json
			return formula.Config{}, fmt.Errorf(loadConfigErrMsg, configPath)
		}
	}

	configFile, err := f.Read(configPath)
	if err != nil {
		return formula.Config{}, err
	}

	var formulaConfig formula.Config
	if configFormat == ConfigYAMLFormat { // formula.yml
		if err := yaml.Unmarshal(configFile, &formulaConfig); err != nil {
			return formula.Config{}, err
		}
	} else if configFormat == ConfigJSONFormat { // config.json
		if err := json.Unmarshal(configFile, &formulaConfig); err != nil {
			return formula.Config{}, err
		}
	}
	return formulaConfig, nil
}
