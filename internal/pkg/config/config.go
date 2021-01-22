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

package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const File = "configs.toml"

type Reader interface {
	Read() (Configs, error)
}

type Writer interface {
	Write(configs Configs) error
}

type Configs struct {
	Language string             `toml:"language"`
	Tutorial string             `toml:"tutorial"`
	Metrics  string             `toml:"metrics"`
	RunType  formula.RunnerType `toml:"runType"`
}

type Manager struct {
	configsPath string
}

func NewManager(ritHome string) Manager {
	return Manager{
		configsPath: filepath.Join(ritHome, File),
	}
}

// Write creates or update the ritchie configuration.
// Case success the error == nil
func (m Manager) Write(configs Configs) error {
	buf := &bytes.Buffer{}
	if err := toml.NewEncoder(buf).Encode(configs); err != nil {
		return err
	}

	if err := ioutil.WriteFile(m.configsPath, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (m Manager) Read() (Configs, error) {
	c := Configs{
		Language: "en",
		Tutorial: "enabled",
		Metrics:  "yes",
		RunType:  formula.DockerRun,
	}

	if _, err := os.Stat(m.configsPath); os.IsNotExist(err) {
		return c, nil
	}

	if _, err := toml.DecodeFile(m.configsPath, &c); err != nil {
		return Configs{}, err
	}

	return c, nil
}
