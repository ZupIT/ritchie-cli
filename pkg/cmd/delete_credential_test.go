package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type credDeleteMock struct {
	deleteMock func() error
}

func (c credDeleteMock) Delete(s string) error {
	return c.deleteMock()
}

type fieldsTestDeleteCredentialCmd struct {
	credDelete credential.CredDelete
	reader     credential.ReaderWriterPather
	envFinder  env.Finder
	inputBool  prompt.InputBool
	inputList  prompt.InputList
}

func TestDeleteCredential(t *testing.T) {
	stdinTest := &deleteCredential{
		Provider: "github",
	}

	deleteSuccess := credDeleteMock{
		deleteMock: func() error {
			return nil
		},
	}

	tests := []struct {
		name       string
		wantErr    bool
		fields     fieldsTestDeleteCredentialCmd
		inputStdin string
	}{
		{
			name:    "execute with success",
			wantErr: false,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: deleteSuccess,
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Env: "default", Credential: "{}"}}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
		{
			name:    "error on find env",
			wantErr: true,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader:     credSettingsMock{},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{}, errors.New("some error on find env")
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
		{
			name:    "error to read credentials value",
			wantErr: true,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{}, errors.New("ReadCredentialsValue error")
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
		{
			name:    "error when there are no credentials in the env",
			wantErr: false,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: "",
		},
		{
			name:    "error on input list",
			wantErr: true,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Env: "default", Credential: "{}"}}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListErrorMock{},
			},
			inputStdin: "",
		},
		{
			name:    "error on input bool",
			wantErr: true,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Env: "default", Credential: "{}"}}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputBoolErrorMock{},
				inputList: inputListMock{},
			},
			inputStdin: "",
		},
		{
			name:    "cancel when input bool is false",
			wantErr: false,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Env: "default", Credential: "{}"}}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputFalseMock{},
				inputList: inputListMock{},
			},
			inputStdin: "",
		},
		{
			name:    "error on Delete",
			wantErr: true,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{
					deleteMock: func() error {
						return errors.New("some error on Delete")
					},
				},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Env: "default", Credential: "{}"}}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
		{
			name:    "error different provider",
			wantErr: false,
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: deleteSuccess,
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInEnvMock: func(path string, env string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "gitlab", Env: "default", Credential: "{}"}}, nil
					},
				},
				envFinder: envFinderCustomMock{
					find: func() (env.Holder, error) {
						return env.Holder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteCredentialCmd := NewDeleteCredentialCmd(tt.fields.credDelete, tt.fields.reader, tt.fields.envFinder, tt.fields.inputBool, tt.fields.inputList)
			deleteCredentialStdin := NewDeleteCredentialCmd(tt.fields.credDelete, tt.fields.reader, tt.fields.envFinder, tt.fields.inputBool, tt.fields.inputList)

			deleteCredentialCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			deleteCredentialStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

			newReader := strings.NewReader(tt.inputStdin)
			deleteCredentialStdin.SetIn(newReader)

			if err := deleteCredentialCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("delete credential command error = %v, wantErr %v", err, tt.wantErr)
			}

			itsTestCaseWithStdin := tt.inputStdin != ""
			if err := deleteCredentialStdin.Execute(); (err != nil) != tt.wantErr && itsTestCaseWithStdin {
				t.Errorf("delete credential stdin command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
