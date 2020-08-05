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
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/color"
)

func TestHello_Run(t *testing.T) {
	type fields struct {
		Text    string
		List    string
		Boolean string
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
	}{
		{
			name: "Run with success",
			fields: fields{
				Text:    "Hello",
				List:    "World",
				Boolean: "true",
			},
			wantWriter: func() string {
				return fmt.Sprintf("Hello world!\n") +
					color.FgGreen.Render("You receive Hello in text.\n") +
					color.FgRed.Render("You receive World in list.\n") +
					color.FgYellow.Render("You receive true in boolean.\n")
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hello{
				Text:    tt.fields.Text,
				List:    tt.fields.List,
				Boolean: tt.fields.Boolean,
			}
			writer := &bytes.Buffer{}
			h.Run(writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Run() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
