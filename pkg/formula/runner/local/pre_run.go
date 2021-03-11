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

package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const loadConfigErrMsg = `Failed to load formula config file
Try running rit update repo
Config file path not found: %s`

var _ formula.PreRunner = PreRunManager{}

type PreRunManager struct {
	ritchieHome string
	make        formula.Builder
	bat         formula.Builder
	shell       formula.Builder
	dir         stream.DirCreateListCopyRemover
	file        stream.FileReadExister
}

func NewPreRun(
	ritchieHome string,
	make formula.Builder,
	bat formula.Builder,
	shell formula.Builder,
	dir stream.DirCreateListCopyRemover,
	file stream.FileReadExister,
) PreRunManager {
	return PreRunManager{
		ritchieHome: ritchieHome,
		make:        make,
		bat:         bat,
		shell:       shell,
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

	if err := pr.checksLatestVersionCompliance(config.RequireLatestVersion, def.RepoName); err != nil {
		return formula.Setup{}, err
	}

	binFilePath := def.BinFilePath(formulaPath)
	if !pr.file.Exists(binFilePath) {
		s := spinner.StartNew("Building formula...")

		if err := pr.buildFormula(formulaPath); err != nil {
			s.Stop()

			// Remove /bin dir to force formula rebuild in next execution
			if err := pr.dir.Remove(def.BinPath(formulaPath)); err != nil {
				return formula.Setup{}, err
			}

			return formula.Setup{}, err
		}

		s.Success(prompt.Green("Formula was successfully built!"))
	}

	s := formula.Setup{
		Pwd:         pwd,
		FormulaPath: formulaPath,
		BinName:     def.BinName(),
		BinPath:     def.BinPath(formulaPath),
		Config:      config,
	}

	return s, nil
}

func (pr PreRunManager) buildFormula(formulaPath string) error {
	info := formula.BuildInfo{FormulaPath: formulaPath}
	if runtime.GOOS == osutil.Windows { // Build formula with build.bat
		if err := pr.bat.Build(info); err != nil {
			return err
		}
		return nil
	}

	buildPath := filepath.Join(formulaPath, "build.sh")
	if pr.file.Exists(buildPath) { // Build formula with build.sh
		if err := pr.shell.Build(info); err != nil {
			return err
		}
	} else if err := pr.make.Build(info); err != nil { // Build formula with Makefile
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

func (pr PreRunManager) checksLatestVersionCompliance(requireLatestVersion bool, repoName string) error {
	if requireLatestVersion {
		repoLister := repo.NewLister(pr.ritchieHome, pr.file)
		repos, _ := repoLister.List()
		repo, _ := repos.Get(repoName)
		if repo.Version.String() != repo.LatestVersion.String() {
			return errors.New("Version of repo installed not is the latest version available, please update the repo to run this formula.")
		}
	}

	return nil
}
