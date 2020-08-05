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
)

var ErrInvalidInput = errors.New("the STDIN inputs weren't informed correctly. Check the JSON used to execute the command")

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return ErrInvalidInput
	}

	return nil
}
