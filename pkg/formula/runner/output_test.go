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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	tmp := os.TempDir()
	home := filepath.Join(tmp, "Output")
	defer os.RemoveAll(home)

	_ = os.MkdirAll(home, os.ModePerm)

	tests := []struct {
		name     string
		data     []string
		filePath string
		want     error
	}{
		{
			name:     "run with success",
			data:     []string{"::output test=test"},
			filePath: home,
			want:     nil,
		},
	}
	for _, tt := range tests {
		got := Output(tt.data, tt.filePath)
		filePath := filepath.Join(home, OutputFile)

		assert.Equal(t, tt.want, got)

		if tt.want == nil {
			assert.FileExists(t, filePath)
		} else {
			assert.NoFileExists(t, filePath)
		}
	}
}
