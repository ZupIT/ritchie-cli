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

package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const loadConfigErrMsg = `Failed to load formula config file
Try running rit update repo
Config file path not found: %s`

type PreRunManager struct {
	ritchieHome string
	make        formula.MakeBuilder
	docker      formula.DockerBuilder
	bat         formula.BatBuilder
	dir         stream.DirCreateListCopyRemover
	file        stream.FileReadExister
}

func NewPreRun(
	ritchieHome string,
	make formula.MakeBuilder,
	docker formula.DockerBuilder,
	bat formula.BatBuilder,
	dir stream.DirCreateListCopyRemover,
	file stream.FileReadExister,
) PreRunManager {
	return PreRunManager{
		ritchieHome: ritchieHome,
		make:        make,
		docker:      docker,
		bat:         bat,
		dir:         dir,
		file:        file,
	}
}

func (pr PreRunManager) PreRun(def formula.Definition, docker bool) (formula.Setup, error) {
	pwd, _ := os.Getwd()
	formulaPath := def.FormulaPath(pr.ritchieHome)

	config, err := pr.loadConfig(formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	binFilePath := def.BinFilePath(formulaPath)
	if !pr.file.Exists(binFilePath) {
		s := spinner.StartNew("Building formula...")
		time.Sleep(2 * time.Second)
		if err := pr.buildFormula(formulaPath, config.DockerIB, docker); err != nil {
			s.Stop()

			// Remove /bin dir to force formula rebuild in next execution
			if err := pr.dir.Remove(def.BinPath(formulaPath)); err != nil {
				return formula.Setup{}, err
			}

			return formula.Setup{}, err
		}
		s.Success(prompt.Green("Formula was successfully built!"))
	}

	tmpDir, err := pr.createWorkDir(pr.ritchieHome, formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	if err := os.Chdir(tmpDir); err != nil {
		return formula.Setup{}, err
	}

	s := formula.Setup{
		Pwd:         pwd,
		FormulaPath: formulaPath,
		BinName:     def.BinName(),
		BinPath:     def.BinPath(formulaPath),
		TmpDir:      tmpDir,
		Config:      config,
	}

	dockerFile := filepath.Join(tmpDir, "Dockerfile")
	if docker && validateDocker() && pr.file.Exists(dockerFile) {
		s.ContainerId, err = buildRunImg(def)
		if err != nil {
			return formula.Setup{}, err
		}
	}

	return s, nil
}

func (pr PreRunManager) buildFormula(formulaPath, dockerIB string, dockerFlag bool) error {
	if dockerFlag && dockerIB != "" && validateDocker() { // Build formula inside docker
		if err := pr.docker.Build(formulaPath, dockerIB); err != nil {
			fmt.Println("\n" + err.Error())
		} else {
			return nil
		}
	}

	if runtime.GOOS == osutil.Windows { // Build formula local with build.bat
		if err := pr.bat.Build(formulaPath); err != nil {
			return err
		}
		return nil
	}

	if err := pr.make.Build(formulaPath); err != nil { // Build formula local with Makefile
		return err
	}

	return nil
}

func (pr PreRunManager) loadConfig(formulaPath string, def formula.Definition) (formula.Config, error) {
	configPath := def.ConfigPath(formulaPath)
	if !pr.file.Exists(configPath) {
		return formula.Config{}, fmt.Errorf(loadConfigErrMsg, configPath)
	}

	configFile, err := pr.file.Read(configPath)
	if err != nil {
		return formula.Config{}, err
	}

	var formulaConfig formula.Config
	if err := json.Unmarshal(configFile, &formulaConfig); err != nil {
		return formula.Config{}, err
	}
	return formulaConfig, nil
}

func (pr PreRunManager) createWorkDir(home, formulaPath string, def formula.Definition) (string, error) {
	tDir := def.TmpWorkDirPath(home)
	if err := pr.dir.Create(tDir); err != nil {
		return "", err
	}

	if err := pr.dir.Copy(def.BinPath(formulaPath), tDir); err != nil {
		return "", err
	}

	return tDir, nil
}

func buildRunImg(def formula.Definition) (string, error) {
	s := spinner.StartNew("Building docker image to run formula...")
	formName := strings.ReplaceAll(def.Path, string(os.PathSeparator), "-")
	containerId := fmt.Sprintf("rit-repo-%s-formula%s", def.RepoName, formName)
	if len(containerId) > 200 {
		containerId = containerId[:200]
	}

	args := []string{"build", "-t", containerId, "."}
	cmd := exec.Command(dockerCmd, args...) // Run command "docker build -t (randomId) ."
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.Stop()
		return "", err
	}

	s.Success(prompt.Green("Docker image successfully built!"))
	return containerId, nil
}

// validate checks if able to run inside docker
func validateDocker() bool {
	args := []string{"version", "--format", "'{{.Server.Version}}'"}
	cmd := exec.Command(dockerCmd, args...)
	output, err := cmd.CombinedOutput()
	if output == nil || err != nil {
		return false
	}

	return true
}
