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
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestListManager_List(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	type in struct {
		ritHome string
		file    stream.FileReadExister
	}
	tests := []struct {
		name    string
		in      in
		want    formula.Repos
		wantErr bool
	}{
		{
			name: "List with success",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-list-repo")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos"))

					repositoryFile := filepath.Join(ritHomePath, "repos", "repositories.json")

					data := `
						[
							{
								"name": "commons",
								"version": "v2.0.0",
								"url": "https://github.com/kaduartur/ritchie-formulas",
								"priority": 0
							}
						]`

					_ = fileManager.Write(repositoryFile, []byte(data))
					return ritHomePath
				}(),
				file: fileManager,
			},
			want: formula.Repos{
				{
					Name:     "commons",
					Version:  "v2.0.0",
					URL:      "https://github.com/kaduartur/ritchie-formulas",
					Token:    "",
					Priority: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail to read reposFilePath",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-list-repo-fail-json")
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos"))

					repositoryFile := filepath.Join(ritHomePath, "repos", "repositories.json")

					data := `not-valid-json`

					_ = fileManager.Write(repositoryFile, []byte(data))
					return ritHomePath
				}(),
				file: fileManager,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Return empty when file not exist",
			in: in{
				ritHome: os.TempDir(),
				file:    fileManager,
			},
			want:    formula.Repos{},
			wantErr: false,
		},
		{
			name: "Return fail when fail to read file",
			in: in{
				ritHome: os.TempDir(),
				file:    fileReadExisterMockWithErrorOnRead{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := NewLister(
				tt.in.ritHome,
				tt.in.file,
			)
			got, err := li.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type fileReadExisterMockWithErrorOnRead struct{}

func (m fileReadExisterMockWithErrorOnRead) Read(path string) ([]byte, error) {
	return nil, errors.New("some error")
}

func (m fileReadExisterMockWithErrorOnRead) Exists(path string) bool {
	return true
}
