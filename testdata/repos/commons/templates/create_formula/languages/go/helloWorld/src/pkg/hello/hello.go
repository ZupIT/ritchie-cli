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

package hello

import (
	"fmt"
	"github.com/gookit/color"
	"io"
)

type Hello struct {
	Text    string
	List    string
	Boolean string
}

func (h Hello) Run(writer io.Writer) {
	var result string
	result += fmt.Sprintf("Hello world!\n")
	result += color.FgGreen.Render(fmt.Sprintf("You receive %s in text.\n", h.Text))
	result += color.FgRed.Render(fmt.Sprintf("You receive %s in list.\n", h.List))
	result += color.FgYellow.Render(fmt.Sprintf("You receive %s in boolean.\n", h.Boolean))

	if _, err := fmt.Fprintf(writer, result); err != nil {
		panic(err)
	}
}
