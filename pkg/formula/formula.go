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

package formula

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const (
	ReposDir           = "repos"
	TmpDir             = "tmp"
	DefaultConfig      = "config.json"
	PwdEnv             = "CURRENT_PWD"
	CtxEnv             = "CONTEXT"
	VerboseEnv         = "VERBOSE_MODE"
	DockerExecutionEnv = "DOCKER_EXECUTION"
	BinUnix            = "run.sh"
	BinWindows         = "run.bat"
	BinDir             = "bin"
	EnvPattern         = "%s=%s"
)

type (
	Input struct {
		Name        string      `json:"name"`
		Type        string      `json:"type"`
		Default     string      `json:"default"`
		Label       string      `json:"label"`
		Items       Items       `json:"items"`
		Cache       Cache       `json:"cache"`
		Condition   Condition   `json:"condition"`
		Pattern     Pattern     `json:"pattern"`
		RequestInfo RequestInfo `json:"requestInfo"`
		Tutorial    string      `json:"tutorial"`
		Required    *bool       `json:"required"`
	}

	Items []string

	RequestInfo struct {
		Url      string `json:"url"`
		JsonPath string `json:"jsonPath"`
	}

	Pattern struct {
		Regex        string `json:"regex"`
		MismatchText string `json:"mismatchText"`
	}
	Cache struct {
		Active   bool   `json:"active"`
		Qty      int    `json:"qty"`
		NewLabel string `json:"newLabel"`
	}
	Condition struct {
		Variable string `json:"variable"`
		Operator string `json:"operator"`
		Value    string `json:"value"`
	}
	Create struct {
		FormulaCmd  string    `json:"formulaCmd"`
		Lang        string    `json:"lang"`
		Workspace   Workspace `json:"workspace"`
		FormulaPath string    `json:"formulaPath"`
	}

	Inputs []Input

	Config struct {
		DockerIB string `json:"dockerImageBuilder"`
		Inputs   Inputs `json:"inputs"`
	}

	// Definition type that represents a Formula
	Definition struct {
		Path     string
		RepoName string
	}

	Setup struct {
		Pwd         string
		FormulaPath string
		BinName     string
		BinPath     string
		TmpDir      string
		Config      Config
		ContainerId string
	}
	Help struct {
		Short string `json:"short"`
		Long  string `json:"long"`
	}

	ExecuteData struct {
		Def     Definition
		InType  api.TermInputType
		RunType RunnerType
		Verbose bool
		Flags   *pflag.FlagSet
	}
)

type BuildInfo struct {
	FormulaPath string
	DockerImg   string
	Workspace   Workspace
}

type Creator interface {
	Create(cf Create) error
}

type Builder interface {
	Build(info BuildInfo) error
}

type CreateBuilder interface {
	Creator
	Builder
}

// FormulaPath builds the formula path from ritchie home
func (d *Definition) FormulaPath(home string) string {
	return filepath.Join(home, ReposDir, d.RepoName, d.Path)
}

// TmpWorkDirPath builds the tmp paths to run formula, first parameter is tmpDir created
// second parameter is tmpBinDir
func (d *Definition) TmpWorkDirPath(home string) string {
	u := uuid.New().String()
	return filepath.Join(home, TmpDir, u)
}

func (d *Definition) UnixBinFilePath(fPath string) string {
	return filepath.Join(fPath, BinDir, BinUnix)
}

// BinFilePath builds the bin file path from formula path
func (d *Definition) BinFilePath(fPath string) string {
	return filepath.Join(fPath, BinDir, d.BinName())
}

func (d *Definition) BinName() string {
	bName := BinUnix
	if runtime.GOOS == osutil.Windows {
		bName = BinWindows
	}
	return bName
}

// BinFilePath builds the bin file path from formula path
func (d *Definition) BinPath(fPath string) string {
	return filepath.Join(fPath, BinDir)
}

// ConfigPath builds the config path from formula path and config name
func (d *Definition) ConfigPath(formulaPath string) string {
	return filepath.Join(formulaPath, DefaultConfig)
}

// FormulaName remove rit from formulaCmd
func (c Create) FormulaCmdName() string {
	d := strings.Split(c.FormulaCmd, " ")
	return strings.Join(d[1:], " ")
}

func (c Create) PkgName() string {
	d := strings.Split(c.FormulaCmd, " ")
	return d[len(d)-1]
}

func (ii Items) Contains(item string) bool {
	return sliceutil.Contains(ii, item)
}
