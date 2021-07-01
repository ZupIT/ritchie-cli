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
	"regexp"
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
	regexCoreFlag              = regexp.MustCompile(`--docker|--local|--stdin|--version|--verbose|--default|--help`)
	//regexFlag                  = regexp.MustCompile("--(.*)=")
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
		FormulaRepo:          d.repoData(),
		Flags:                flags(),
	}

	metric := APIData{
		Id:         Id(metricId()),
		UserId:     userId,
		Os:         runtime.GOOS,
		RitVersion: ritVersion,
		Timestamp:  time.Now(),
		Data:       data,
	}
	return metric, nil
}

func (d DataCollectorManager) repoData() formula.Repo {
	repoBytes, _ := d.file.Read(
		filepath.Join(d.ritchieHomeDir, formula.ReposDir, "repositories.json"),
	)
	repos := formula.Repos{}
	_ = json.Unmarshal(repoBytes, &repos)

	for _, r := range repos {
		if string(r.Name) == RepoName && r.Token == "" {
			return r
		}
	}
	return formula.Repo{}
}

func metricId() string {
	args := os.Args
	args[0] = "rit"
	var metricID []string
	for _, element := range args {
		if !strings.Contains(element, "--") {
			metricID = append(metricID, element)
		}
	}
	return strings.Join(metricID, "_")
}

func flags() []string {
	args := os.Args
	var flags []string
	for _, element := range args {
		switch {
		case regexCoreFlag.MatchString(element):
			element = strings.Replace(element, "--", "", -1)
			flags = append(flags, element)
			// Remove comment to keep input flags keys on datas
			// case strings.Contains(element, "--"):
			// 	flag := regexFlag.FindStringSubmatch(element)[1]
			// 	flags = append(flags, flag)
		}
	}
	return flags
}
