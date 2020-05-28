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

	args := strings.Fields(scenario.Steps[0].Value)
	cmd, stdin, err, out := funcHitTerminal("rit", args)

	if err == nil {
		for _, step := range scenario.Steps {
			if step.Action == "sendkey" {
				err = funcSendKeys(step, out, stdin)
				if err != nil {
					break
				}
			} else if step.Action == "select" {
				err = funcSelect(step, out, stdin)
				if err != nil {
					break
				}
			}
		}
	}

	defer stdin.Close()

	resp := ""
	scanner := funcScannerTerminal(out)
	for scanner.Scan() {
		m := scanner.Text()
		funcShowTerminal(m)
		resp = fmt.Sprint(resp, m, "\n")
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while running: %q", err)
	}

	fmt.Println(resp)
	fmt.Println("--------")
	return resp, err
}

func (scenario *Scenario) RunStdin() (string, error) {
	fmt.Println("Running: "+ scenario.Entry)

	echo := strings.Fields(scenario.Steps[0].Value)
	rit := strings.Fields(scenario.Steps[1].Value)


	commandEcho := exec.Command("echo", echo...)
	commandRit := exec.Command("rit", rit...)

	pipeReader, pipeWriter := io.Pipe()
	commandEcho.Stdout = pipeWriter
	commandRit.Stdin = pipeReader

	var b2 bytes.Buffer
	commandRit.Stdout = &b2

	errorEcho := commandEcho.Start()
	if errorEcho != nil {
		log.Printf("Error while running: %q", errorEcho)
	}

	errorRit := commandRit.Start()
	if errorRit != nil {
		log.Printf("Error while running: %q", errorRit)
	}

	errorEcho = commandEcho.Wait()
	if errorEcho != nil {
		log.Printf("Error while running: %q", errorEcho)
	}

	pipeWriter.Close()

	errorRit = commandRit.Wait()
	if errorRit != nil {
		log.Printf("Error while running: %q", errorRit)
	}

	fmt.Println(&b2)
	fmt.Println("--------")
	return b2.String(), errorRit
}

func FuncValidateLoginRequired() {
	login := []string{"show", "context"}
	_, stdin, _, out := funcHitTerminal("rit", login)
	scanner := funcScannerTerminal(out)
	for scanner.Scan() {
		m := scanner.Text()
		funcShowTerminal(m)
		if strings.Contains(m, "To use this command, you need to start a session on Ritchie") {
			err := inputCommand(stdin, "12345\n")
			if err != nil {
				log.Printf("Error when input number: %q", err)
			}
			break
		}
	}
	scanner = funcScannerTerminal(out)
	for scanner.Scan() {
		m := scanner.Text()
		funcShowTerminal(m)
	}
}

func funcHitTerminal(app string, args []string) (*exec.Cmd, io.WriteCloser, error, io.Reader) {
	cmd := exec.Command(app, args...)
	stdin, err, out, cmd := commandInit(cmd)
	if err != nil {
		log.Panic(err)
	}
	return cmd, stdin, nil, out
}

func funcSelect(step Step, out io.Reader, stdin io.WriteCloser) error {
	scanner := funcScannerTerminal(out)
	startKey := false
	optionNumber := 0
	for scanner.Scan() {
		m := scanner.Text()
		funcShowTerminal(m)
		if strings.Contains(m, step.Key) {
			startKey = true
		}
		if startKey {
			if strings.Contains(m, step.Value) {
				err := inputCommand(stdin, "\n")
				if err != nil {
					return err
				}
				break
			} else if optionNumber >= 1 {
				err := inputCommand(stdin, "j")
				if err != nil {
					return err
				}
			}

			optionNumber++
		}
	}
	return nil
}

func funcSendKeys(step Step, out io.Reader, stdin io.WriteCloser) error {
	valueFinal := step.Value + "\n"
	scanner := funcScannerTerminal(out)
	startKey := false
	//Need to work on this possibility
	// optionNumber := 0
	for scanner.Scan() {
		m := scanner.Text()
		funcShowTerminal(m)
		if strings.Contains(m, step.Key) {
			startKey = true
		}
		if startKey {
			// if strings.Contains(m, "Type new value.") {
			// 	err = inputCommand(err, stdin, "\n")
				err := inputCommand(stdin, valueFinal)
				if err != nil {
					return err
				}
				break
			// } else if optionNumber >= 1 {
			// 	err = inputCommand(err, stdin, "j")
			// }
			// optionNumber++
		}
	}
	return nil
}

func inputCommand(stdin io.WriteCloser, command string) error {
	time.Sleep(1000 * time.Millisecond)
	_, err := io.Copy(stdin, bytes.NewBuffer([]byte(command)))
	if err != nil {
		log.Printf("Error when giving inputs: %q", err)
	}
	return err
}

func LoadScenarios(file string) []Scenario {
	scaffoldJson, _ := os.Open(file)
	fmt.Println(scaffoldJson)
	defer scaffoldJson.Close()
	var result []Scenario
	byteValue, _ := ioutil.ReadAll(scaffoldJson)
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		log.Printf("Error unmarshal json: %q", err)
		os.Exit(1)
	}
	return result
}

func funcScannerTerminal(out io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	return scanner
}

func funcShowTerminal(message string) {
	fmt.Println(message)
}