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

package credential

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestCredentialResolver(t *testing.T) {
	tempDirectory := os.TempDir()
	home := filepath.Join(tempDirectory, "CredDelete")
	defer os.RemoveAll(home)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	envFinder := env.NewFinder(home, fileManager)
	credentialSetter := NewSetter(home, envFinder, dirManager)
	credentialSetterError := NewSetter(home+"/wrong_path", envFinder, dirManager)
	credentialFinder := NewFinder(home, envFinder)

	var tests = []struct {
		name       string
		credSetter Setter
		credential string
		expected   string
		err        error
	}{
		{
			name:       "resolve new provider",
			credential: "CREDENTIAL_PROVIDER_KEY",
			credSetter: credentialSetter,
			expected:   "key",
		},
		{
			name:       "resolve new key",
			credential: "CREDENTIAL_PROVIDER_KEY2",
			credSetter: credentialSetter,
			expected:   "key2",
		},
		{
			name:       "resolve existing key",
			credential: "CREDENTIAL_PROVIDER_KEY2",
			credSetter: credentialSetter,
			expected:   "key2",
		},
		{
			name:       "error to read password",
			credential: "CREDENTIAL_PROVIDER_NEW_ERROR",
			credSetter: credentialSetter,
			err:        errors.New("error to read password"),
		},
		{
			name:       "error to set credential",
			credential: "CREDENTIAL_AWS_USER",
			credSetter: credentialSetterError,
			err:        errors.New("error to set credential"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &mocks.InputPasswordMock{}
			pass.On("Password", mock.Anything, mock.Anything).Return(tt.expected, tt.err)

			credentialResolver := NewResolver(credentialFinder, tt.credSetter, pass)
			got, err := credentialResolver.Resolve(tt.credential)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
