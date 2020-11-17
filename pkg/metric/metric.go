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

type UserId string

func (u UserId) String() string {
	return string(u)
}

type Command struct {
	Id                string    `json:"id"`
	UserID            UserId    `json:"userId"`
	Timestamp         time.Time `json:"timestamp"`
	Command           string    `json:"command"`
	ExecutionTime     float64   `json:"executionTime"`
	Error             string    `json:"error,omitempty"`
	CommonsRepoAdded  string    `json:"commonsRepoAdded,omitempty"`
	MetricsAcceptance string    `json:"metricsAcceptance,omitempty"`
}

type User struct {
	Id            UserId `json:"userId"`
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

type Sender interface {
	Send(metric APIData)
}

type UserIdGenerator interface {
	Generate() UserId
}

type Checker interface {
	Check() bool
}

type Collector interface {
	Collect(commandExecutionTime float64, ritVersion string, commandError ...string) (APIData, error)
}
