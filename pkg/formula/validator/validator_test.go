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

package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspaceManagerAdd(t *testing.T) {
	tests := []struct {
		name    string
		formula string
		err     error
	}{
		{
			name:    "success",
			formula: "rit test success",
		},
		{
			name:    "error when formula empty",
			formula: "",
			err:     ErrFormulaCmdNotBeEmpty,
		},
		{
			name:    "error when formula prefix isn't ritchie",
			formula: "zup test error",
			err:     ErrFormulaCmdMustStartWithRit,
		},
		{
			name:    "error when formula doesn't have the minimum size",
			formula: "rit test",
			err:     ErrInvalidFormulaCmdSize,
		},
		{
			name:    "error when formula contains invalid characters",
			formula: "rit test @ccount",
			err:     ErrInvalidCharactersFormulaCmd,
		},
		{
			name:    "error when formula is core command",
			formula: "rit add repo",
			err:     fmt.Errorf(MsgErrFormulaCanBeCoreCommand, "add"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()

			err := validator.FormulaCommmandValidator(tt.formula)

			assert.Equal(t, tt.err, err)
		})
	}
}
