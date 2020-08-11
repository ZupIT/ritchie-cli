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
		file      stream.FileWriteReadExister
		InputList prompt.InputList
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
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "yes", nil
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
				InputList: inputListCustomMock{
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
				InputList: inputListErrorMock{},
			},
		},
		{
			name: "success when metrics file exist",
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("no"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail on read when metrics file exist",
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("no"), errors.New("error reading file")
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail on write when metrics file exist",
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("no"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return errors.New("error writing file")
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsCmd := NewMetricsCmd(tt.in.file, tt.in.InputList)
			if err := metricsCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("metrics command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
