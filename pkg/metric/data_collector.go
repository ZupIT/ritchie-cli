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

package metric

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	_                Collector = DataCollectorManager{}
	CommonsRepoAdded           = ""
	RepoName                   = ""
)

type DataCollectorManager struct {
	userId         UserIdGenerator
	ritchieHomeDir string
	file           stream.FileReader
}

func NewDataCollector(
	userId UserIdGenerator,
	ritchieHomeDir string,
	file stream.FileReader,
) DataCollectorManager {
	return DataCollectorManager{
		userId:         userId,
		ritchieHomeDir: ritchieHomeDir,
		file:           file,
	}
}

func (d DataCollectorManager) Collect(
	commandExecutionTime float64,
	ritVersion string,
	commandError ...string,
) (APIData, error) {
	userId, err := d.userId.Generate()
	if err != nil {
		return APIData{}, err
	}

	commandExecutionTime = math.Round(commandExecutionTime*100) / 100

	data := Data{
		CommandError:         strings.Join(commandError, " "),
		CommonsRepoAdded:     CommonsRepoAdded,
		CommandExecutionTime: commandExecutionTime,
		FormulaRepo:          d.commandRepo(),
	}

	metric := APIData{
		Id:         Id(metricID()),
		UserId:     userId,
		Os:         runtime.GOOS,
		RitVersion: ritVersion,
		Timestamp:  time.Now(),
		Data:       data,
	}
	return metric, nil
}

func (d DataCollectorManager) CollectUserState(ritVersion string) User {
	userId, _ := d.userId.Generate()
	user := User{
		Id:            userId,
		Os:            runtime.GOOS,
		Version:       ritVersion,
		DefaultRunner: d.defaultRunner(),
		Repos:         d.userRepos(),
	}

	return user
}

func (d DataCollectorManager) defaultRunner() string {
	runnerBytes, _ := d.file.Read(
		filepath.Join(d.ritchieHomeDir, "default-formula-runner"),
	)
	if string(runnerBytes) == "0" {
		return "local"
	}
	return "docker"
}

func (d DataCollectorManager) userRepos() Repos {
	repoBytes, _ := d.file.Read(
		filepath.Join(d.ritchieHomeDir, formula.ReposDir, "repositories.json"),
	)
	repos := formula.Repos{}
	metricsRepo := Repo{}
	metricsRepos := Repos{}
	_ = json.Unmarshal(repoBytes, &repos)

	for _, r := range repos {
		metricsRepo.URL = r.URL
		metricsRepo.Name = string(r.Name)
		metricsRepo.Private = true
		if r.Token == "" {
			metricsRepo.Private = false
		}
		metricsRepos = append(metricsRepos, metricsRepo)
	}
	return metricsRepos
}

func (d DataCollectorManager) commandRepo() formula.Repo {
	repos := d.readRepos()
	for _, r := range repos {
		if string(r.Name) == RepoName && r.Token == "" {
			return r
		}
	}
	return formula.Repo{}
}

func (d DataCollectorManager) readRepos() formula.Repos {
	repoBytes, _ := d.file.Read(
		filepath.Join(d.ritchieHomeDir, formula.ReposDir, "repositories.json"),
	)
	repos := formula.Repos{}
	_ = json.Unmarshal(repoBytes, &repos)

	return repos
}

func metricID() string {
	args := os.Args
	args[0] = "rit"
	return strings.Join(args, "_")
}
