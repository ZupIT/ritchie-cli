package functional

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func (scenario *Scenario) runStdinForWindows() (bytes.Buffer, error) {
	writeOutput := scenario.Steps[0].Value
	rit := strings.Fields(scenario.Steps[1].Value)
	args := append([]string{"-Command", "Write-Output", "'" + writeOutput + "'", "|", "rit"}, rit...)
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
	return b2, err
}

func setUpRitWin() {
	fmt.Println("Running Setup for Windows..")
	args := []string{"-Command", "Write-Output", "'{\"passphrase\":\"test\"}'", "|", "rit", "init", "--stdin"}
	cmd := exec.Command("powershell", args...)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error when input number: %q", err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
