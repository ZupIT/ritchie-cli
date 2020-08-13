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

package metrics

import (
	"path/filepath"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

var FilePath = filepath.Join(api.RitchieHomeDir(), "metrics")

type Id string

type Metric struct {
	Id        Id          `json:"metricId"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type UserId string

type OS string

type Command string

type Data struct {
	UserId       UserId `json:"userId"`
	OS           OS     `json:"operationalSystem:"`
	Command      string
	CommandError string
}

type Sender interface {
	Send(dataset Metric)
}

type UserIdGenerator interface {
	Generate() (UserId, error)
}

type Checker interface {
	Check() (bool, error)
}

type Collector interface {
	Collect() (Metric, error)
}
