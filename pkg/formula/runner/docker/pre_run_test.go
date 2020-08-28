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
	ritHomeName := ".rit-pre-run-docker"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	repoPath := filepath.Join(ritHome, "repos", "commons")
	dockerBuilder := builder.NewBuildDocker()

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	var config, invalidConfig formula.Config
	_ = json.Unmarshal([]byte(configJson), &config)
	_ = json.Unmarshal([]byte(invalidConfigJson), &invalidConfig)

	type in struct {
		def         formula.Definition
		dockerBuild formula.DockerBuilder
		file        stream.FileReadExister
		dir         stream.DirCreateListCopyRemover
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
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				dockerBuild: dockerBuilder,
				file:        fileManager,
				dir:         dirManager,
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
				dockerBuild: dockerBuildMock{
					build: func(formulaPath, dockerImg string) error {
						return builder.ErrDockerBuild
					},
				},
				file: fileManager,
				dir:  dirManager,
			},
			out: out{
				want:    formula.Setup{},
				wantErr: true,
				err:     nil,
			},
		},
		{
			name: "error not found dockerImageBuilder field in config.json",
			in: in{
				def:         formula.Definition{Path: "testing/without-dockerimg", RepoName: "commons"},
				dockerBuild: dockerBuilder,
				file:        fileManager,
				dir:         dirManager,
			},
			out: out{
				want:    formula.Setup{},
				wantErr: true,
				err:     ErrDockerImageNotFound,
			},
		},
		{
			name: "error remove bin dir when build formula fail",
			in: in{
				def:         formula.Definition{Path: "testing/without-dockerimg", RepoName: "commons"},
				dockerBuild: dockerBuilder,
				file:        fileManager,
				dir:         dirManagerMock{rmErr: errors.New("error to remove bin dir")},
			},
			out: out{
				want:    formula.Setup{},
				wantErr: true,
				err:     errors.New("error to remove bin dir"),
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
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				dockerBuild: dockerBuilder,
				file:        fileManager,
				dir:         dirManagerMock{createErr: errors.New("error to create dir")},
			},
			out: out{
				wantErr: true,
				err:     errors.New("error to create dir"),
			},
		},
		{
			name: "copy work dir error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				dockerBuild: dockerBuilder,
				file:        fileManager,
				dir:         dirManagerMock{copyErr: errors.New("error to copy dir")},
			},
			out: out{
				wantErr: true,
				err:     errors.New("error to copy dir"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			_ = dirManager.Remove(filepath.Join(in.def.FormulaPath(ritHome), "bin"))
			preRun := NewPreRun(ritHome, in.dockerBuild, in.dir, in.file)
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

type dockerBuildMock struct {
	build func(formulaPath, dockerImg string) error
}

func (do dockerBuildMock) Build(formulaPath, dockerImg string) error {
	return do.build(formulaPath, dockerImg)
}

type dirManagerMock struct {
	copyErr   error
	createErr error
	rmErr     error
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
	return di.rmErr
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

const (
	configJson = `{
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
	invalidConfigJson = `{
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
)
