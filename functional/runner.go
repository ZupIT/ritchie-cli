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
	initCmd = "init"
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

	os := runtime.GOOS
	if  os == "windows" && len(scenario.Steps) >= 2 {
		ginkgo.Skip("Scenarios with multi steps for windows doesnt work")
		return "", nil
	} else {
		err, resp := scenario.runStepsForUnix()
		return resp, err
	}
}

func (scenario *Scenario) RunStdin() (string, error) {
	fmt.Println("Running: " + scenario.Entry)
	os := runtime.GOOS
	if  os == "windows" {
		b2, err := scenario.runStdinForWindows()
		return b2.String(), err
	} else {
		b2, errorRit := scenario.runStdinForUnix()
		return b2.String(), errorRit
	}

}

func RitInit() {
	os := runtime.GOOS
	if  os == "windows" {
		setUpRitWin()
	} else {
		setUpRitUnix()
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
	err = json.Unmarshal([]byte(b), &res)
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
