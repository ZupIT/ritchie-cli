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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	home := filepath.Join(tmp, "CredFinder")
	defer os.RemoveAll(home)

	githubCred := Detail{Service: "github"}
	envFinder := env.NewFinder(home, fileManager)
	dirManager := stream.NewDirManager(fileManager)

	setter := NewSetter(home, envFinder, dirManager)
	_ = setter.Set(githubCred)

	tests := []struct {
		name     string
		provider string
		cred     Detail
		err      error
	}{
		{
			name:     "Run with success",
			provider: githubCred.Service,
			cred:     githubCred,
			err:      nil,
		},
		{
			name:     "Return err when file not exist",
			provider: "aws",
			cred:     Detail{Credential: Credential{}},
			err:      errors.New(prompt.Red(fmt.Sprintf(errNotFoundTemplate, "aws"))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := NewFinder(home, envFinder)
			got, err := finder.Find(tt.provider)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.cred, got)
		})
	}
}
