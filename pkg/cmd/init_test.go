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
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

func Test_initCmd_runAnyEntry(t *testing.T) {
	type in struct {
		repoAddErr     error
		gitTag         git.Tag
		gitTagErr      error
		tutorialHolder rtutorial.TutorialHolder
		tutorialErr    error
		fileWriteErr   error
		configRunErr   error
		inBool         bool
		inBoolErr      error
		inList         prompt.InputList
		ritConfigErr   error
	}

	tests := []struct {
		name       string
		in         in
		wantErr    error
		inputStdin string
	}{
		{
			name: "success to add commons repo and accept to send metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return AcceptOpt, nil
					},
				},
			},
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"local\"}\n",
		},
		{
			name: "success when not add commons and not accept to send metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         false,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return DockerRunType, nil
						}
						return DeclineOpt, nil
					},
				},
			},
			inputStdin: "{\"addCommons\": false,\"sendMetrics\": false, \"runType\": \"docker\" }\n",
		},
		{
			name: "warning when call git.LatestTag",
			in: in{
				gitTagErr:      errors.New("error to get latest tag"),
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return DockerRunType, nil
						}
						return DeclineOpt, nil
					},
				},
			},
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"docker\"}\n",
		},
		{
			name: "warning when call repo.Add",
			in: in{
				repoAddErr:     errors.New("error to add commons repo"),
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return AcceptOpt, nil
					},
				},
			},
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "find tutorial error",
			in: in{
				gitTag:      git.Tag{Name: "1.0.0"},
				tutorialErr: errors.New("error to find tutorial"),
				inBool:      true,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return AcceptOpt, nil
					},
				},
			},
			wantErr:    errors.New("error to find tutorial"),
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error in select response of metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBoolErr:      errors.New("error to select metrics response"),
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return "", nil
					},
				},
			},
			wantErr: errors.New("error to select metrics response"),
		},
		{
			name: "error in select response of run tyoe",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return "", errors.New("error to select run type")
						}
						return AcceptOpt, nil
					},
				},
			},
			wantErr: errors.New("error to select run type"),
		},
		{
			name: "error to write metric file",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				fileWriteErr:   errors.New("error to write metric file"),
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return AcceptOpt, nil
					},
				},
			},
			wantErr:    errors.New("error to write metric file"),
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error to create runner config",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				configRunErr:   errors.New("error to create runner config"),
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return AcceptOpt, nil
					},
				},
			},
			wantErr:    errors.New("error to create runner config"),
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error to select response of commons repo",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBoolErr:      errors.New("error to select commons repo response"),
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return LocalRunType, nil
						}
						return AcceptOpt, nil
					},
				},
			},
			wantErr: errors.New("error to select commons repo response"),
		},
		{
			name: "error stdin invalid formula run type",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBool:         true,
				inList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						if name == SelectFormulaTypeQuestion {
							return "invalid", nil
						}
						return AcceptOpt, nil
					},
				},
			},
			wantErr:    ErrInvalidRunType,
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"invalid\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			repoMock := new(mocks.RepoManager)
			repoMock.On("Add", mock.Anything).Return(in.repoAddErr)
			gitMock := new(mocks.GitRepositoryMock)
			gitMock.On("LatestTag", mock.Anything).Return(in.gitTag, in.gitTagErr)
			tutorialMock := new(mocks.TutorialFindSetterMock)
			tutorialMock.On("Find").Return(in.tutorialHolder, in.tutorialErr)
			fileMock := new(mocks.FileManager)
			fileMock.On("Write", mock.Anything, mock.Anything).Return(in.fileWriteErr)
			configRunnerMock := new(mocks.ConfigRunnerMock)
			configRunnerMock.On("Create", mock.Anything).Return(in.configRunErr)
			inBoolMock := new(mocks.InputBoolMock)
			inBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(in.inBool, in.inBoolErr)
			ritConfigMock := new(mocks.RitConfigMock)
			ritConfigMock.On("Write", mock.Anything).Return(in.ritConfigErr)
			metricSender := metric.NewHttpSender("", http.DefaultClient)

			initPrompt := NewInitCmd(
				repoMock,
				gitMock,
				tutorialMock,
				configRunnerMock,
				fileMock,
				in.inList,
				inBoolMock,
				metricSender,
				ritConfigMock,
			)
			initStdin := NewInitCmd(
				repoMock,
				gitMock,
				tutorialMock,
				configRunnerMock,
				fileMock,
				in.inList,
				inBoolMock,
				metricSender,
				ritConfigMock,
			)

			initPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			initStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			initStdin.SetIn(newReader)

			got := initPrompt.Execute()
			assert.Equal(t, tt.wantErr, got)

			if tt.inputStdin != "" {
				gotStdin := initStdin.Execute()
				assert.Equal(t, tt.wantErr, gotStdin)
			}
		})
	}
}
