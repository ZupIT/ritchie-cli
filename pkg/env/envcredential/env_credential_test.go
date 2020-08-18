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

package envcredential

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestCredentialResolver(t *testing.T) {
	fileManager := stream.NewFileManager()
	tempDirectory := os.TempDir()
	contextFinder := rcontext.NewFinder(tempDirectory, fileManager)
	credentialSetter := credential.NewSetter(tempDirectory, contextFinder)
	credentialFinder := credential.NewFinder(tempDirectory, contextFinder, fileManager)

	defer os.RemoveAll(credential.File(tempDirectory, "", ""))

	var tests = []struct {
		name            string
		credentialField string
		output          string
	}{
		{
			name: "Test resolve new provider",
			credentialField: "CREDENTIAL_PROVIDER_KEY",
			output: "key",
		},
		{
			name: "Test resolve new key",
			credentialField: "CREDENTIAL_PROVIDER_KEY2",
			output: "key2",
		},
		{
			name: "Test resolve existing key",
			credentialField: "CREDENTIAL_PROVIDER_KEY2",
			output: "key2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentialResolver := NewResolver(credentialFinder, credentialSetter, passwordMock{tt.output})
			credentialValue, err := credentialResolver.Resolve(tt.credentialField)
			if err != nil {
				t.Errorf("Resolve credentials error = %v", err)
			}
			if credentialValue != tt.output {
				t.Errorf("Resolve credentials failed to retrieve. Expected %v, got %v", tt.output, credentialValue)
			}
		})
	}
}

type passwordMock struct {
	value string
}

func (pass passwordMock) Password(string) (string, error) {
	return pass.value, nil
}
