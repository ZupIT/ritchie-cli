package cmd

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func Test_metricsCmd_runPrompt(t *testing.T) {
	type in struct {
		file  stream.FileWriteReadExister
		input prompt.InputList
	}

	var tests = []struct {
		name    string
		wantErr bool
		in      in
	}{
		{
			name:    "success when metrics file dont exist",
			wantErr: false,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "yes", nil
					},
				},
			},
		},
		{
			name:    "success when not accept send metrics",
			wantErr: false,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "no", nil
					},
				},
			},
		},
		{
			name:    "fail on write file when metrics file dont exist",
			wantErr: true,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return errors.New("reading file error")
					},
				},
				input: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "yes", nil
					},
				},
			},
		},
		{
			name:    "fail on input list when metrics file dont exist",
			wantErr: true,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				input: inputListErrorMock{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsCmd := NewMetricsCmd(tt.in.file, tt.in.input)
			if err := metricsCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("metrics command error = %v | error wanted: %v", err, tt.wantErr)
			}
		})
	}
}

type CheckerMock struct {
	CheckMock func() bool
}

func (cm CheckerMock) Check() bool {
	return cm.CheckMock()
}
