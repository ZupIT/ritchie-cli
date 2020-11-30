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

package metric

import (
	"path/filepath"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

var (
	ServerRestURL = ""
	FilePath      = filepath.Join(api.RitchieHomeDir(), "metrics")
)

type Id string

func (i Id) String() string {
	return string(i)
}

type UserID string

func (u UserID) String() string {
	return string(u)
}

type Command struct {
	Id               string    `json:"id"`
	UserID           UserID    `json:"userId"`
	Timestamp        time.Time `json:"timestamp"`
	Command          string    `json:"command"`
	ExecutionTime    float64   `json:"executionTime"`
	Error            string    `json:"error,omitempty"`
	CommonsRepoAdded string    `json:"commonsRepoAdded,omitempty"`
}

type User struct {
	Id            UserID `json:"userId"`
	Os            string `json:"os"`
	Version       string `json:"version"`
	DefaultRunner string `json:"defaultRunner"`
	Repos         Repos  `json:"repos"`
}

type Repos []Repo

type Repo struct {
	Private bool   `json:"private"`
	URL     string `json:"url,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Metadata struct {
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}

type SendCommandDataParams struct {
	ExecutionTime float64
	Error         string
}

type Sender interface {
	SendUserState(ritVersion string)
	SendCommandData(cmd SendCommandDataParams)
}

type UserIdGenerator interface {
	Generate() UserID
}

type Checker interface {
	Check() bool
}

type Collector interface {
	CollectCommandData(commandExecutionTime float64, commandError ...string) Command
	CollectUserState(ritVersion string) User
}
