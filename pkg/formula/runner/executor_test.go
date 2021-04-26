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

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestExecute(t *testing.T) {
	type in struct {
		runners            formula.Runners
		config             formula.ConfigRunner
		exe                formula.ExecuteData
		preRunBuilderError error
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
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.LocalRun},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.LocalRun,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "execute docker success",
			in: in{
				runners: formula.Runners{
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.LocalRun},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.DockerRun,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "execute default success",
			in: in{
				runners: formula.Runners{
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.LocalRun},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.DefaultRun,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "execute repo local success",
			in: in{
				runners: formula.Runners{
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.LocalRun},
				exe: formula.ExecuteData{
					Def:     formula.Definition{RepoName: "local-user"},
					InType:  0,
					RunType: formula.DefaultRun,
					Verbose: false,
				},
			},
			want: nil,
		},
		{
			name: "find default runner error",
			in: in{
				runners: formula.Runners{
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{},
				},
				config: configRunnerMock{findErr: ErrConfigNotFound},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.DefaultRun,
					Verbose: false,
				},
			},
			want: ErrConfigNotFound,
		},
		{
			name: "run formula error",
			in: in{
				runners: formula.Runners{
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{err: errors.New("error to run formula")},
				},
				config: configRunnerMock{runType: formula.DockerRun},
				exe: formula.ExecuteData{
					Def:     formula.Definition{},
					InType:  0,
					RunType: formula.DefaultRun,
					Verbose: false,
				},
			},
			want: errors.New("error to run formula"),
		},
		{
			name: "run pre run builder error",
			in: in{
				runners: formula.Runners{
					formula.LocalRun:  localRunnerMock{},
					formula.DockerRun: dockerRunnerMock{},
				},
				config: configRunnerMock{runType: formula.DockerRun},
				exe: formula.ExecuteData{
					Def:     formula.Definition{RepoName: "local-teste"},
					InType:  0,
					RunType: formula.DefaultRun,
					Verbose: false,
				},
				preRunBuilderError: errors.New("error to pre run builder formula"),
			},
			want: errors.New("error to pre run builder formula"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preRunBuilder := new(mocks.PreRunBuilder)
			preRunBuilder.On("Build", mock.Anything).Return(tt.in.preRunBuilderError)

			executorManager := NewExecutor(tt.in.runners, preRunBuilder, tt.in.config)
			got := executorManager.Execute(tt.in.exe)

			assert.Equal(t, tt.want, got)
		})
	}
}

type localRunnerMock struct {
	err error
}

func (l localRunnerMock) Run(def formula.Definition, inputType api.TermInputType, verbose bool, flags *pflag.FlagSet) error {
	return l.err
}

type dockerRunnerMock struct {
	err error
}

func (d dockerRunnerMock) Run(def formula.Definition, inputType api.TermInputType, verbose bool, flags *pflag.FlagSet) error {
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
