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

package rmetrics

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	metricsFile string
	homePath    string
	fr          stream.FileReadExister
}

func NewFinder(homePath string, fr stream.FileReadExister) FindManager {
	return FindManager{
		metricsFile: fmt.Sprintf(MetricsPath, homePath),
		homePath:    homePath,
		fr:          fr,
	}
}

func (f FindManager) Find() (MetricsHolder, error) {
	metricsHolder := MetricsHolder{Current: DefaultMetrics}

	if !f.fr.Exists(f.metricsFile) {
		return metricsHolder, nil
	}

	file, err := f.fr.Read(f.metricsFile)
	if err != nil {
		return metricsHolder, err
	}

	if err := json.Unmarshal(file, &metricsHolder); err != nil {
		return metricsHolder, err
	}

	return metricsHolder, nil
}
