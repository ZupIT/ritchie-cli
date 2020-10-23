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

package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	loadConfigErrMsg = `Failed to load formula config file
Try running rit update repo
Config file path not found: %s`
	dockerCmd = "docker"
)

var (
	ErrDockerNotInstalled = errors.New("you must have the docker installed to run formulas inside it," +
		" check how to install it at: [https://docs.docker.com/get-docker]")
	ErrDockerImageNotFound = errors.New("config.json does not contain the" +
		" \"dockerImageBuilder\" field, to run this formula with docker add a docker image name to it")
	ErrDockerfileNotFound = errors.New("the formula cannot be executed inside the docker," +
		" you must add a \"Dockerfile\" to execute the formula inside the docker")
)

var _ formula.PreRunner = PreRunManager{}

type PreRunManager struct {
	ritchieHome string
	docker      formula.DockerBuilder
	dir         stream.DirCreateListCopyRemover
	file        stream.FileReadExister
}

func NewPreRun(
	ritchieHome string,
	docker formula.DockerBuilder,
	dir stream.DirCreateListCopyRemover,
	file stream.FileReadExister,
) PreRunManager {
	return PreRunManager{
		ritchieHome: ritchieHome,
		docker:      docker,
		dir:         dir,
		file:        file,
	}
}

func (pr PreRunManager) PreRun(def formula.Definition) (formula.Setup, error) {
	pwd, _ := os.Getwd()
	formulaPath := def.FormulaPath(pr.ritchieHome)

	config, err := pr.loadConfig(formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	binFilePath := def.UnixBinFilePath(formulaPath)
	if !pr.file.Exists(binFilePath) {
		s := spinner.StartNew("Building formula...")
		time.Sleep(2 * time.Second)

		if err := pr.buildFormula(formulaPath, config.DockerIB); err != nil {
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
	if !pr.file.Exists(dockerFile) {
		return formula.Setup{}, ErrDockerfileNotFound
	}

	s.ContainerId, err = buildRunImg(def)
	if err != nil {
		return formula.Setup{}, err
	}

	return s, nil
}

func (pr PreRunManager) buildFormula(formulaPath, dockerImg string) error {
	if err := validateDocker(dockerImg); err != nil {
		return err
	}

	if err := pr.docker.Build(formulaPath, dockerImg); err != nil {
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

	binPath := def.BinPath(formulaPath)
	if err := pr.dir.Copy(binPath, tDir); err != nil {
		return "", err
	}

	return tDir, nil
}

func buildRunImg(def formula.Definition) (string, error) {
	prompt.Info("Docker image build started")
	formName := strings.ReplaceAll(def.Path, string(os.PathSeparator), "-")
	containerId := fmt.Sprintf("rit-repo-%s-formula%s", def.RepoName, formName)
	if len(containerId) > 200 {
		containerId = containerId[:200]
	}

	metric.RepoName = def.RepoName

	containerId = strings.ToLower(containerId)
	args := []string{"build", "-t", containerId, "."}
	//nolint:gosec
	cmd := exec.Command(dockerCmd, args...) // Run command "docker build -t (randomId) ."
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	prompt.Success("Docker image successfully built!")
	return containerId, nil
}

// validate checks if able to run inside docker
func validateDocker(dockerImg string) error {
	args := []string{"version", "--format", "'{{.Server.Version}}'"}
	cmd := exec.Command(dockerCmd, args...)
	output, err := cmd.CombinedOutput()
	if output == nil || err != nil {
		return ErrDockerNotInstalled
	}

	if dockerImg == "" {
		return ErrDockerImageNotFound
	}

	return nil
}
