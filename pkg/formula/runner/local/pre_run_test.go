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
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestPreRun(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHomeName := ".rit-pre-run-local"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	repoPath := filepath.Join(ritHome, "repos", "commons")

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	var config formula.Config
	_ = json.Unmarshal([]byte(configJson), &config)

	type in struct {
		def        formula.Definition
		makeBuild  formula.MakeBuilder
		batBuild   formula.BatBuilder
		shellBuild formula.ShellBuilder
		file       stream.FileReadExister
		dir        stream.DirCreateListCopyRemover
	}

	type out struct {
		want    formula.Setup
		wantErr bool
		err     error
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
				wantErr: false,
				err:     nil,
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
				wantErr: true,
				err:     builder.ErrBuildFormulaMakefile,
			},
		},
		{
			name: "not found config error",
			in: in{
				def:  formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file: fileManagerMock{exist: false},
			},
			out: out{
				wantErr: true,
				err:     fmt.Errorf(loadConfigErrMsg, filepath.Join(tmpDir, ritHomeName, "repos", "commons", "testing", "formula", "config.json")),
			},
		},
		{
			name: "read config error",
			in: in{
				def:  formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file: fileManagerMock{exist: true, rErr: errors.New("error to read config")},
			},
			out: out{
				wantErr: true,
				err:     errors.New("error to read config"),
			},
		},
		{
			name: "unmarshal config error",
			in: in{
				def:  formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file: fileManagerMock{exist: true, rBytes: []byte("error")},
			},
			out: out{
				wantErr: true,
				err:     errors.New("invalid character 'e' looking for beginning of value"),
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
				wantErr: true,
				err:     errors.New("error to create dir"),
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
				wantErr: true,
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
				want:    formula.Setup{},
				wantErr: true,
				err:     errors.New("remove bin dir error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			_ = dirManager.Remove(filepath.Join(in.def.FormulaPath(ritHome), "bin"))
			preRun := NewPreRun(ritHome, in.makeBuild, in.batBuild, in.shellBuild, in.dir, in.file)
			got, err := preRun.PreRun(in.def)

			if tt.out.wantErr {
				if tt.out.err == nil && err == nil {
					t.Errorf("PreRun(%s) want a error", tt.name)
				}

				if tt.out.err != nil && err != nil && tt.out.err.Error() != err.Error() {
					t.Errorf("PreRun(%s) got %v, want %v", tt.name, err, tt.out.err)
				}
			}

			if !reflect.DeepEqual(tt.out.want.Config, got.Config) {
				t.Errorf("PreRun(%s) got %v, want %v", tt.name, got.Config, tt.out.want.Config)
			}

			_ = os.Chdir(got.Pwd) // Return to test folder
		})
	}
}

type makeBuildMock struct {
	build func(formulaPath string) error
}

func (ma makeBuildMock) Build(formulaPath string) error {
	return ma.build(formulaPath)
}

type batBuildMock struct {
	build func(formulaPath string) error
}

func (ba batBuildMock) Build(formulaPath string) error {
	return ba.build(formulaPath)
}

type shellBuildMock struct {
	build func(formulaPath string) error
}

func (sh shellBuildMock) Build(formulaPath string) error {
	return sh.build(formulaPath)
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

const configJson = `{
  "dockerImageBuilder": "cimg/go:1.14",
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
