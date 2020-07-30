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

package fileutil

type Service interface {
	ReadFile(path string) ([]byte, error)
	WriteFilePerm(path string, content []byte, perm int32) error
}

type DefaultService struct{}

func (s DefaultService) ReadFile(path string) ([]byte, error) {
	return ReadFile(path)
}

func (s DefaultService) WriteFilePerm(path string, content []byte, perm int32) error {
	return WriteFilePerm(path, content, perm)
}
