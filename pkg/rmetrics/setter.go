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

type SetterManager struct {
	metricsFile string
	fw          stream.FileWriter
}

func NewSetter(homePath string, fw stream.FileWriter) Setter {
	return SetterManager{
		metricsFile: fmt.Sprintf(MetricsPath, homePath),
		fw:          fw,
	}
}

func (s SetterManager) Set(tutorial string) (MetricsHolder, error) {
	metricsHolder := MetricsHolder{Current: DefaultMetrics}
	tutorialHolderDefault := MetricsHolder{Current: DefaultMetrics}

	metricsHolder.Current = tutorial

	b, err := json.Marshal(&metricsHolder)
	if err != nil {
		return tutorialHolderDefault, err
	}

	if err := s.fw.Write(s.metricsFile, b); err != nil {
		return tutorialHolderDefault, err
	}

	return metricsHolder, nil
}
