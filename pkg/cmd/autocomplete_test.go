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
	"testing"
)

func TestNewAutocompleteCmd(t *testing.T) {
	cmd := NewAutocompleteCmd()
	if cmd == nil {
		t.Errorf("NewAutocompleteCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewAutocompleteZsh(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompleteZsh(mock)
	if cmd == nil {
		t.Errorf("NewAutocompleteZsh got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewAutocompleteBash(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompleteBash(mock)
	if cmd == nil {
		t.Errorf("NewAutocompleteBash got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewAutocompleteFish(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompleteFish(mock)
	if cmd == nil {
		t.Errorf("NewAutocompleteFish got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewAutocompletePowerShell(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompletePowerShell(mock)
	if cmd == nil {
		t.Errorf("NewAutocompletePowerShell got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
