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

	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
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
	repoLister := repo.NewLister(ritHome, fileManager)

	reposPath := filepath.Join(ritHome, "repos")
	repoPath := filepath.Join(reposPath, "commons")
	repoPathUpdated := filepath.Join(reposPath, "commonsUpdated")
	repoPathLocal := filepath.Join(reposPath, "local")

	createSaved := func(path string) {
		_ = dirManager.Remove(path)
		_ = dirManager.Create(path)
	}
	createSaved(ritHome)
	createSaved(repoPath)
	createSaved(repoPathUpdated)
	createSaved(repoPathLocal)

	zipFile := filepath.Join("../../../testdata", "ritchie-formulas-test.zip")
	zipRepositories := filepath.Join("../../../testdata", "repositories.zip")
	_ = streams.Unzip(zipFile, repoPath)
	_ = streams.Unzip(zipFile, repoPathUpdated)
	_ = streams.Unzip(zipFile, repoPathLocal)
	_ = streams.Unzip(zipRepositories, reposPath)

	tests := []struct {
		name                 string
		repoName             string
		requireLatestVersion bool
		outErr               error
	}{
		{
			name:                 "Return nil when require latest version is true and repository is updated",
			repoName:             "commonsUpdated",
			requireLatestVersion: true,
		},
		{
			name:     "Return nil when require latest version is false",
			repoName: "commons",
		},
		{
			name:                 "Return nil when require latest version is true and repository it's local",
			repoName:             "local",
			requireLatestVersion: true,
		},
		{
			name:                 "Return error version when require latest version is true and repository is outdated",
			repoName:             "commons",
			requireLatestVersion: true,

			outErr: fmt.Errorf(ErrPreRunCheckerVersion, currentVersionCommonsInRepositoriesZip, latestVersionCommonsInRepositoriesZip),
		},
		{
			name:                 "Return error repo when require latest version is true and repository not be identify",
			repoName:             "otherRepo",
			requireLatestVersion: true,

			outErr: fmt.Errorf(ErrPreRunCheckerRepo, "otherRepo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := NewPreRunBuilderChecker(repoLister)
			err := check.CheckVersionCompliance(tt.repoName, tt.requireLatestVersion)

			assert.Equal(t, tt.outErr, err)
		})
	}
}
