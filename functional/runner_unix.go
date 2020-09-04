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
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func (scenario *Scenario) runStdinForUnix() (bytes.Buffer, error) {
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
	var stderr bytes.Buffer
	commandRit.Stderr = &stderr

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
		b2 = stderr
	}

	fmt.Println(&b2)
	fmt.Println("--------")
	return b2, errorRit
}

func setUpRitSingleUnix() {
	fmt.Println("Running Setup for Unix..")

	fmt.Println("Running INIT")
	initStepRit := Step{Key: "", Value: "init", Action: "rit"}
	initAcceptsMetrics := Step{
		Key:    "and feature use statistics and crash reports?\u001B[0m  \u001B[36m[Use arrows to move, type to filter]\u001B[0m",
		Value:  "Yes, I agree to contribute with data anonymously",
		Action: "select"}
	initAddRepo := Step{
		Key:    "Would you like to add the community repository? [https://github.com/ZupIT/ritchie-formulas]",
		Value:  "yes",
		Action: "select"}
	initSetRunType := Step{
		Key:    "Select a default formula run type:",
		Value:  "local",
		Action: "select",
	}

	init := Scenario{Entry: "Running Init", Result: "", Steps: []Step{initStepRit, initAcceptsMetrics, initAddRepo, initSetRunType}}

	err, _ := init.runStepsForUnix()
	if err != nil {
		log.Printf("Error when do init: %q", err)
	}
}

func setUpClearSetupUnix() {
	fmt.Println("Running Clear for Unix..")
	myPath := ".rit"
	usr, _ := user.Current()
	dir := filepath.Join(usr.HomeDir, myPath)

	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
}

func (scenario *Scenario) runStepsForUnix() (error, string) {
	args := strings.Fields(scenario.Steps[0].Value)
	cmd, stdin, out, err := execRit(args)
	if err == nil {
		for _, step := range scenario.Steps {
			if step.Action == "sendkey" {
				err = sendKeys(step, out, stdin)
				if err != nil {
					break
				}
			} else if step.Action == "select" {
				err = selectOption(step, out, stdin)
				if err != nil {
					break
				}
			}
		}
	}

	defer stdin.Close()

	resp := ""
	scanner := scannerTerminal(out)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		resp = fmt.Sprint(resp, m, "\n")
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while running: %q", err)
	}

	fmt.Println(resp)
	fmt.Println("--------")
	return err, resp
}

func commandInit(cmdIn *exec.Cmd) (stdin io.WriteCloser, out io.Reader, err error) {
	stdin, err = cmdIn.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	stdout, _ := cmdIn.StdoutPipe()

	err = cmdIn.Start()
	if err != nil {
		return nil, nil, err
	}

	return stdin, stdout, nil
}

func execRit(args []string) (*exec.Cmd, io.WriteCloser, io.Reader, error) {
	cmd := exec.Command(rit, args...)
	stdin, out, err := commandInit(cmd)
	if err != nil {
		log.Panic(err)
	}
	return cmd, stdin, out, err
}

func selectOption(step Step, out io.Reader, stdin io.WriteCloser) error {
	scanner := scannerTerminal(out)
	startKey := false
	optionNumber := 0
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
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

func sendKeys(step Step, out io.Reader, stdin io.WriteCloser) error {
	valueFinal := step.Value + "\n"
	scanner := scannerTerminal(out)
	startKey := false
	// Need to work on this possibility
	// optionNumber := 0
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
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
