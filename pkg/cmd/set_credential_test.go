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
					text: func(name string, required bool, defaultValue string) (string, error) {
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
					text: func(name string, required bool, defaultValue string) (string, error) {
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
					text: func(name string, required bool, defaultValue string) (string, error) {
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
					text: func(name string, required bool, defaultValue string) (string, error) {
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
