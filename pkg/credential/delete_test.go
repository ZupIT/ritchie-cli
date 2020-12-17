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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/env"
)

func TestCredDelete(t *testing.T) {
	tmp := os.TempDir()
	defer os.RemoveAll(tmp)

	service := "github"
	credentialFolder := filepath.Join(tmp, credentialDir, env.Default)
	_ = os.MkdirAll(credentialFolder, os.ModePerm)
	credentialFile := filepath.Join(credentialFolder, service)
	envFinder := env.NewFinder(homeDir, fileManager)

	cred := Detail{
		Credential: Credential{},
		Service:    service,
	}

	tests := []struct {
		name    string
		service string
		err     string
	}{
		{
			name:    "run with success",
			service: service,
		},
		{
			name: "error on file remover",
			err:  "remove " + credentialFolder + ": directory not empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(cred)
			err := ioutil.WriteFile(credentialFile, jsonData, os.ModePerm)
			assert.NoError(t, err)

			deleteCredential := NewCredDelete(tmp, envFinder)

			err = deleteCredential.Delete(tt.service)
			if err == nil {
				assert.Empty(t, tt.err)
				assert.NoFileExists(t, credentialFile)
			} else {
				assert.EqualError(t, err, tt.err)
				assert.FileExists(t, credentialFile)
			}
		})
	}
}
