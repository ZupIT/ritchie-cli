package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
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
	ctxFinder  rcontext.Finder
	inputBool  prompt.InputBool
	inputList  prompt.InputList
}

const provider = "github"

// TODO: remove upon stdin deprecation
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Context: "default", Credential: "{}"}}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
		{
			name:    "error on find context",
			wantErr: "some error on find Context",
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader:     credSettingsMock{},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{}, errors.New("some error on find Context")
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{}, errors.New("ReadCredentialsValue error")
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
					},
				},
				inputBool: inputTrueMock{},
				inputList: inputListMock{},
			},
			inputStdin: createJSONEntry(stdinTest),
		},
		{
			name:    "error when there are no credentials in the context",
			wantErr: "",
			fields: fieldsTestDeleteCredentialCmd{
				credDelete: credDeleteMock{},
				reader: credSettingsCustomMock{
					CredentialsPathMock: func() string {
						return ""
					},
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Context: "default", Credential: "{}"}}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Context: "default", Credential: "{}"}}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Context: "default", Credential: "{}"}}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "github", Context: "default", Credential: "{}"}}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
					ReadCredentialsValueInContextMock: func(path string, context string) ([]credential.ListCredData, error) {
						return []credential.ListCredData{{Provider: "gitlab", Context: "default", Credential: "{}"}}, nil
					},
				},
				ctxFinder: ctxFinderCustomMock{
					findMock: func() (rcontext.ContextHolder, error) {
						return rcontext.ContextHolder{Current: ""}, nil
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
			deleteCredentialCmd := NewDeleteCredentialCmd(tt.fields.credDelete, tt.fields.reader, tt.fields.ctxFinder, tt.fields.inputBool, tt.fields.inputList)
			deleteCredentialStdin := NewDeleteCredentialCmd(tt.fields.credDelete, tt.fields.reader, tt.fields.ctxFinder, tt.fields.inputBool, tt.fields.inputList)

			deleteCredentialCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			deleteCredentialStdin.PersistentFlags().Bool("stdin", true, "input by stdin")

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
	credentialPath := filepath.Join(ritHomeDir, "credentials", rcontext.DefaultCtx)
	credentialFile := filepath.Join(credentialPath, provider)
	_ = os.MkdirAll(credentialPath, os.ModePerm)
	defer os.RemoveAll(ritHomeDir)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	ctxFinder := rcontext.NewFinder(ritHomeDir, fileManager)
	credDeleter := credential.NewCredDelete(ritHomeDir, ctxFinder, fileManager)
	credSettings := credential.NewSettings(fileManager, dirManager, homeDir)

	listMock := &InputListMock{}
	listMock.On("List", mock.Anything).Return(provider, nil)

	boolMock := &InputBoolMock{}
	boolMock.On("Bool", mock.Anything).Return(true, nil)

	tests := []struct {
		name      string
		inputBool prompt.InputBool
		inputList prompt.InputList
		wantErr   string
	}{
		{
			name:      "execute with success",
			wantErr:   "",
			inputList: listMock,
			inputBool: boolMock,
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
			jsonData, err := json.Marshal(cred)
			err = ioutil.WriteFile(credentialFile, jsonData, os.ModePerm)
			require.NoError(t, err)

			cmd := NewDeleteCredentialCmd(credDeleter, credSettings, ctxFinder, tt.inputBool, tt.inputList)
			// TODO: remove stdin flag after deprecation
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs([]string{})

			err = cmd.Execute()
			if err != nil {
				require.Equal(t, err.Error(), tt.wantErr)
			} else {
				require.Empty(t, tt.wantErr)
				require.NoFileExists(t, credentialFile)
			}
		})
	}
}

type InputListMock struct {
	mock.Mock
}

func (m *InputListMock) List(string, []string, ...string) (string, error) {
	args := m.Called()
	return args.String(0), nil
}

type InputBoolMock struct {
	mock.Mock
}

func (m *InputBoolMock) Bool(string, []string, ...string) (bool, error) {
	args := m.Called()
	return args.Bool(0), nil
}
