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

package credential

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var fileManager = stream.NewFileManager()
var dirManager = stream.NewDirManager(fileManager)
var homeDir, _ = os.UserHomeDir()
var credSettings = NewSettings(fileManager, dirManager, homeDir)

func providersPath() string {
	tempDir := os.TempDir()
	path := fmt.Sprintf("%s/providers.json", tempDir)
	return path
}

func TestSettings_ReadCredentialsFields(t *testing.T) {
	credentials, err := credSettings.ReadCredentialsFields("../../testdata/credentials.json")
	if err != nil {
		t.Errorf("Error reading credentials fields")
	}

	if credentials == nil || len(credentials) <= 0 {
		t.Errorf("Error reading credentials fields, cannot be empty or null")
	}
}

func TestSettings_ReadCredentialsValue(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "run with success",
			path:    "../../testdata/.rit/credentials/",
			wantErr: false,
		},
		{
			name:    "error on json unmarshal",
			path:    "../../testdata/.rit/credentialserr/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, err := credSettings.ReadCredentialsValue(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read credentials value error = %s, wantErr %v", err, tt.wantErr)
			}

			if (credentials == nil || len(credentials) <= 0) != tt.wantErr {
				t.Errorf("Error reading credentials, cannot be empty or null %v", len(credentials))
			}
		})
	}
}

func TestSettings_ReadCredentialsValueInEnv(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		env     string
		wantErr bool
	}{
		{
			name:    "run with success",
			path:    "../../testdata/.rit/credentials/",
			env:     "default",
			wantErr: false,
		},
		{
			name:    "error on json unmarshal",
			path:    "../../testdata/.rit/credentialserr/",
			env:     "error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, err := credSettings.ReadCredentialsValueInEnv(tt.path, tt.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read credentials value in env error = %s, wantErr %v", err, tt.wantErr)
			}

			if (credentials == nil || len(credentials) <= 0) != tt.wantErr {
				t.Errorf("Error reading credentials, cannot be empty or null")
			}
		})
	}
}

func TestSettings_WriteCredentialsFields(t *testing.T) {
	defer os.Remove(providersPath())
	var tests = []struct {
		name    string
		path    string
		fields  Fields
		wantErr bool
	}{
		{
			name:    "Run with success",
			path:    providersPath(),
			fields:  NewDefaultCredentials(),
			wantErr: false,
		},
		{
			name:    "Error with invalid path",
			path:    "",
			fields:  NewDefaultCredentials(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := credSettings.WriteCredentialsFields(tt.fields, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write credentials fields error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestSettings_WriteDefaultCredentialsFields(t *testing.T) {
	err := credSettings.WriteDefaultCredentialsFields(providersPath())
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error writing credentials: %s", err)
	}
}

func TestNewDefaultCredentials(t *testing.T) {
	defaultCredentials := NewDefaultCredentials()

	if defaultCredentials == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(defaultCredentials) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}
}

func TestSingleSettings_WriteDefaultCredentialsOnExistingFile(t *testing.T) {
	credentials := Fields{
		"customField": []Field{},
	}
	fieldsData, err := json.Marshal(credentials)
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}

	// Write an initial credential file
	err = ioutil.WriteFile(providersPath(), fieldsData, os.ModePerm)
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}

	// Call the method
	err = credSettings.WriteDefaultCredentialsFields(providersPath())
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}

	// Reopen file and check if previous config was not lost
	file, _ := ioutil.ReadFile(providersPath())
	var fields Fields
	err = json.Unmarshal(file, &fields)
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}
	if len(fields) != len(NewDefaultCredentials())+1 {
		t.Errorf("Writing existing credentials did not succeed in adding a field")
	}
	if fields["customField"] == nil {
		t.Errorf("Writing existing credentials did not save custom field")
	}
}

func TestProviderPath(t *testing.T) {
	provider := credSettings.ProviderPath()
	slicedPath := strings.Split(provider, string(os.PathSeparator))
	providersJson := slicedPath[len(slicedPath)-1]

	if providersJson != "providers.json" {
		t.Errorf("Providers path must end on providers.json")
	}
}

func TestCredentialsPath(t *testing.T) {
	credentials := credSettings.CredentialsPath()
	slicedPath := strings.Split(credentials, string(os.PathSeparator))
	fmt.Println(slicedPath)
	providersDir := slicedPath[len(slicedPath)-1]

	if providersDir != "credentials" {
		t.Errorf("Providers path must end on credentials dir")
	}
}

func TestProvidersArr(t *testing.T) {
	credentials := NewDefaultCredentials()
	providersArray := NewProviderArr(credentials)

	if providersArray[len(providersArray)-1] != AddNew {
		t.Errorf("%q option must be the last one", AddNew)
	}

	if providersArray == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(providersArray) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}

}
