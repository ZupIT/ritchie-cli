package functional

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Step struct {
	Key string `json:"key"`
	Value string `json:"value"`
	Action string `json:"action"`
}

type Scenario struct {
	Entry string `json:"entry"`
	Steps []Step `json:"steps"`
	Result string `json:"result"`
}

func commandInit(cmdIn *exec.Cmd) (stdin io.WriteCloser, err error, out io.Reader, cmd *exec.Cmd) {
	stdin, err = cmdIn.StdinPipe()
	if err != nil {
		return nil, err, nil, cmdIn
	}

	stdout, _ := cmdIn.StdoutPipe()

	err = cmdIn.Start()
	if err != nil {
		return nil, err, nil, cmdIn

	}

	return stdin, nil, stdout, cmdIn
}

func (scenario *Scenario) RunSteps() (string, error) {
	fmt.Println("Running: "+ scenario.Entry)
	steps := scenario.Steps

	cmd, stdin, err, out := funcHitRit(steps)

	for _, step := range steps {
		if step.Action == "sendkey" {
			err = funcSendKeys(step, out, err, stdin)
		} else if step.Action == "select" {
			err = funcSelect(step, out, err, stdin)
		}
	}


	defer stdin.Close()

	resp := ""
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		resp = fmt.Sprint(resp, m, "\n")
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while running: %q", err)
		os.Exit(1)
	}

	fmt.Println(resp)
	fmt.Println("--------")
	return resp, nil
}

func funcHitRit(steps []Step) (*exec.Cmd, io.WriteCloser, error, io.Reader) {
	args := strings.Fields(steps[0].Value)
	cmd := exec.Command("rit", args...)
	stdin, err, out, cmd := commandInit(cmd)
	if err != nil {
		log.Panic(err)
	}
	return cmd, stdin, err, out
}

func funcSelect(step Step, out io.Reader, err error, stdin io.WriteCloser) error {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	startKey := false
	optionNumber := 0
	for scanner.Scan() {
		m := scanner.Text()
		if strings.Contains(m, step.Key) {
			startKey = true
		}
		if startKey {
			if strings.Contains(m, step.Value) {
				err = inputCommand(err, stdin, "\n")
				break
			} else if optionNumber >= 1 {
				err = inputCommand(err, stdin, "j")
			}
			optionNumber++
		}
	}
	return err
}

func funcSendKeys(step Step, out io.Reader, err error, stdin io.WriteCloser) error {
	valueFinal := step.Value + "\n"
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	startKey := false
	//Need to work on this possibility
	// optionNumber := 0
	for scanner.Scan() {
		m := scanner.Text()
		if strings.Contains(m, step.Key) {
			startKey = true
		}
		if startKey {
			// if strings.Contains(m, "Type new value.") {
			// 	err = inputCommand(err, stdin, "\n")
				err = inputCommand(err, stdin, valueFinal)
				break
			// } else if optionNumber >= 1 {
			// 	err = inputCommand(err, stdin, "j")
			// }
			// optionNumber++
		}
	}
	return err
}

func inputCommand(err error, stdin io.WriteCloser, impcommand string) error {
	time.Sleep(1000 * time.Millisecond)
	_, err = io.Copy(stdin, bytes.NewBuffer([]byte(impcommand)))
	if err != nil {
		log.Printf("Error when giving inputs: %q", err)
		os.Exit(1)
	}
	return err
}

func LoadScenarios(file string) []Scenario {
	scaffoldJson, _ := os.Open(file)
	fmt.Println(scaffoldJson)
	defer scaffoldJson.Close()
	var result []Scenario
	byteValue, _ := ioutil.ReadAll(scaffoldJson)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}