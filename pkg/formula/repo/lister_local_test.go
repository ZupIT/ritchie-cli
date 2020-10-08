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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestListManager_ListLocal(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	type in struct {
		ritHome string
	}
	tests := []struct {
		name    string
		in      in
		want    formula.RepoName
		wantErr bool
	}{
		{
			name: "List Local with success",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-list-local-repo")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos", "local"))

					return ritHomePath
				}(),
			},
			want:    "local",
			wantErr: false,
		},
		{
			name: "Return empty when local folder not exist",
			in: in{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-list-local-repo-fail")
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos"))

					return ritHomePath
				}(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := NewListerLocal(
				tt.in.ritHome,
			)
			got, err := li.ListLocal()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLocal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLocal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
