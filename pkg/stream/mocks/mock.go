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

package stream

type FileReaderCustomMock struct {
	ReadMock func(path string) ([]byte, error)
}

func (fmr FileReaderCustomMock) Read(path string) ([]byte, error) {
	return fmr.ReadMock(path)
}

type FileReadExisterCustomMock struct {
	ReadMock   func(path string) ([]byte, error)
	ExistsMock func(path string) bool
}

// Read of FileManagerCustomMock
func (fmc FileReadExisterCustomMock) Read(path string) ([]byte, error) {
	return fmc.ReadMock(path)
}

// Exists of FileManagerCustomMock
func (fmc FileReadExisterCustomMock) Exists(path string) bool {
	return fmc.ExistsMock(path)
}

type FileWriterCustomMock struct {
	WriteMock func(path string, content []byte) error
}

func (wcm FileWriterCustomMock) Write(path string, content []byte) error {
	return wcm.WriteMock(path, content)
}

type FileWriterMock struct{}

func (FileWriterMock) Write(path string, content []byte) error {
	return nil
}

type FileWriteReadExisterMock struct{}

func (FileWriteReadExisterMock) Write(path string, content []byte) error {
	return nil
}

func (FileWriteReadExisterMock) Read(path string) ([]byte, error) {
	return []byte(""), nil
}

func (FileWriteReadExisterMock) Exists(path string) bool {
	return true
}

type FileWriteReadExisterCustomMock struct {
	WriteMock  func(path string, content []byte) error
	ReadMock   func(path string) ([]byte, error)
	ExistsMock func(path string) bool
}

func (f FileWriteReadExisterCustomMock) Read(path string) ([]byte, error) {
	return f.ReadMock(path)
}

func (f FileWriteReadExisterCustomMock) Exists(path string) bool {
	return f.ExistsMock(path)
}

func (f FileWriteReadExisterCustomMock) Write(path string, content []byte) error {
	return f.WriteMock(path, content)
}
