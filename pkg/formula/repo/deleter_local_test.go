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

package repo

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestDeleteLocalWithSuccess(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	ritHomePath := filepath.Join(os.TempDir(), "TestDeleteManager_DeleteWithSuccess")
	repoName := "local"
	repoPath := filepath.Join(ritHomePath, "repos", repoName)

	_ = dirManager.Remove(ritHomePath)
	_ = dirManager.Create(ritHomePath)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)

	deleter := NewLocalDeleter(ritHomePath, dirManager)
	err := deleter.DeleteLocal()
	if err != nil {
		t.Errorf("Delete return err %v", err)
	}

	if dirManager.Exists(repoPath) {
		t.Errorf("Repopath should not exist.")
	}
}

func TestDeleteLocalWhenErr(t *testing.T) {
	type in struct {
		ritHome  string
		dir      stream.DirRemover
		repoName formula.RepoName
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Return err when remove fail",
			in: in{
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := DeleteLocalManager{
				ritHome: tt.in.ritHome,
				dir:     tt.in.dir,
			}
			if err := dm.DeleteLocal(); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
