package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_RunCmd(t *testing.T) {
	var runCmdCases = []struct {
		fileMock     stream.FileWriteReadExister
		inputMock    prompt.InputList
		name         string
		choose       string
		want         error
		inputFileErr error
		inputFlag    []string
	}{
		{
			name:         "success with prompt",
			choose:       "yes",
			inputFileErr: nil,
			want:         nil,
		},
		{
			name:         "failed  with prompt",
			choose:       "yes",
			inputFileErr: errors.New("failed to write file"),
			want:         errors.New("failed to write file"),
		},
		{
			name:         "success with flags",
			inputFlag:    []string{"--metrics=no"},
			inputFileErr: nil,
			want:         nil,
		},
		{
			name:         "fail with flags",
			inputFlag:    []string{"--metrics=no"},
			inputFileErr: errors.New("failed to write file"),
			want:         errors.New("failed to write file"),
		},
	}

	for _, tt := range runCmdCases {
		t.Run(tt.name, func(t *testing.T) {
			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.choose, tt.want)

			inputFileMock := new(mocks.FileManagerMock)
			inputFileMock.On("Write", mock.Anything, mock.Anything).Return(tt.inputFileErr)

			metricsCmd := NewMetricsCmd(inputFileMock, inputListMock)
			metricsCmd.ParseFlags(tt.inputFlag)

			got := metricsCmd.Execute()
			assert.Equal(t, tt.want, got)

		})
	}
}

func Test_ResolveInput(t *testing.T) {

	var tests = []struct {
		name      string
		choose    string
		fileMock  stream.FileWriteReadExister
		want      error
		inputFlag []string
	}{
		{
			name:   "test run prompt",
			choose: "yes",
			want:   nil,
		},
		{
			name:      "test run flag",
			inputFlag: []string{"--metrics=no"},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.choose, tt.want)

			newMetrics := metricsCmd{
				file:  tt.fileMock,
				input: inputListMock,
			}

			metricsCmd := NewMetricsCmd(tt.fileMock, inputListMock)
			metricsCmd.ParseFlags(tt.inputFlag)

			got, err := newMetrics.resolveInput(metricsCmd)
			fmt.Println(got, err)
		})
	}
}

func Test_RunPrompt(t *testing.T) {
	var runPromptCases = []struct {
		fileMock  stream.FileWriteReadExister
		inputMock prompt.InputList
		name      string
		choose    string
		want      error
	}{
		{
			name:   "success with yes",
			choose: "yes",
			want:   nil,
		},
		{
			name:   "fail on choose with yes",
			choose: "yes",
			want:   errors.New("error on choose"),
		},
	}
	for _, tt := range runPromptCases {
		t.Run(tt.name, func(t *testing.T) {
			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.choose, tt.want)

			newMetrics := metricsCmd{
				file:  tt.fileMock,
				input: inputListMock,
			}

			_, err := newMetrics.runPrompt()

			assert.Equal(t, tt.want, err)

		})
	}

}

func Test_RunFlag(t *testing.T) {
	var runFlagCases = []struct {
		fileMock  stream.FileWriteReadExister
		inputMock prompt.InputList
		choose    string
		name      string
		want      error
	}{
		{
			name:   "success with yes",
			choose: "yes",
			want:   nil,
		},
		{
			name:   "error invalid flag value",
			choose: "invalid",
			want:   errors.New("please provide a valid value to the flag metrics"),
		},
	}

	for _, tt := range runFlagCases {
		t.Run(tt.name, func(t *testing.T) {
			newMetrics := metricsCmd{
				file:  tt.fileMock,
				input: tt.inputMock,
			}
			cobraCmd := NewMetricsCmd(tt.fileMock, tt.inputMock)
			cobraCmd.Flags().Set(metricsFlagName, tt.choose)
			got, err := newMetrics.runFlag(cobraCmd)
			fmt.Println(got, err)

		})
	}
}
