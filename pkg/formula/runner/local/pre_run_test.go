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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestPreRun(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHomeName := ".rit-pre-run-local"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	reposPath := filepath.Join(ritHome, "repos")
	repoPath := filepath.Join(reposPath, "commons")
	repoPathOutdated := filepath.Join(reposPath, "commonsUpdated")

	defer os.RemoveAll(ritHome)
	_ = dirManager.Remove(ritHome)

	createSaved := func(path string) {
		_ = dirManager.Remove(path)
		_ = dirManager.Create(path)
	}
	createSaved(repoPath)
	createSaved(repoPathOutdated)

	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	zipRepositories := filepath.Join("..", "..", "..", "..", "testdata", "repositories.zip")
	_ = streams.Unzip(zipFile, repoPath)
	_ = streams.Unzip(zipFile, repoPathOutdated)
	_ = streams.Unzip(zipRepositories, reposPath)

	var config formula.Config
	_ = json.Unmarshal([]byte(configJSON), &config)

	configWithLatestTagRequired := config
	configWithLatestTagRequired.RequireLatestVersion = true

	repoLister := repo.NewLister(ritHome, fileManager)
	preRunChecker := runner.NewPreRunBuilderChecker(repoLister)

	type in struct {
		def        formula.Definition
		makeBuild  formula.Builder
		batBuild   formula.Builder
		shellBuild formula.Builder
		file       stream.FileReadExister
		dir        stream.DirCreateListCopyRemover
	}

	type out struct {
		want             formula.Setup
		err              error
		NonSpecificError bool
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "local build success",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file: fileManager,
				dir:  dirManager,
			},
			out: out{
				want: formula.Setup{
					Config: config,
				},
			},
		},
		{
			name: "local build error",
			in: in{
				def: formula.Definition{Path: "testing/without-build-sh", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return builder.ErrBuildFormulaMakefile
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return builder.ErrBuildFormulaMakefile
					},
				},
				file: fileManager,
				dir:  dirManager,
			},
			out: out{
				err: builder.ErrBuildFormulaMakefile,
			},
		},
		{
			name: "not found config error",
			in: in{
				def:  formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file: fileManagerMock{exist: false},
			},
			out: out{
				err: fmt.Errorf(loadConfigErrMsg, filepath.Join(tmpDir, ritHomeName, "repos", "commons", "testing", "formula", "config.json")),
			},
		},
		{
			name: "read config error",
			in: in{
				def:  formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file: fileManagerMock{exist: true, rErr: errors.New("error to read config")},
			},
			out: out{
				err: errors.New("error to read config"),
			},
		},
		{
			name: "unmarshal config error",
			in: in{
				def:  formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file: fileManagerMock{exist: true, rBytes: []byte("error")},
			},
			out: out{
				err: errors.New("invalid character 'e' looking for beginning of value"),
			},
		},
		{
			name: "create work dir error",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file: fileManager,
				dir:  dirManagerMock{createErr: errors.New("error to create dir")},
			},
			out: out{
				err: errors.New("error to create dir"),
			},
		},
		{
			name: "copy work dir error",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file: fileManager,
				dir:  dirManagerMock{copyErr: errors.New("error to copy dir")},
			},
			out: out{
				err: errors.New("error to copy dir"),
			},
		},
		{
			name: "Chdir error",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file: fileManager,
				dir:  dirManagerMock{},
			},
			out: out{
				err:              errors.New("chdir /tmp/.rit-pre-run-local/tmp/"),
				NonSpecificError: true,
			},
		},
		{
			name: "local build error delete bin dir",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return builder.ErrBuildFormulaMakefile
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return builder.ErrBuildFormulaMakefile
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return builder.ErrBuildFormulaMakefile
					},
				},
				file: fileManager,
				dir:  dirManagerMock{removeErr: errors.New("remove bin dir error")},
			},
			out: out{
				err: errors.New("remove bin dir error"),
			},
		},
		{
			name: "local build success with latest version required and repository is updated",
			in: in{
				def: formula.Definition{Path: "testing/withLatestVersionRequired", RepoName: "commonsUpdated"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file: fileManager,
				dir:  dirManager,
			},
			out: out{
				want: formula.Setup{
					Config: configWithLatestTagRequired,
				},
			},
		},
		{
			name: "local build failed with latest version required and repository is outdated",
			in: in{
				def:  formula.Definition{Path: "testing/withLatestVersionRequired", RepoName: "commons"},
				file: fileManager,
				dir:  dirManager,
			},
			out: out{
				err: fmt.Errorf("This formula needs run in the latest version of the repository\n\tCurrent version: 2.15.1\n\tLatest version: 3.0.0\nRun 'rit update repo' to update."),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			_ = dirManager.Remove(filepath.Join(in.def.FormulaPath(ritHome), "bin"))
			preRun := NewPreRun(ritHome, in.makeBuild, in.batBuild, in.shellBuild, in.dir, in.file, preRunChecker)

			got, err := preRun.PreRun(in.def)

			if err != nil || tt.out.err != nil {
				if tt.out.NonSpecificError {
					assert.Contains(t, err.Error(), tt.out.err.Error())
				} else {
					assert.Equal(t, tt.out.err.Error(), err.Error())
				}
			}

			assert.Equal(t, tt.out.want.Config, got.Config)

			_ = os.Chdir(got.Pwd) // Return to test folder
		})
	}
}

