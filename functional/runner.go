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
	"runtime"
	"strings"
	"time"

	"github.com/onsi/ginkgo"

	"github.com/ZupIT/ritchie-cli/pkg/cmd"
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
	args := strings.Fields(scenario.Steps[0].Value)
	cmd, stdin, out, err := execRit(args)

	os := runtime.GOOS
	if  os == "windows" && len(scenario.Steps) >= 2 {
		ginkgo.Skip("Scenarios with multi steps for windows doesnt work")
	}

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
	return resp, err
}

func (scenario *Scenario) RunStdin() (string, error) {
	fmt.Println("Running: " + scenario.Entry)
	os := runtime.GOOS
	if  os == "windows" {
		writeOutput := strings.ReplaceAll(scenario.Steps[0].Value, "\"","\"\"\"")
		rit := strings.Fields(scenario.Steps[1].Value)
		args := append([]string{"--%%", "Write-Output", "'"+writeOutput+"'", "|", "rit"}, rit...)
		cmd := exec.Command("powershell", args...)
		_, pipeWriter := io.Pipe()
		cmd.Stdout = pipeWriter

		var b2 bytes.Buffer
		cmd.Stdout = &b2

		err := cmd.Start()
		if err != nil {
			log.Printf("Error while running: %q", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Printf("Error while running: %q", err)
		}

		pipeWriter.Close()

		fmt.Println(&b2)
		fmt.Println("--------")

		return b2.String(), err
	} else {
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



}

func RitInit() {
	os := runtime.GOOS
	if  os == "windows" {
		fmt.Println("Running Setup for Windows..")
		args := []string{"--%%", "Write-Output", "'{\"\"\"passphrase\"\"\":\"\"\"test\"\"\"}'", "|", "rit", "init", "--stdin"}
		cmd := exec.Command("powershell", args...)

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error while running: %q", err)
		}
		fmt.Printf("%s\n", stdoutStderr)

	} else {
		fmt.Println("Running Setup for Unix..")
		command := []string{initCmd}
		_, stdin, out, _ := execRit(command)
		scanner := scannerTerminal(out)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			if strings.Contains(m, cmd.MsgPhrase) {
				err := inputCommand(stdin, "12345\n")
				if err != nil {
					log.Printf("Error when input number: %q", err)
				}
				break
			}
		}
		scanner = scannerTerminal(out)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
	}
	fmt.Println("Setup Done..")
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
	//Need to work on this possibility
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
