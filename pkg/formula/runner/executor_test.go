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
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestExecute(t *testing.T) {

	type in struct {
		runners formula.Runners
		config  formula.ConfigRunner
		exe     formula.ExecuteData
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "execute local success",
			in: in{
				runners: formula.Runners{
					formula.Local:  localRunnerMock{},
					formula.Docker: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.Local},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.Local,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "execute docker success",
			in: in{
				runners: formula.Runners{
					formula.Local:  localRunnerMock{},
					formula.Docker: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.Local},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.Docker,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "execute default success",
			in: in{
				runners: formula.Runners{
					formula.Local:  localRunnerMock{},
					formula.Docker: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.Local},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: -1,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "find default runner error",
			in: in{
				runners: formula.Runners{
					formula.Local:  localRunnerMock{},
					formula.Docker: dockerRunnerMock{},
				},
				config: configRunnerMock{findErr: ErrConfigNotFound},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: -1,
					Verbose: false,
				},
			},
			want: ErrConfigNotFound,
		},
		{
			name: "run formula error",
			in: in{
				runners: formula.Runners{
					formula.Local:  localRunnerMock{},
					formula.Docker: dockerRunnerMock{err: errors.New("error to run formula")},
				},
				config: configRunnerMock{runType: formula.Docker},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: -1,
					Verbose: false,
				},
			},
			want: errors.New("error to run formula"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executorManager := NewExecutor(tt.in.runners, tt.in.config)
			got := executorManager.Execute(tt.in.exe)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Execute(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

type localRunnerMock struct {
	err error
}

func (l localRunnerMock) Run(formula.Definition, api.TermInputType, bool) error {
	return l.err
}

type dockerRunnerMock struct {
	err error
}

func (d dockerRunnerMock) Run(formula.Definition, api.TermInputType, bool) error {
	return d.err
}

type configRunnerMock struct {
	runType   formula.RunnerType
	createErr error
	findErr   error
}

func (c configRunnerMock) Create(runType formula.RunnerType) error {
	return c.createErr
}

func (c configRunnerMock) Find() (formula.RunnerType, error) {
	return c.runType, c.findErr
}
