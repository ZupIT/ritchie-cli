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

package sliceutil

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

// Contains tells whether a contains s.
func Contains(aa []string, s string) bool {
	for _, v := range aa {
		if s == v {
			return true
		}
	}
	return false
}

// ContainsCmd tells whether a contains c.
func ContainsCmd(aa api.Commands, c api.Command) bool {
	for _, v := range aa {
		if c.Parent == v.Parent && c.Usage == v.Usage {
			return true
		}

		coreCmd := fmt.Sprintf("%s_%s", v.Parent, v.Usage)
		if c.Parent == coreCmd { // Ensures that no core commands will be added
			return true
		}
	}
	return false
}

func Remove(ss []string, r string) []string {
	for i, s := range ss {
		if s == r {
			return append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}
