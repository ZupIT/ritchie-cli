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

package mocks

import "github.com/stretchr/testify/mock"

type DirManagerMock struct {
	mock.Mock
}

func (d *DirManagerMock) Remove(dir string) error {
	args := d.Called(dir)
	return args.Error(0)
}

type FileManagerMock struct {
	mock.Mock
}

func (fi *FileManagerMock) Write(string, []byte) error {
	args := fi.Called()
	return args.Error(0)
}

func (fi *FileManagerMock) Read(string) ([]byte, error) {
	args := fi.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (fi *FileManagerMock) Exists(string) bool {
	args := fi.Called()
	return args.Bool(0)
}

func (fi *FileManagerMock) Append(path string, content []byte) error {
	args := fi.Called(path, content)
	return args.Error(0)
}

func (fi *FileManagerMock) Move(oldPath, newPath string, files []string) error {
	args := fi.Called(oldPath, newPath, files)
	return args.Error(0)
}

func (fi *FileManagerMock) Remove(path string) error {
	args := fi.Called(path)
	return args.Error(0)
}

func (fi *FileManagerMock) ListNews(oldPath, newPath string) ([]string, error) {
	args := fi.Called(oldPath, newPath)
	return args.Get(0).([]string), args.Error(1)
}
