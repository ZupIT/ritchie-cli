package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
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

const provider = "github"

// TODO: remove upon stdin deprecation, reduce dependencies
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
		wantErr    string
		fields     fieldsTestDeleteCredentialCmd
		inputStdin string
	}{
		{
			name:    "execute with success",
			wantErr: "",
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
			wantErr: "some error on find env",
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
			wantErr: "ReadCredentialsValue error",
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
			wantErr: "you have no defined credentials in this env",
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
			wantErr: "some error",
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
			wantErr: "error on boolean list",
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
			wantErr: "",
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
			wantErr: "some error on Delete",
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
			wantErr: "",
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
			deleteCredentialCmd.SetArgs([]string{})
			deleteCredentialStdin.SetArgs([]string{})

			newReader := strings.NewReader(tt.inputStdin)
			deleteCredentialStdin.SetIn(newReader)

			err := deleteCredentialCmd.Execute()
			if err != nil {
				require.Equal(t, err.Error(), tt.wantErr)
			} else {
				require.Empty(t, tt.wantErr)
			}

			itsTestCaseWithStdin := tt.inputStdin != ""
			err = deleteCredentialStdin.Execute()
			if itsTestCaseWithStdin {
				if err != nil {
					require.Equal(t, err.Error(), tt.wantErr)
				} else {
					require.Empty(t, tt.wantErr)
				}
			}
		})
	}
}

func TestDeleteCredentialViaPrompt(t *testing.T) {
	homeDir := os.TempDir()
	ritHomeDir := filepath.Join(homeDir, ".rit")
	credentialPath := filepath.Join(ritHomeDir, "credentials", env.Default)
	credentialFile := filepath.Join(credentialPath, provider)
	_ = os.MkdirAll(credentialPath, os.ModePerm)
	defer os.RemoveAll(ritHomeDir)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	ctxFinder := env.NewFinder(ritHomeDir, fileManager)
	credDeleter := credential.NewCredDelete(ritHomeDir, ctxFinder, fileManager)
	credSettings := credential.NewSettings(fileManager, dirManager, homeDir)

	tests := []struct {
		name            string
		inputBoolResult bool
		inputListError  error
		fileShouldExist bool
		args            string
		wantErr         string
	}{
		{
			name:            "execute prompt with success",
			inputBoolResult: true,
		},
		{
			name: "execute flag with success",
			args: "--provider=github",
		},
		{
			name:            "fail on input list error",
			wantErr:         "some error",
			inputListError:  errors.New("some error"),
			fileShouldExist: true,
		},
		{
			name:            "do nothing on input bool refusal",
			inputBoolResult: false,
			fileShouldExist: true,
		},
	}

	cred := credential.Detail{
		Username: "",
		Credential: credential.Credential{
			"username": "my user",
		},
		Service: provider,
		Type:    "text",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(cred)
			err := ioutil.WriteFile(credentialFile, jsonData, os.ModePerm)
			assert.NoError(t, err)

			listMock := &mocks.InputListMock{}
			listMock.On("List", mock.Anything).Return(provider, tt.inputListError)
			boolMock := &mocks.InputBoolMock{}
			boolMock.On("Bool", mock.Anything).Return(tt.inputBoolResult, nil)

			cmd := NewDeleteCredentialCmd(credDeleter, credSettings, ctxFinder, boolMock, listMock)
			// TODO: remove stdin flag after  deprecation
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs([]string{tt.args})

			err = cmd.Execute()
			if err != nil {
				assert.Equal(t, err.Error(), tt.wantErr)
			} else {
				assert.Empty(t, tt.wantErr)
			}

			if tt.fileShouldExist {
				assert.FileExists(t, credentialFile)
			} else {
				assert.NoFileExists(t, credentialFile)
			}
		})
	}
}
