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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	ritHome := filepath.Join(os.TempDir(), "tutorial")
	_ = os.Mkdir(ritHome, os.ModePerm)
	defer os.RemoveAll(ritHome)

	tests := []struct {
		name        string
		holderState string
		fileContent string
		err         string
	}{
		{
			name:        "With no tutorial file",
			holderState: "enabled",
		},
		{
			name:        "With existing tutorial file",
			holderState: "disabled",
			fileContent: `{"tutorial": "disabled"}`,
		},
		{
			name:        "Error reading the tutorial file",
			fileContent: "error",
			err:         "invalid character 'e' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fileContent != "" {
				err := ioutil.WriteFile(fmt.Sprintf(TutorialPath, ritHome), []byte(tt.fileContent), os.ModePerm)
				assert.NoError(t, err)
			}

			finder := NewFinder(ritHome)

			got, err := finder.Find()

			if err != nil {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Empty(t, tt.err)
				assert.Equal(t, tt.holderState, got.Current)
			}
		})
	}
}
