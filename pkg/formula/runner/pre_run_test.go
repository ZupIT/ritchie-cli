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
	ritHome := filepath.Join(tmpDir, ".rit-setup")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	makeBuilder := builder.NewBuildMake()
	batBuilder := builder.NewBuildBat(fileManager)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	var config formula.Config
	_ = json.Unmarshal([]byte(configJson), &config)

	type in struct {
		def         formula.Definition
		makeBuild   formula.MakeBuilder
		batBuild    formula.BatBuilder
		dockerBuild formula.DockerBuilder
		file        stream.FileReadExister
		dir         stream.DirCreateListCopyRemover
		dockerFlag  bool
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
			name: "docker build success",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				dockerBuild: dockerBuildMock{
					build: func(formulaPath, dockerImg string) error {
						return dirManager.Create(filepath.Join(formulaPath, "bin"))
					},
				},
				file:       fileManager,
				dir:        dirManager,
				dockerFlag: true,
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
			name: "docker build fallback docker",
			in: in{
				def: formula.Definition{Path: "testing/formula", RepoName: "commons"},
				dockerBuild: dockerBuildMock{
					build: func(formulaPath, dockerImg string) error {
						return builder.ErrDockerBuild
					},
				},
				makeBuild: makeBuildMock{
					build: makeBuilder.Build,
				},
				batBuild: batBuildMock{
					build: batBuilder.Build,
				},
				file:       fileManager,
				dir:        dirManager,
				dockerFlag: true,
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
			name: "docker build success",
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
				file:       fileManager,
				dir:        dirManager,
				dockerFlag: false,
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
			name: "docker build error",
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
				file:       fileManager,
				dir:        dirManager,
				dockerFlag: false,
			},
			out: out{
				wantErr: true,
				err:     builder.ErrBuildFormulaMakefile,
			},
		},
		{
			name: "not found config error",
			in: in{
				def:        formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file:       fileManagerMock{exist: false},
				dockerFlag: false,
			},
			out: out{
				wantErr: true,
				err:     fmt.Errorf(loadConfigErrMsg, filepath.Join(os.TempDir(), ".rit-setup", "repos", "commons", "testing", "formula", "config.json")),
			},
		},
		{
			name: "read config error",
			in: in{
				def:        formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file:       fileManagerMock{exist: true, rErr: errors.New("error to read config")},
				dockerFlag: false,
			},
			out: out{
				wantErr: true,
				err:     errors.New("error to read config"),
			},
		},
		{
			name: "unmarshal config error",
			in: in{
				def:        formula.Definition{Path: "testing/formula", RepoName: "commons"},
				file:       fileManagerMock{exist: true, rBytes: []byte("error")},
				dockerFlag: false,
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
				file:       fileManager,
				dir:        dirManagerMock{createErr: errors.New("error to create dir")},
				dockerFlag: false,
			},
			out: out{
				wantErr: true,
				err:     errors.New("error to create dir"),
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
				file:       fileManager,
				dir:        dirManagerMock{copyErr: errors.New("error to copy dir")},
				dockerFlag: false,
			},
			out: out{
				wantErr: true,
				err:     errors.New("error to copy dir"),
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
				file:       fileManager,
				dir:        dirManagerMock{},
				dockerFlag: false,
			},
			out: out{
				wantErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			_ = dirManager.Remove(filepath.Join(in.def.FormulaPath(ritHome), "bin"))
			preRun := NewPreRun(ritHome, in.makeBuild, in.dockerBuild, in.batBuild, in.dir, in.file)
			got, err := preRun.PreRun(in.def, in.dockerFlag)

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

type dockerBuildMock struct {
	build func(formulaPath, dockerImg string) error
}

func (do dockerBuildMock) Build(formulaPath, dockerImg string) error {
	return do.build(formulaPath, dockerImg)
}

type batBuildMock struct {
	build func(formulaPath string) error
}

func (ba batBuildMock) Build(formulaPath string) error {
	return ba.build(formulaPath)
}

type dirManagerMock struct {
	copyErr   error
	createErr error
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
	return nil
}

const configJson = `{
    "dockerImageBuilder": "cimg/go:1.14",
    "inputs": [
      {
        "name": "sample_text",
        "type": "text",
        "default": "",
        "label": "Type : ",
        "items": null,
        "cache": {
          "active": true,
          "qty": 6,
          "newLabel": "Type new value. "
        }
      },
      {
        "name": "sample_list",
        "type": "text",
        "default": "in1",
        "label": "Pick your : ",
        "items": [
          "in_list1",
          "in_list2",
          "in_list3",
          "in_listN"
        ],
        "cache": {
          "active": false,
          "qty": 0,
          "newLabel": ""
        }
      },
      {
        "name": "sample_bool",
        "type": "bool",
        "default": "false",
        "label": "Pick: ",
        "items": [
          "false",
          "true"
        ],
        "cache": {
          "active": false,
          "qty": 0,
          "newLabel": ""
        }
      }
    ]
  }`
