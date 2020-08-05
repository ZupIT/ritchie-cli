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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestDeleteWithSuccess(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	ritHomePath := filepath.Join(os.TempDir(), "TestDeleteManager_DeleteWithSuccess")
	repoName := "some_repo_name"
	repoPath := filepath.Join(ritHomePath, "repos", repoName)

	repoJson := formula.Repos{
		{
			Name: formula.RepoName(repoName),
		},
		{
			Name: formula.RepoName("some other repo"),
		},
	}

	_ = dirManager.Remove(ritHomePath)
	_ = dirManager.Create(ritHomePath)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)

	repoData, _ := json.Marshal(repoJson)
	repoFilePath := filepath.Join(ritHomePath, "repos", reposFileName)
	_ = fileManager.Remove(repoFilePath)
	_ = fileManager.Write(repoFilePath, repoData)

	deleter := NewDeleter(ritHomePath, fileManager, dirManager)
	err := deleter.Delete(formula.RepoName(repoName))
	if err != nil {
		t.Errorf("Delete return err %v", err)
	}

	if dirManager.Exists(repoPath) {
		t.Errorf("Repopath should not exist.")
	}

	newRepoData, err := fileManager.Read(repoFilePath)
	if err != nil {
		t.Errorf("Read repofilePath return err %v", err)
	}

	newRepoJson := formula.Repos{}
	err = json.Unmarshal(newRepoData, &newRepoJson)
	if err != nil {
		t.Errorf("Unmarshal repofilePath return err %v", err)
	}

	if len(newRepoJson) != 1 {
		t.Errorf("new repofilePath should have only not removed repo")
	}

}

func TestDeleteWhenErr(t *testing.T) {
	type in struct {
		ritHome  string
		file     stream.FileWriteReadExister
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
		{
			name: "Return err when read fail",
			in: in{
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					read: func(path string) ([]byte, error) {
						return nil, errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when fail to parse json",
			in: in{
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					read: func(path string) ([]byte, error) {
						return []byte("not a json data"), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Return err when fail to write",
			in: in{
				dir: DirCreateListCopyRemoverCustomMock{
					remove: func(dir string) error {
						return nil
					},
				},
				file: FileWriteCreatorReadExistRemover{
					read: func(path string) ([]byte, error) {
						data := `
						[
								{
									"name": "commons",
									"version": "v2.0.0",
									"url": "https://github.com/kaduartur/ritchie-formulas",
									"priority": 0
								}
						]
						`
						return []byte(data), nil
					},
					write: func(path string, content []byte) error {
						return errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := DeleteManager{
				ritHome: tt.in.ritHome,
				file:    tt.in.file,
				dir:     tt.in.dir,
			}
			if err := dm.Delete(tt.in.repoName); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
