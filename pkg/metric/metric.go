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
	ServerRestURL    = ""
	ServerGrpcURL    = ""
	FilePath         = filepath.Join(api.RitchieHomeDir(), "metrics")
	CommonsRepoAdded = true
)

type Id string

func (i Id) String() string {
	return string(i)
}

type UserId string

func (u UserId) String() string {
	return string(u)
}

type APIData struct {
	Id         Id          `json:"metricId"`
	UserId     UserId      `json:"userId"`
	Timestamp  time.Time   `json:"timestamp"`
	Os         string      `json:"os"`
	RitVersion string      `json:"ritVersion"`
	Data       interface{} `json:"data"`
}

type Data struct {
	CommandError     string `json:"commandError"`
	CommonsRepoAdded bool   `json:"commonsRepoAdded"`
}

type Sender interface {
	Send(metric APIData)
}

type UserIdGenerator interface {
	Generate() (UserId, error)
}

type Checker interface {
	Check() bool
}

type Collector interface {
	Collect(ritVersion string, commandError ...string) (APIData, error)
}
