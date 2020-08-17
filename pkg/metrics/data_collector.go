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
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/cmd"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type DataCollectorManager struct {
	userId UserIdGenerator
	file stream.FileReadExister
}

func NewDataCollector(userId UserIdGenerator) DataCollectorManager {
	return DataCollectorManager{
		userId: userId,

	}
}

func (d DataCollectorManager) Collect(commandError ...string) (APIData, error) {

	userId, err := d.userId.Generate()
	if err != nil {
		return APIData{}, err
	}

	data := Data{
		CommandError: CommandError(strings.Join(commandError, " ")),
		Version: Version(cmd.Version),

	}

	metric := APIData{
		ID: ID(metricID()),
		UserId:       userId,
		OS:           OS(runtime.GOOS),
		Timestamp: time.Now(),
		Data:      data,

	}



	return metric, nil
}

func metricID() string {
	args := os.Args
	return strings.Join(args, "_")

}
