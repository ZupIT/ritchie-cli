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
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

func Test_initCmd_runAnyEntry(t *testing.T) {
	type inBool struct {
		label string
		value bool
		err   error
	}

	type in struct {
		repoAddErr     error
		gitTag         git.Tag
		gitTagErr      error
		tutorialHolder rtutorial.TutorialHolder
		tutorialErr    error
		fileWriteErr   error
		configRunErr   error
		inBools        []inBool
		inList         string
		inListErr      error
		ritConfigErr   error
	}

	tests := []struct {
		name       string
		in         in
		wantErr    error
		inputStdin string
		inputFlag  []string
	}{
		{
			name: "success to add commons repo and accept to send metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         LocalRunType,
			},
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"local\"}\n",
		},
		{
			name: "success when not add commons and not accept to send metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBools: []inBool{
					{
						label: AddTheCommunityRepo,
						value: false,
					},
					{
						label: AgreeSendMetrics,
						value: false,
					},
				},
				inList: DockerRunType,
			},
			inputStdin: "{\"addCommons\": false,\"sendMetrics\": false, \"runType\": \"docker\" }\n",
		},
		{
			name: "warning when call git.LatestTag",
			in: in{
				gitTagErr:      errors.New("error to get latest tag"),
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         DockerRunType,
			},
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"docker\"}\n",
		},
		{
			name: "warning when call repo.Add",
			in: in{
				repoAddErr:     errors.New("error to add commons repo"),
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         LocalRunType,
			},
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
			inputFlag: []string{"--sendMetrics=yes", "--addCommons=yes", "--runType=local"},
		},
		{
			name: "find tutorial error",
			in: in{
				gitTag:      git.Tag{Name: "1.0.0"},
				tutorialErr: errors.New("error to find tutorial"),
				inList:      LocalRunType,
			},
			wantErr:    errors.New("error to find tutorial"),
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error in select response of metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBools: []inBool{
					{
						label: AgreeSendMetrics,
						value: false,
						err:   errors.New("error to select metrics response"),
					},
				},
				inList: LocalRunType,
			},
			wantErr: errors.New("error to select metrics response"),
		},
		{
			name: "error in select response of run tyoe",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inListErr:      errors.New("error to select run type"),
			},
			wantErr: errors.New("error to select run type"),
		},
		{
			name: "error to write metric file",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				fileWriteErr:   errors.New("error to write metric file"),
				inList:         LocalRunType,
			},
			wantErr:    errors.New("error to write metric file"),
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error to create runner config",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				configRunErr:   errors.New("error to create runner config"),
				inList:         LocalRunType,
			},
			wantErr:    errors.New("error to create runner config"),
			inputStdin: "{\"addCommons\": true, \"sendMetrics\": false, \"runType\": \"local\"}\n",
		},
		{
			name: "error to select response of commons repo",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inBools: []inBool{
					{
						label: AddTheCommunityRepo,
						value: false,
						err:   errors.New("error to select commons repo response"),
					},
				},
				inList: LocalRunType,
			},
			wantErr: errors.New("error to select commons repo response"),
		},
		{
			name: "error stdin invalid formula run type",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         "invalid",
			},
			wantErr:    ErrInvalidRunType,
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"invalid\"}\n",
		},
		{
			name: "error to write ritchie configs",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         LocalRunType,
				ritConfigErr:   errors.New("error to write ritchie configs"),
			},
			wantErr:    errors.New("error to write ritchie configs"),
			inputStdin: "{\"addCommons\": true,\"sendMetrics\": true, \"runType\": \"local\"}\n",
		},
		{
			name: "success with flags, add commons and send metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         LocalRunType,
			},
			inputFlag: []string{"--sendMetrics=yes", "--addCommons=yes", "--runType=local"},
		},
		{
			name: "success with flags, when not add commons and not send metrics",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				inList:         LocalRunType,
			},
			inputFlag: []string{"--sendMetrics=no", "--addCommons=no", "--runType=docker"},
		},
		{
			name: "error with flags, invalid metrics value",
			in: in{
				fileWriteErr:   errors.New("please provide a valid value to the flag 'sendmetrics'"),
				inList:         LocalRunType,

			},
			inputFlag: []string{"--sendMetrics=invalidValue", "--addCommons=no", "--runType=local"},
			wantErr: errors.New("please provide a valid value to the flag 'sendmetrics'"),
		},
		{
			name: "error with flags, invalid commons value",
			in: in{
				fileWriteErr: errors.New("please provide a valid value to the flag 'addCommons'"),
			},
			inputFlag: []string{"--sendMetrics=no", "--addCommons=invalidValue", "--runType=local"},
			wantErr: errors.New("please provide a valid value to the flag 'addCommons'"),
		},
		{
			name: "error with flags, error to write metric file",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				fileWriteErr:   errors.New("error to write metric file"),
			},
			inputFlag: []string{"--sendMetrics=yes", "--addCommons=no", "--runType=local"},
			wantErr: errors.New("error to write metric file"),
		},
		{
			name: "error with flags, error to write metric file",
			in: in{
				gitTag:         git.Tag{Name: "1.0.0"},
				tutorialHolder: rtutorial.TutorialHolder{Current: rtutorial.DefaultTutorial},
				fileWriteErr:   errors.New("error to write metric file"),
			},
			inputFlag: []string{"--sendMetrics=yes", "--addCommons=no", "--runType=local"},
			wantErr: errors.New("error to write metric file"),
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
			for _, in := range in.inBools {
				inBoolMock.On("Bool", in.label, mock.Anything, mock.Anything).Return(in.value, in.err)
			}
			inBoolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

			inListMock := new(mocks.InputListMock)
			inListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(in.inList, in.inListErr)

			ritConfigMock := new(mocks.RitConfigMock)
			ritConfigMock.On("Write", mock.Anything).Return(in.ritConfigErr)
			metricSender := metric.NewHttpSender("", http.DefaultClient)

			initPrompt := NewInitCmd(
				repoMock,
				gitMock,
				tutorialMock,
				configRunnerMock,
				fileMock,
				inListMock,
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
				inListMock,
				inBoolMock,
				metricSender,
				ritConfigMock,
			)

			initFlag := NewInitCmd(
				repoMock,
				gitMock,
				tutorialMock,
				configRunnerMock,
				fileMock,
				inListMock,
				inBoolMock,
				metricSender,
				ritConfigMock,
			)

			initPrompt.PersistentFlags().Bool("stdin", false, "input by stdin")
			initStdin.PersistentFlags().Bool("stdin", true, "input by stdin")
			initFlag.PersistentFlags().Bool("stdin", false, "input by stdin")
			initFlag.SetArgs(tt.inputFlag)

			newReader := strings.NewReader(tt.inputStdin)
			initStdin.SetIn(newReader)

			got := initPrompt.Execute()
			assert.Equal(t, tt.wantErr, got)

			if tt.inputStdin != "" {
				gotStdin := initStdin.Execute()
				assert.Equal(t, tt.wantErr, gotStdin)
			} else if len(tt.inputFlag) != 0 {
				gotFlag := initFlag.Execute()
				assert.Equal(t, tt.wantErr, gotFlag)
			}
		})
	}
}
