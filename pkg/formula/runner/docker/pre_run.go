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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	loadConfigErrMsg = `Failed to load formula config file
Try running rit update repo
Config file path not found: %s`
	dockerLegacyBinaryName = "docker"
	dockerModernBinaryName = "com.docker.cli"
)

var (
	ErrDockerNotInstalled = errors.New(
		"you must have the docker installed to run formulas inside it, check how to install it at: [https://docs.docker.com/get-docker]",
	)
	ErrDockerImageNotFound = errors.New(
		"config.json does not contain the \"dockerImageBuilder\" field, to run this formula with docker add a docker image name to it",
	)
	ErrDockerfileNotFound = errors.New(
		"the formula cannot be executed inside the docker, you must add a \"Dockerfile\" to execute the formula inside the docker",
	)
	ErrInvalidVolume = errors.New(
		"config.json file does not contain a valid volume to be mounted",
	)

	dockerVersionRegexp *regexp.Regexp
	once sync.Once
)

var _ formula.PreRunner = PreRunManager{}

type PreRunManager struct {
	ritchieHome string
	docker      formula.Builder
	dir         stream.DirCreateListCopyRemover
	file        stream.FileReadExister
	checker     runner.PreRunCheckerManager
}

func NewPreRun(
	ritchieHome string,
	docker formula.Builder,
	dir stream.DirCreateListCopyRemover,
	file stream.FileReadExister,
	checker runner.PreRunCheckerManager,
) PreRunManager {
	return PreRunManager{
		ritchieHome: ritchieHome,
		docker:      docker,
		dir:         dir,
		file:        file,
		checker:     checker,
	}
}

func (pr PreRunManager) PreRun(def formula.Definition) (formula.Setup, error) {
	pwd, _ := os.Getwd()
	formulaPath := def.FormulaPath(pr.ritchieHome)

	config, err := pr.loadConfig(formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	if err := pr.checker.CheckVersionCompliance(def.RepoName, config.RequireLatestVersion); err != nil {
		return formula.Setup{}, err
	}

	binFilePath := def.UnixBinFilePath(formulaPath)
	if !pr.file.Exists(binFilePath) {
		s := spinner.StartNew("Building formula...")
		time.Sleep(2 * time.Second)

		if err := pr.buildFormula(formulaPath, config.DockerIB, config.Volumes); err != nil {
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

func (pr PreRunManager) buildFormula(formulaPath, dockerImg string, dockerVolumes []string) error {
	if err := validateDocker(dockerImg); err != nil {
		return err
	}

	if err := validateVolumes(dockerVolumes); err != nil {
		return err
	}

	info := formula.BuildInfo{FormulaPath: formulaPath, DockerImg: dockerImg}
	if err := pr.docker.Build(info); err != nil {
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
	cmd := exec.Command(getDockerCmd(), args...) // Run command "docker build -t (randomId) ."
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
	cmd := exec.Command(getDockerCmd(), args...)
	output, err := cmd.CombinedOutput()
	if output == nil || err != nil {
		return ErrDockerNotInstalled
	}

	if dockerImg == "" {
		return ErrDockerImageNotFound
	}

	return nil
}

// validate checks if volumes is not null
func validateVolumes(dockerVolumes []string) error {
	for _, volume := range dockerVolumes {
		if !strings.Contains(volume, ":") {
			return ErrInvalidVolume
		}
	}
	return nil
}

func getDockerCmd() string {
	if runtime.GOOS == "linux" {
		_, isRunningInWsl1 := os.LookupEnv("IS_WSL")
		_, isRunningInWsl2 := os.LookupEnv("WSL_DISTRO_NAME")

		if isRunningInWsl1 || isRunningInWsl2 {
			return getDockerCmdBasedOnEngineVersion()
		} else {
			return dockerLegacyBinaryName
		}
	} else {
		return getDockerCmdBasedOnEngineVersion()
	}
}

func getDockerEngineVersion() string {
	once.Do(func() {
		dockerVersionRegexp = regexp.MustCompile("[0-9]*\\.[0-9]*\\.[0-9]*")
	})

	cmd := exec.Command(dockerLegacyBinaryName, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err == nil {
		// we try invoking 'docker --version' with the legacy binary; if it fails, we can assume either Docker isn't
		// installed or 'com.docker.cli' is the right binary; if it succeeds, we check the version instead of using
		// 'docker' straight away because, according to previous testing, it seems some simple commands such as
		// --version could function with some of the latest Docker Engine versions but 'com.docker.cli' is the only bin
		// that works for actual interactions with containers and so on. in the future, we may perform more extensive
		// testing and remove the elaborate regex and version checking if the invocation of 'docker --version' alone is
		// a sufficient test to know which binary name to use.
		return dockerVersionRegexp.FindString(out.String())
	} else {
		return ""
	}
}

// starting with Docker Desktop 2.5.0 (engine 19.03.13), the 'docker' executable in Windows and MacOS
// has been renamed to 'com.docker.cli'
func getDockerCmdBasedOnEngineVersion() string {
	engineVersion := getDockerEngineVersion()

	if engineVersion == "" {
		// ideally, we should not need a guard here; instead, we must not allow users to invoke the formula
		// with the Docker runner in the first place, but that solution needs more planning
		return dockerModernBinaryName
	}

	splitEngineVersion := strings.Split(engineVersion, ".")

	engineMajorVersionInt, _ := strconv.Atoi(splitEngineVersion[0])
	engineMinorVersionInt, _ := strconv.Atoi(splitEngineVersion[1])
	enginePatchVersionInt, _ := strconv.Atoi(splitEngineVersion[2])

	if engineMajorVersionInt >= 20 {
		return dockerModernBinaryName
	} else if engineMajorVersionInt <= 18 {
		return dockerLegacyBinaryName
	} else { // if 19.x.x
		if engineMinorVersionInt >= 3 && enginePatchVersionInt >= 13 {
			return dockerModernBinaryName
		} else {
			return dockerLegacyBinaryName
		}
	}
}
