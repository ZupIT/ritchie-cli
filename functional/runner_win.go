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
)

func (scenario *Scenario) runStdinForWindows() (bytes.Buffer, error) {
	writeOutput := scenario.Steps[0].Value
	rit := strings.Fields(scenario.Steps[1].Value)
	args := append([]string{"-Command", "Write-Output", "'" + writeOutput + "'", "|", "rit"}, rit...)
	cmd := exec.Command("powershell", args...)
	_, pipeWriter := io.Pipe()
	defer pipeWriter.Close()

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

	fmt.Println(&b2)
	fmt.Println("--------")
	return b2, err
}

func setUpRitSingleWin() {
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
