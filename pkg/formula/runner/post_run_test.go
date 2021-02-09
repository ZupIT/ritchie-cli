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
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestPostRun(t *testing.T) {
	type in struct {
		file   stream.FileNewListMoveRemover
		dir    stream.DirRemover
		setup  formula.Setup
		docker bool
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				file:   fileManagerMock{},
				dir:    dirManagerMock{},
				setup:  formula.Setup{},
				docker: false,
			},
			want: nil,
		},
		{
			name: "success docker",
			in: in{
				file:   fileManagerMock{},
				dir:    dirManagerMock{},
				setup:  formula.Setup{},
				docker: true,
			},
			want: nil,
		},
		{
			name: "error remove .env file docker",
			in: in{
				file:   fileManagerMock{rmErr: errors.New("error to remove .env file")},
				dir:    dirManagerMock{},
				setup:  formula.Setup{},
				docker: true,
			},
			want: errors.New("error to remove .env file"),
		},
		{
			name: "error list new files",
			in: in{
				file:   fileManagerMock{lErr: errors.New("error to list new files")},
				dir:    dirManagerMock{},
				setup:  formula.Setup{},
				docker: false,
			},
			want: errors.New("error to list new files"),
		},
		{
			name: "error move new files",
			in: in{
				file:   fileManagerMock{mErr: errors.New("error to move new files")},
				dir:    dirManagerMock{},
				setup:  formula.Setup{},
				docker: false,
			},
			want: errors.New("error to move new files"),
		},
		{
			name: "error remove work dir",
			in: in{
				file:   fileManagerMock{},
				dir:    dirManagerMock{rmErr: errors.New("error to remove workdir")},
				setup:  formula.Setup{},
				docker: false,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewPostRunner(tt.in.file, tt.in.dir)
			got := runner.PostRun(tt.in.setup, tt.in.docker)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("PostRun(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}

}

type dirManagerMock struct {
	rmErr error
}

func (d dirManagerMock) Remove(dir string) error {
	return d.rmErr
}

type fileManagerMock struct {
	rBytes   []byte
	rErr     error
	wErr     error
	aErr     error
	mErr     error
	rmErr    error
	lErr     error
	exist    bool
	listNews []string
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

func (fi fileManagerMock) Move(oldPath, newPath string, files []string) error {
	return fi.mErr
}

func (fi fileManagerMock) Remove(path string) error {
	return fi.rmErr
}

func (fi fileManagerMock) ListNews(oldPath, newPath string) ([]string, error) {
	return fi.listNews, fi.lErr
}
