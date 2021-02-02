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

type FindManager struct {
	tutorialFile string
}

func NewFinder(homePath string) FindManager {
	return FindManager{
		tutorialFile: fmt.Sprintf(TutorialPath, homePath),
	}
}

func (f FindManager) Find() (TutorialHolder, error) {
	tutorialHolder := TutorialHolder{Current: DefaultTutorial}

	if _, err := os.Stat(f.tutorialFile); os.IsNotExist(err) {
		return tutorialHolder, nil
	}

	file, err := ioutil.ReadFile(f.tutorialFile)
	if err != nil {
		return tutorialHolder, err
	}

	if err := json.Unmarshal(file, &tutorialHolder); err != nil {
		return tutorialHolder, err
	}

	return tutorialHolder, nil
}
