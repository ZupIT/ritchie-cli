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

package functional

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/onsi/ginkgo"
)

const (
	rit     = "rit"
	windows = runtime.GOOS == "windows"
	// initCmd = "init"
)

type Step struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Action string `json:"action"`
}

type Scenario struct {
	Entry  string `json:"entry"`
	Steps  []Step `json:"steps"`
	Result string `json:"result"`
}

func (scenario *Scenario) RunSteps() (string, error) {
	fmt.Println("Running: " + scenario.Entry)

	if windows && len(scenario.Steps) >= 2 {
		ginkgo.Skip("Scenarios with multi steps for windows doesnt work")
		return "", nil
	} else {
		err, resp := scenario.runStepsForUnix()
		return resp, err
	}
}

func (scenario *Scenario) RunStdin() (string, error) {
	fmt.Println("Running STDIN: " + scenario.Entry)
	if windows {
		b2, err := scenario.runStdinForWindows()
		return b2.String(), err
	} else {
		b2, errorRit := scenario.runStdinForUnix()
		return b2.String(), errorRit
	}

}

func RitSingleInit() {
	if windows {
		setUpRitSingleWin()
	} else {
		setUpRitSingleUnix()
	}
	fmt.Println("Setup Done..")
}

func RitClearConfigs() {
	if windows {
		setUpClearSetupWindows()
	} else {
		setUpClearSetupUnix()
	}
	fmt.Println("Setup Done..")
}

func LoadScenarios(file string) []Scenario {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatal("Error openning scenarios file:", err)
	}
	fmt.Println(jsonFile)
	defer jsonFile.Close()
	var res []Scenario
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error reading scenarios json:", err)
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Fatal("Error unmarshal json:", err)
	}
	return res
}

func scannerTerminal(out io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	return scanner
}