func TestPreRunYAML(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHomeName := ".rit-pre-run-local"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	reposPath := filepath.Join(ritHome, "repos")
	repoPath := filepath.Join(reposPath, "commons")
	repoPathOutdated := filepath.Join(reposPath, "commonsUpdated")

	defer os.RemoveAll(ritHome)
	_ = dirManager.Remove(ritHome)

	createSaved := func(path string) {
		_ = dirManager.Remove(path)
		_ = dirManager.Create(path)
	}
	createSaved(repoPath)
	createSaved(repoPathOutdated)

	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	zipRepositories := filepath.Join("..", "..", "..", "..", "testdata", "repositories.zip")
	_ = streams.Unzip(zipFile, repoPath)
	_ = streams.Unzip(zipFile, repoPathOutdated)
	_ = streams.Unzip(zipRepositories, reposPath)

	var config formula.Config
	_ = yaml.Unmarshal([]byte(configYAML), &config)

	configWithLatestTagRequired := config
	configWithLatestTagRequired.RequireLatestVersion = true

	repoLister := repo.NewLister(ritHome, fileManager)
	preRunChecker := runner.NewPreRunBuilderChecker(repoLister)

	type in struct {
		def        formula.Definition
		makeBuild  formula.Builder
		batBuild   formula.Builder
		shellBuild formula.Builder
		file       stream.FileReadExister
		dir        stream.DirCreateListCopyRemover
	}

	type out struct {
		want             formula.Setup
		err              error
		NonSpecificError bool
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "local build success",
			in: in{
				def: formula.Definition{Path: "testing/formulayml", RepoName: "commons"},
				makeBuild: makeBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				batBuild: batBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				shellBuild: shellBuildMock{
					build: func(formulaPath string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file: fileManager,
				dir:  dirManager,
			},
			out: out{
				want: formula.Setup{
					Config: config,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			_ = dirManager.Remove(filepath.Join(in.def.FormulaPath(ritHome), "bin"))
			preRun := NewPreRun(ritHome, in.makeBuild, in.batBuild, in.shellBuild, in.dir, in.file, preRunChecker)

			got, err := preRun.PreRun(in.def)

			if err != nil || tt.out.err != nil {
				if tt.out.NonSpecificError {
					assert.Contains(t, err.Error(), tt.out.err.Error())
				} else {
					assert.Equal(t, tt.out.err.Error(), err.Error())
				}
			}

			assert.Equal(t, tt.out.want.Config, got.Config)

			_ = os.Chdir(got.Pwd) // Return to test folder
		})
	}
}

type makeBuildMock struct {
	build func(formulaPath string) error
}

func (ma makeBuildMock) Build(info formula.BuildInfo) error {
	return ma.build(info.FormulaPath)
}

type batBuildMock struct {
	build func(formulaPath string) error
}

func (ba batBuildMock) Build(info formula.BuildInfo) error {
	return ba.build(info.FormulaPath)
}

type shellBuildMock struct {
	build func(formulaPath string) error
}

func (sh shellBuildMock) Build(info formula.BuildInfo) error {
	return sh.build(info.FormulaPath)
}

type dirManagerMock struct {
	copyErr   error
	createErr error
	removeErr error
}

func (di dirManagerMock) Create(dir string) error {
	return di.createErr
}

func (di dirManagerMock) Copy(src, dst string) error {
	return di.copyErr
}

func (di dirManagerMock) List(dir string, hiddenDir bool) ([]string, error) {
	return nil, nil
}

func (di dirManagerMock) Remove(dir string) error {
	return di.removeErr
}

type fileManagerMock struct {
	rBytes []byte
	rErr   error
	wErr   error
	aErr   error
	exist  bool
}

func (fi fileManagerMock) Write(string, []byte) error {
	return fi.wErr
}

func (fi fileManagerMock) Read(string) ([]byte, error) {
	return fi.rBytes, fi.rErr
}

func (fi fileManagerMock) Exists(string) bool {
	return fi.exist
}

func (fi fileManagerMock) Append(path string, content []byte) error {
	return fi.aErr
}

const configJSON = `{
  "dockerImageBuilder": "cimg/go:1.14",
  "dockerVolumes": [
	],
  "inputs": [
    {
      "cache": {
        "active": true,
        "newLabel": "Type new value. ",
        "qty": 3
      },
      "label": "Type your name: ",
      "name": "input_text",
      "type": "text"
    },
    {
      "default": "false",
      "items": [
        "false",
        "true"
      ],
      "label": "Have you ever used Ritchie? ",
      "name": "input_boolean",
      "type": "bool"
    },
    {
      "default": "everything",
      "items": [
        "daily tasks",
        "workflows",
        "toils",
        "everything"
      ],
      "label": "What do you want to automate? ",
      "name": "input_list",
      "type": "text"
    },
    {
      "label": "Tell us a secret: ",
      "name": "input_password",
      "type": "password"
    }
  ]
}`

const configYAML = `dockerImageBuilder: cimg/go:1.14
dockerVolumes: []
inputs:
- cache:
    active: true
    newLabel: 'Type new value. '
    qty: 3
  label: 'Type your name: '
  name: input_text
  type: text
- default: 'false'
  items:
  - 'false'
  - 'true'
  label: 'Have you ever used Ritchie? '
  name: input_boolean
  type: bool
- default: everything
  items:
  - daily tasks
  - workflows
  - toils
  - everything
  label: 'What do you want to automate? '
  name: input_list
  type: text
- label: 'Tell us a secret: '
  name: input_password
  type: password`
