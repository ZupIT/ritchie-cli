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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
	"github.com/stretchr/testify/assert"
)

const (
	currentVersionCommonsInRepositoriesZip = "2.15.1"
	latestVersionCommonsInRepositoriesZip  = "3.0.0"
)

func TestCheckVersionCompliance(t *testing.T) {
	tmpDir := os.TempDir()
	ritHomeName := ".rit-pre-run-checker"
	ritHome := filepath.Join(tmpDir, ritHomeName)
	defer os.RemoveAll(ritHome)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	reposPath := filepath.Join(ritHome, "repos")
	repoPath := filepath.Join(reposPath, "commons")
	repoPathOutdated := filepath.Join(reposPath, "commonsOutdated")

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	_ = dirManager.Remove(repoPathOutdated)
	_ = dirManager.Create(repoPathOutdated)
	zipFile := filepath.Join("../../../testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)
	_ = streams.Unzip(zipFile, repoPathOutdated)
	zipRepositories := filepath.Join("../../../testdata", "repositories.zip")
	_ = streams.Unzip(zipRepositories, reposPath)
	type in struct {
		repoName             string
		requirelatestVersion bool
	}

	tests := []struct {
		name   string
		in     in
		outErr error
	}{
		{
			name: "Return nil when require latest version in true and repository is updated",
			in: in{
				repoName:             "commonsOutdated",
				requirelatestVersion: true,
			},
		},
		{
			name: "Return nil when require latest version in true and repository is outdated",
			in: in{
				repoName:             "commons",
				requirelatestVersion: true,
			},
			outErr: fmt.Errorf(versionError, currentVersionCommonsInRepositoriesZip, latestVersionCommonsInRepositoriesZip),
		},
		{
			name: "Return nil when require latest version in false",
			in: in{
				repoName:             "commons",
				requirelatestVersion: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := NewPreRunBuilderChecker(ritHome, tt.in.repoName, tt.in.requirelatestVersion, fileManager)
			err := check.CheckVersionCompliance()

			assert.Equal(t, tt.outErr, err)
		})
	}
}
