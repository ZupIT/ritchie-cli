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

package runner

import (
	"os/exec"
	"testing"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestInputResolver(t *testing.T) {
	types := formula.TermInputTypes{
		api.Prompt: inputRunnerMock{},
		api.Flag:   inputRunnerMock{},
	}

	inputResolver := NewInputResolver(types)

	tests := []struct {
		name   string
		inType api.TermInputType
		want   error
	}{
		{
			name:   "success",
			inType: api.Prompt,
			want:   nil,
		},
		{
			name:   "invalid",
			inType: -1,
			want:   ErrInputNotRecognized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := inputResolver.Resolve(tt.inType)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Resolve(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

type inputRunnerMock struct {
}

func (inputRunnerMock) Inputs(cmd *exec.Cmd, setup formula.Setup, flags *pflag.FlagSet) error {
	return nil
}
