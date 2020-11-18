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
	"errors"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func Test_initCmd_runAnyEntry(t *testing.T) {
	someError := errors.New("some error")
	type fields struct {
		repo     formula.RepositoryAdder
		git      git.Repositories
		tutorial rtutorial.Finder
		config   formula.ConfigRunner
		file     stream.FileWriteReadExister
		inList   prompt.InputList
		inBool   prompt.InputBool
	}

	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		inputStdin string
	}{
		{
			name: "success to add commons repo and accept to send metrics",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"local\"}\n",
		},
		{
			name: "success when not add commons and not accept to send metrics",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputFalseMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.DockerRun.String(), nil
						}
						return DoNotAcceptMetrics, nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": false,\"sendMetrics\": false, \"runType\": \"docker\" }\n",
		},
		{
			name: "warning when call git.LatestTag",
			fields: fields{
				repo: defaultRepoAdderMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				git: GitRepositoryMock{
					latestTag: func(info git.RepoInfo) (git.Tag, error) {
						return git.Tag{}, someError
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.DockerRun.String(), nil
						}
						return DoNotAcceptMetrics, nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"docker\"}\n",
		},
		{
			name: "warning when call repo.Add",
			fields: fields{
				repo: repoListerAdderCustomMock{
					add: func(d formula.Repo) error {
						return someError
					},
				},
				git: defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    false,
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "find tutorial error",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				tutorial: TutorialFindSetterCustomMock{
					find: func() (rtutorial.TutorialHolder, error) {
						return rtutorial.TutorialHolder{}, errors.New("not found tutorial")
					},
				},
				inBool: inputTrueMock{},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    true,
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error in select response of metrics",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return "", someError
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error in select response of run tyoe",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return "", someError
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error to write file",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return errors.New("error to write file")
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    true,
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error to create config",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: errors.New("error to create config"),
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    true,
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error to select response of commons repo",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputBoolErrorMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return formula.LocalRun.String(), nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error stdin invalid formula run type",
			fields: fields{
				repo: defaultRepoAdderMock,
				git:  defaultGitRepositoryMock,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(string, []byte) error {
						return nil
					},
				},
				config: ConfigRunnerMock{
					createErr: nil,
				},
				tutorial: TutorialFinderMock{},
				inBool:   inputTrueMock{},
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return "invalid", nil
						}
						return AcceptMetrics, nil
					},
				},
			},
			wantErr:    true,
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"invalid\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields

			metricSender := SenderMock{}

			initPrompt := NewInitCmd(
				field.repo,
				field.git,
				field.tutorial,
				field.config,
				field.file,
				field.inList,
				field.inBool,
				metricSender,
			)
			initStdin := NewInitCmd(
				field.repo,
				field.git,
				field.tutorial,
				field.config,
				field.file,
				field.inList,
				field.inBool,
				metricSender,
			)

			initPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			initStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			initStdin.SetIn(newReader)

			if err := initPrompt.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := initStdin.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("init_runStdin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type SenderMock struct {
	SendUserStateMock   func()
	SendCommandDataMock func()
}

func (s SenderMock) SendUserState(ritVersion string) {
	s.SendUserStateMock()
}

func (s SenderMock) SendCommandData(cmd metric.SendCommandDataParams) {
	s.SendCommandDataMock()
}
