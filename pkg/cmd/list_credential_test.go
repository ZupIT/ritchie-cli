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
)

// func Test_ListCredentialCmd(t *testing.T) {
// 	fileManager := stream.NewFileManager()
// 	dirManager := stream.NewDirManager(fileManager)
// 	homeDir, _ := os.UserHomeDir()
// 	credSettings := credential.NewSettings(fileManager, dirManager, homeDir)
//
// 	t.Run("Success case", func(t *testing.T) {
// 		o := NewListCredentialCmd(credSettings)
// 		if err := o.Execute(); err != nil {
// 			t.Errorf("Test_ListCredentialCmd error = %s", err)
// 		}
// 	})
//
// }

//  todo change
func Test_ListCredentialCmCHANGE(t *testing.T) {
	type in struct {
		credFile credential.ReaderWriterPather
	}
	var tests = []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "success run",
			in: in{credFile: credSettingsCustomMock{
				ReadCredentialsValueMock: func(path string) ([]credential.ListCredData, error) {
					credData := credential.ListCredData{
						Provider:   "",
						Credential: "",
						Context:    "",
					}
					credDataArr := []credential.ListCredData{credData}
					return credDataArr, nil
				},
			}},
			wantErr: false,
		},
		{
			name: "success run with no credentials",
			in: in{credFile: credSettingsCustomMock{
				ReadCredentialsValueMock: func(path string) ([]credential.ListCredData, error) {
					return []credential.ListCredData{}, nil
				},
			}},
			wantErr: false,
		},
		{
			name: "fail on read credentials",
			in: in{credFile: credSettingsCustomMock{
				ReadCredentialsValueMock: func(path string) ([]credential.ListCredData, error) {
					return []credential.ListCredData{}, errors.New("errors reading credentials")
				},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewListCredentialCmd(tt.in.credFile)
			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("list credential command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
