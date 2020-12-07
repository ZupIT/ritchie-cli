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

package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestBuildFormulaCmd(t *testing.T) {
	cmd := NewBuildFormulaCmd()
	b := bytes.NewBufferString("")
	cmd.SetArgs([]string{""})
	cmd.SetOut(b)
	if err := cmd.Execute(); err != nil {
		t.Errorf(err.Error())
	}
	out, _ := ioutil.ReadAll(b)
	outString := string(out)
	if !strings.Contains(outString, deprecatedText) {
		t.Errorf(
			"Error on build formula command, expected %s but got %s",
			deprecatedText,
			outString,
		)
	}
}
