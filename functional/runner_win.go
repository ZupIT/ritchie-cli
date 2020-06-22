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
)

func (scenario *Scenario) runStdinForWindows() (bytes.Buffer, error) {
	writeOutput := scenario.Steps[0].Value
	rit := strings.Fields(scenario.Steps[1].Value)
	args := append([]string{"-Command", "Write-Output", "'" + writeOutput + "'", "|", "rit"}, rit...)
	cmd := exec.Command("powershell", args...)
	_, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var b2 bytes.Buffer
	cmd.Stdout = &b2

	err := cmd.Start()
	if err != nil {
		log.Printf("Error while running: %q", err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while running: %q", err)
		b2 = stderr
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

func setUpClearSetupWindows() {
	fmt.Println("Running Clear for Windows..")
	myPath := "\\.rit\\"
	usr, _ := user.Current()
	dir := usr.HomeDir + myPath

	d, err := os.Open(dir)
	if err != nil {
		log.Printf("Error Open dir: %q", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Printf("Error Readdirnames: %q", err)
	}
	for _, name := range names {
		err := os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.Printf("Error cleaning repo rit: %q", err)
		}
	}
}
