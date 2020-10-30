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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var creds = make(map[string][]credential.Field)

func Test_setCredentialCmd_runPrompt(t *testing.T) {
	type in struct {
		Setter        credential.Setter
		credFile      credential.ReaderWriterPather
		file          stream.FileReadExister
		InputText     prompt.InputText
		InputBool     prompt.InputBool
		InputList     prompt.InputList
		InputPassword prompt.InputPassword
	}
	var tests = []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "success run with no data",
			in: in{
				Setter:    credSetterMock{},
				credFile:  credSettingsMock{},
				InputText: inputSecretMock{},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credential.AddNew, nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: false,
		},
		{
			name: "success run with full data",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "plain text",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["file"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
				},
				InputText: inputTextMock{},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "file", nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: false,
		},
		{
			name: "fail text with full data and file input",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "plain text",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["file"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
				},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						return "", errors.New("text error")
					},
				},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "file", nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "fail to read file",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "plain text",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["type"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return nil, errors.New("error reading file")
					},
				},
				InputText: inputTextMock{},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "file", nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "fail empty credential file",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "plain text",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["type"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), nil
					},
				},
				InputText: inputTextMock{},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "file", nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "fail no file to read",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "plain text",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["type"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
				},
				InputText: inputTextMock{},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "file", nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "fail cannot find any credential in path ",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "plain text",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["type"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
				},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						return "", errors.New("text error")
					},
				},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "type", nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "fail when password return err",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						cred := credential.Field{
							Name: "accesskeyid",
							Type: "secret",
						}
						credArr := []credential.Field{}
						credArr = append(credArr, cred)
						creds["type"] = credArr
						return creds, nil
					},
				},
				file: sMocks.FileReadExisterCustomMock{},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						return "./path/to/my/credentialFile", nil
					},
				},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "type", nil
					},
				},
				InputPassword: inputPasswordErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "fail when write credential fields return err",
			in: in{
				Setter: credSetterMock{},
				credFile: credSettingsCustomMock{
					ReadCredentialsFieldsMock: func(path string) (credential.Fields, error) {
						return credential.Fields{}, errors.New("error reading credentials")
					},
				},
				file: sMocks.FileReadExisterCustomMock{},
				InputText: inputTextCustomMock{
					text: func(name string, required bool) (string, error) {
						return "./path/to/my/credentialFile", nil
					},
				},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return "type", nil
					},
				},
				InputPassword: inputPasswordErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "fail when list return err",
			in: in{
				Setter:        credSetterMock{},
				credFile:      credSettingsMock{},
				InputText:     inputSecretMock{},
				InputBool:     inputFalseMock{},
				InputList:     inputListErrorMock{},
				InputPassword: inputPasswordErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "fail when text return err",
			in: in{
				Setter:    credSetterMock{},
				credFile:  credSettingsMock{},
				InputText: inputTextErrorMock{},
				InputBool: inputFalseMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credential.AddNew, nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
		{
			name: "fail when text bool err",
			in: in{
				Setter:    credSetterMock{},
				credFile:  credSettingsMock{},
				InputText: inputTextMock{},
				InputBool: inputBoolErrorMock{},
				InputList: inputListCustomMock{
					list: func(name string, items []string) (string, error) {
						return credential.AddNew, nil
					},
				},
				InputPassword: inputPasswordMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewSetCredentialCmd(
				tt.in.Setter,
				tt.in.credFile,
				tt.in.file,
				tt.in.InputText,
				tt.in.InputBool,
				tt.in.InputList,
				tt.in.InputPassword,
			)
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("set credential command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// --------------- TESTS -------------->

func (spec *SetCredentialSuite) TestSetCredentialWithSucess() {
	spec.TestInText.On("Text", mock.Anything, mock.Anything).Return("username=ritchie", nil)
	spec.TestInBool.On("Bool", mock.Anything).Return(false)
	spec.TestInList.On("List", mock.Anything, mock.Anything).Return(string(credential.AddNew), nil)
	spec.TestInPassword.On("Password", mock.Anything).Return("s3cr3t", nil)
	spec.TestInSetter.On("Set", mock.Anything).Return(nil)

	spec.TestInReader.On("ReadCredentialsFields", mock.Anything).Return(credential.Fields{}, nil)
	spec.TestInReader.On("ReadCredentialsValue", mock.Anything).Return([]credential.ListCredData{}, nil)
	spec.TestInReader.On("WriteDefaultCredentialsFields", mock.Anything).Return(nil)
	spec.TestInReader.On("WriteCredentialsFields", mock.Anything, mock.Anything).Return(nil)
	spec.TestInReader.On("ProviderPath").Return("")
	spec.TestInReader.On("CredentialsPath").Return("")

	cmd := NewSetCredentialCmd(
		spec.TestInSetter,
		spec.TestInReader,
		spec.TestInFile,
		spec.TestInText,
		spec.TestInBool,
		spec.TestInList,
		spec.TestInPassword,
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	spec.Nil(cmd.Execute())
}

// func (spec *SetCredentialSuite) TestSetCredentialWithSucess2() {

// }

// <-------------- TESTS ---------------

// --------------- RUNNER -------------->

func TestSetCredentialSuite(t *testing.T) {
	suite.Run(t, &SetCredentialSuite{})
}

// <-------------- RUNNER ---------------

// --------------- SETUP -------------->

type SetCredentialSuite struct {
	suite.Suite

	TestInList     *inputList
	TestInBool     *inputBool
	TestInText     *inputText
	TestInPassword *inputPassword
	TestInSetter   *setter
	TestInFile     *fileReadExister
	TestInReader   *readerWriterPather
}

func (suite *SetCredentialSuite) SetupTest() {
	suite.TestInList = new(inputList)
	suite.TestInBool = new(inputBool)
	suite.TestInText = new(inputText)
	suite.TestInPassword = new(inputPassword)
	suite.TestInSetter = new(setter)
	suite.TestInFile = new(fileReadExister)
	suite.TestInReader = new(readerWriterPather)
}

// <-------------- SETUP ---------------

// --------------- MOCKS -------------->

// INPUTS ----->

type inputPassword struct {
	mock.Mock
}

type inputText struct {
	mock.Mock
}

type inputBool struct {
	mock.Mock
}

type inputList struct {
	mock.Mock
}

func (i *inputPassword) Password(label string, helper ...string) (string, error) {
	args := i.Called(label, helper)

	return args.String(0), args.Error(0)
}

func (i *inputText) Text(name string, required bool, helper ...string) (string, error) {
	args := i.Called(name, required, helper)

	return args.String(0), args.Error(0)
}

func (l *inputList) List(name string, items []string, helper ...string) (string, error) {
	args := l.Called(name, items)

	return args.String(0), args.Error(0)
}

func (i *inputBool) Bool(name string, items []string, helper ...string) (bool, error) {
	args := i.Called(name, items, helper)

	return args.Bool(0), args.Error(0)
}

// INPUTS <------------

// OUTRAS FUNCOES ----->

type setter struct {
	mock.Mock
}

type readerWriterPather struct {
	mock.Mock
}

type fileReadExister struct {
	mock.Mock
}

func (s *setter) Set(d credential.Detail) error {
	args := s.Called(d)

	return args.Error(0)
}

func (r *readerWriterPather) ReadCredentialsFields(path string) (credential.Fields, error) {
	args := r.Called(path)

	return args.Get(0).(credential.Fields), args.Error(0)
}

func (r *readerWriterPather) ReadCredentialsValue(path string) ([]credential.ListCredData, error) {
	args := r.Called(path)

	return args.Get(0).([]credential.ListCredData), args.Error(0)
}

func (r *readerWriterPather) WriteCredentialsFields(fields credential.Fields, path string) error {
	args := r.Called(fields, path)

	return args.Error(0)
}

func (r *readerWriterPather) WriteDefaultCredentialsFields(path string) error {
	args := r.Called(path)

	return args.Error(0)
}

func (r *readerWriterPather) ProviderPath() string {
	args := r.Called()

	return args.String(0)
}

func (r *readerWriterPather) CredentialsPath() string {
	args := r.Called()

	return args.String(0)
}

func (f *fileReadExister) Exists(path string) bool {
	args := f.Called(path)

	return args.Bool(0)
}

func (f *fileReadExister) Read(path string) ([]byte, error) {
	args := f.Called(path)

	return args.Get(0).([]byte), args.Error(0)
}

// OUTRAS FUNCOES <-----

// <-------------- MOCKS ---------------
