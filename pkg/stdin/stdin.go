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
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var ErrInvalidInput = errors.New("the STDIN inputs weren't informed correctly. Check the JSON used to execute the command")

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return ErrInvalidInput
	}

	return nil
}

// This function aims to facilitate stdin manipulations especially on tests
// inputs:
// 		content string: the string content for the stdin
// outputs:
//		tmpfile File: the tmp file read by stdin
// 		oldstdin File: the original stdin being overwritten
// NOTICE: do not forget to remove the tmpfile and restore the original stdin, use the commands
// 		defer func() { os.Stdin = oldStdin }()
// 		defer os.Remove(tmpfile.Name())
func WriteToStdin(content string) (tmpfile *os.File, oldStdin *os.File, error error) {
	fileContent := []byte(content)
	tmpfile, err := ioutil.TempFile("", "stdin")
	if err != nil {
		return nil, nil, err
	}

	if _, err := tmpfile.Write(fileContent); err != nil {
		return nil, nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, nil, err
	}

	oldstdin := os.Stdin
	os.Stdin = tmpfile

	return tmpfile, oldstdin, nil
}
