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

package stdin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

const msg = "read stdin test"

type TestReader struct {
	Test string `json:"test"`
}

type StdinFile struct {
	stdin       *os.File
	stdinWriter *os.File
}

func TestReadJson(t *testing.T) {

	// Convert interface to Json for test
	i := TestReader{Test: msg}
	jsonData, _ := json.Marshal(i)

	// Insert Json inside a new Reader (simulating os.Stdin)
	var stdin bytes.Buffer
	stdin.Write(jsonData)
	reader := bufio.NewReader(&stdin)

	tr := TestReader{}

	// ReadJson through Reader and convert to chosen interface

	if err := ReadJson(reader, &tr); err != nil {
		t.Errorf("Got error %v", err)
	}

	// Assert the decoder result is the initial message
	if msg != tr.Test {
		t.Errorf("Expected : %v but got %v", msg, tr.Test)
	}
}

func TestExistsEntry(t *testing.T) {
	var tests = []struct {
		name           string
		expectedResult bool
		inputMsg       string
	}{
		{
			name:           "return true when json data inputed",
			expectedResult: true,
			inputMsg:       "{\"sendMetrics\": true, \"runType\": \"invalid\"}\n",
		},
		{
			name:           "return false when json data not inputed",
			expectedResult: false,
			inputMsg:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StdinFile{}
			isTestWithStdin := len(tt.inputMsg) > 0

			if !StdinIsClean() {
				t.Errorf("Expects stdin to be clean before the test, but it is not.")
			}

			if isTestWithStdin {
				addEntryToSdin(s, tt.inputMsg)
			}

			input := ExistsEntry()

			if input != tt.expectedResult {
				t.Errorf("Wanted: %v, Got: %v", input, tt.expectedResult)
			}

			if isTestWithStdin {
				resetStdin(s)
			}

			if !StdinIsClean() {
				t.Errorf("Expects stdin to be clean after the test, but it is not")
			}
		})
	}
}

func StdinIsClean() bool {
	return !ExistsEntry()
}

func addEntryToSdin(s *StdinFile, entry string) {
	r, w, _ := os.Pipe()

	s.stdinWriter = w
	s.stdin = os.Stdin
	os.Stdin = r
	_, _ = s.stdinWriter.Write([]byte(entry))
}

func resetStdin(s *StdinFile) {
	s.stdinWriter.Close()
	os.Stdin = s.stdin

	s.stdinWriter = nil
	s.stdin = nil
}
