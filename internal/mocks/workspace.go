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

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/stretchr/testify/mock"
)

type WorkspaceListHasherMock struct {
	mock.Mock
}

func (w *WorkspaceListHasherMock) List() (formula.Workspaces, error) {
	args := w.Called()
	return args.Get(0).(formula.Workspaces), args.Error(1)
}

func (w *WorkspaceListHasherMock) CurrentHash(formulaPath string) (string, error) {
	args := w.Called(formulaPath)
	return args.String(0), args.Error(1)
}

func (w *WorkspaceListHasherMock) PreviousHash(formulaPath string) (string, error) {
	args := w.Called(formulaPath)
	return args.String(0), args.Error(1)
}

func (w *WorkspaceListHasherMock) UpdateHash(formulaPath string, hash string) error {
	args := w.Called(formulaPath, hash)
	return args.Error(0)
}
