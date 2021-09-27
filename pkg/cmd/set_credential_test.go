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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestSetCredentialCmd(t *testing.T) {
	tmp := os.TempDir()
	home := filepath.Join(tmp, "SetCredential")
	credentialFile := filepath.Join(home, "credentials", env.Default, provider)
	credentialFileCircle := filepath.Join(home, "credentials", env.Default, "circleci")
	os.MkdirAll(filepath.Join(home, ".rit"), os.ModePerm)
	defer os.RemoveAll(home)

	provider := "github"

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	envFinder := env.NewFinder(home, fileManager)
	setter := credential.NewSetter(home, envFinder, dirManager)
	credSettings := credential.NewSettings(fileManager, dirManager, home)

	var tests = []struct {
		name             string
		args             []string
		provider         string
		credentialFile   string
		passError        error
		textError        error
		addMoreCredError error
		fieldNameError   error
		fieldTypeError   error
		err              error
	}{
		{
			name:           "success run with selected provider",
			args:           []string{},
			provider:       provider,
			credentialFile: credentialFile,
		},
		{
			name:           "success run with new provider",
			args:           []string{},
			provider:       credential.AddNew,
			credentialFile: credentialFileCircle,
		},
		{
			name:           "fail defining field name",
			args:           []string{},
			provider:       credential.AddNew,
			fieldNameError: errors.New("fail defining field name"),
			err:            errors.New("fail defining field name"),
		},
		{
			name:           "fail defining field type",
			args:           []string{},
			provider:       credential.AddNew,
			fieldTypeError: errors.New("fail defining field type"),
			err:            errors.New("fail defining field type"),
		},
		{
			name:             "fail adding more credentials",
			args:             []string{},
			provider:         credential.AddNew,
			addMoreCredError: errors.New("fail add more credential"),
			err:              errors.New("fail add more credential"),
		},
		{
			name:      "fail to provide text",
			args:      []string{},
			provider:  provider,
			textError: errors.New("text error"),
			err:       errors.New("text error"),
		},
		{
			name:      "fail to provide password",
			args:      []string{},
			provider:  provider,
			passError: errors.New("pass error"),
			err:       errors.New("pass error"),
		},
		{
			name: "error provider flag empty",
			args: []string{"--provider="},
			err:  errors.New("please provide a value for 'provider'"),
		},
		{
			name: "error fields flag empty",
			args: []string{"--provider=something", "--fields="},
			err:  errors.New("please provide a value for 'fields'"),
		},
		{
			name: "error values flag empty",
			args: []string{"--provider=something", "--fields=field1", "--values="},
			err:  errors.New("please provide a value for 'values'"),
		},
		{
			name: "error unequal length of fields and values flag",
			args: []string{"--provider=something", "--fields=field1,field2", "--values=value1"},
			err:  errors.New("number of fields does not match with number of values"),
		},
		{
			name:           "success flags",
			args:           []string{"--provider=" + provider, "--fields=field1,field2", "--values=value1,value2"},
			credentialFile: credentialFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Remove(credentialFile)

			inputText := &mocks.InputTextMock{}
			inputBool := &mocks.InputBoolMock{}
			inputList := &mocks.InputListMock{}
			inputPassword := &mocks.InputPasswordMock{}
			inputList.On(
				"List", "Select your provider", mock.Anything, mock.Anything,
			).Return(tt.provider, nil)
			inputText.On("Text", "Define your provider name:", true, mock.Anything).Return("circleci", nil)
			inputText.On("Text", "username:", true, mock.Anything).Return("user", tt.textError)
			inputText.On("Text", "email:", true, mock.Anything).Return("my email", nil)
			inputText.On(
				"Text", "Define your field name: (ex.:token, secretAccessKey)", true, mock.Anything,
			).Return("token", tt.fieldNameError)
			inputList.On(
				"List", "Select your field type:", mock.Anything, mock.Anything,
			).Return("secret", tt.fieldTypeError)
			inputBool.On(
				"Bool", "Add more credentials fields to this provider?", []string{"no", "yes"}, mock.Anything,
			).Return(false, tt.addMoreCredError)
			inputPassword.On("Password", "token:", mock.Anything).Return("some pass", tt.passError)

			cmd := NewSetCredentialCmd(
				setter,
				credSettings,
				inputText,
				inputBool,
				inputList,
				inputPassword,
			)

			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err != nil {
				assert.Equal(t, err, tt.err)
			} else {
				assert.Nil(t, tt.err)
				assert.FileExists(t, tt.credentialFile)
			}
		})
	}
}
