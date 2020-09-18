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

package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestMain(m *testing.M) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, ".rit-builder")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	os.Exit(m.Run())
}

func TestBuildShell(t *testing.T) {
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, ".rit-builder")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	buildShell := NewBuildShell()

	type in struct {
		formPath string
	}

	type out struct {
		wantErr bool
		err     error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success",
			in: in{
				formPath: filepath.Join(repoPath, "testing", "formula"),
			},
			out: out{wantErr: false},
		},
		{
			name: "shell error",
			in: in{
				formPath: repoPath,
			},
			out: out{wantErr: true},
		},
		{
			name: "Chdir error",
			in: in{
				formPath: filepath.Join(repoPath, "invalid"),
			},
			out: out{wantErr: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildShell.Build(tt.in.formPath)

			if got != nil && !tt.out.wantErr {
				t.Errorf("Run(%s) got %v, want not nil error", tt.name, got)
			}
		})
	}
}
