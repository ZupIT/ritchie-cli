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

package rtutorial

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {

	tests := []struct {
		name string
		in   string
	}{
		{
			name: "Set enabled tutorial",
			in:   "enabled",
		},
		{
			name: "Set disabled tutorial",
			in:   "disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := os.TempDir()
			tmpTutorial := filepath.Join(tmp, TutorialFile)
			defer os.Remove(tmpTutorial)

			setter := NewSetter(tmp)
			got, err := setter.Set(tt.in)
			assert.NoError(t, err)

			expected := TutorialHolder{Current: tt.in}
			assert.Equal(t, expected, got)

			file, err := ioutil.ReadFile(tmpTutorial)
			assert.NoError(t, err)

			var holder = TutorialHolder{}
			err = json.Unmarshal(file, &holder)
			assert.NoError(t, err)
			assert.Equal(t, expected, holder)
		})
	}

}
