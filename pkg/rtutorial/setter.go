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

package rtutorial

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type SetterManager struct {
	tutorialFile string
}

func NewSetter(homePath string) Setter {
	return SetterManager{
		tutorialFile: fmt.Sprintf(TutorialPath, homePath),
	}
}

func (s SetterManager) Set(tutorial string) (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}

	tutorialHolder.Current = tutorial

	b, err := json.Marshal(&tutorialHolder)
	if err != nil {
		return tutorialHolder, err
	}

	if err := ioutil.WriteFile(s.tutorialFile, b, os.ModePerm); err != nil {
		return tutorialHolder, err
	}

	return tutorialHolder, nil
}
