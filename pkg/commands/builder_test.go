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

package commands

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Run("Build commands successfully", func(t *testing.T) {
		cmd := Build()
		b := bytes.NewBufferString("")
		cmd.SetArgs([]string{"--version"})
		cmd.SetOut(b)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
		out, err := ioutil.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		outString := string(out)
		if !strings.Contains(outString, "rit version") {
			t.Fatalf("expected \"%s\" got \"%s\"", "rit version", outString)
		}
	})
}
