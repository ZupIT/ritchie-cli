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
)

type DataCollectorManager struct {
	userId UserIdGenerator
}

func NewDataCollector() DataCollectorManager {
	return DataCollectorManager{}
}

func (d DataCollectorManager) Collect() (Metric, error) {

	userId, err := d.userId.Generate()
	if err != nil {
		return Metric{}, err
	}

	data := Data{
		Command: Command(joinArgs(" ")),
		UserId: userId,
		OS:     OS(runtime.GOOS),
	}

	metric := Metric{
		Id:        Id(joinArgs("-")),
		Timestamp: time.Now(),
		Data:      data,
	}
	return metric, nil
}

func joinArgs(sep string) string {
	args := os.Args
	return strings.Join(args, sep)

}
