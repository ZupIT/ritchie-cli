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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	defer os.RemoveAll(tmp)

	githubCred := Detail{Service: "github"}

	envFinder := env.NewFinder(tmp, fileManager)
	dirManager := stream.NewDirManager(fileManager)

	setter := NewSetter(tmp, envFinder, dirManager, fileManager)
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
			finder := NewFinder(tmp, envFinder)
			got, err := finder.Find(tt.provider)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.cred, got)
		})
	}
}
