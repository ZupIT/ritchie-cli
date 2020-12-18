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

	"github.com/stretchr/testify/assert"

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

func TestReadCredentialsFields(t *testing.T) {
	credentials, err := credSettings.ReadCredentialsFields("./testdata/credentials.json")

	assert.NoError(t, err)
	assert.NotNil(t, credentials)
	assert.Greater(t, len(credentials), 0)
}

func TestReadCredentialsValues(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "run with success",
			path:    "./testdata/.rit/credentials/",
			wantErr: "",
		},
		{
			name:    "error on json unmarshal",
			path:    "./testdata/.rit/credentialserr/",
			wantErr: "invalid character 'e' looking for beginning of object key string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, err := credSettings.ReadCredentialsValue(tt.path)
			if err == nil {
				assert.Empty(t, tt.wantErr)
				assert.Greater(t, len(credentials), 0)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestReadCredentialsValueInEnv(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		env     string
		wantErr string
	}{
		{
			name:    "run with success",
			path:    "./testdata/.rit/credentials/",
			env:     "default",
			wantErr: "",
		},
		{
			name:    "error on json unmarshal",
			path:    "./testdata/.rit/credentialserr/",
			env:     "error",
			wantErr: "invalid character 'e' looking for beginning of object key string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials, err := credSettings.ReadCredentialsValueInEnv(tt.path, tt.env)
			if err == nil {
				assert.Empty(t, tt.wantErr)
				assert.Greater(t, len(credentials), 0)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestWriteCredentialsFields(t *testing.T) {
	providersPath := providersPath()
	defaultCred := NewDefaultCredentials()
	defer os.Remove(providersPath)
	var tests = []struct {
		name    string
		path    string
		fields  Fields
		wantErr string
	}{
		{
			name:    "Run with success",
			path:    providersPath,
			fields:  defaultCred,
			wantErr: "",
		},
		{
			name:    "Error with invalid path",
			path:    "",
			fields:  defaultCred,
			wantErr: "open : no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Remove(providersPath)
			err := credSettings.WriteCredentialsFields(tt.fields, tt.path)
			if err == nil {
				assert.Empty(t, tt.wantErr)
				file, err := ioutil.ReadFile(providersPath)
				assert.NoError(t, err)

				fields := Fields{}
				err = json.Unmarshal(file, &fields)
				assert.NoError(t, err)

				assert.Equal(t, defaultCred, fields)
			} else {
				assert.EqualError(t, err, tt.wantErr)
				assert.NoFileExists(t, providersPath)
			}
		})
	}
}

func TestWriteDefaultCredentialsFields(t *testing.T) {
	path := providersPath()
	defaultCred := NewDefaultCredentials()
	defer os.Remove(path)
	err := credSettings.WriteDefaultCredentialsFields(providersPath())
	assert.NoError(t, err)
	file, err := ioutil.ReadFile(path)
	assert.NoError(t, err)

	fields := Fields{}
	err = json.Unmarshal(file, &fields)
	assert.NoError(t, err)

	assert.Equal(t, defaultCred, fields)
}

func TestNewDefaultCredentials(t *testing.T) {
	defaultCredentials := NewDefaultCredentials()

	assert.NotNil(t, defaultCredentials)
	assert.Greater(t, len(defaultCredentials), 0)
}

func TestWriteDefaultCredentialsOnExistingFile(t *testing.T) {
	credentials := Fields{
		"customField": []Field{},
	}
	fieldsData, err := json.Marshal(credentials)
	assert.NoError(t, err)

	// Write an initial credential file
	path := providersPath()
	err = ioutil.WriteFile(path, fieldsData, os.ModePerm)
	defer os.Remove(path)
	assert.NoError(t, err)

	// Call the method
	err = credSettings.WriteDefaultCredentialsFields(path)
	assert.NoError(t, err)

	// Reopen file and check if previous config was not lost
	file, _ := ioutil.ReadFile(path)
	var fields Fields
	err = json.Unmarshal(file, &fields)
	assert.NoError(t, err)
	assert.Equal(t, len(NewDefaultCredentials())+1, len(fields), "Writing existing credentials did not succeed in adding a field")
	assert.NotNil(t, fields["customField"], "Writing existing credentials did not save custom field")
}

func TestProviderPath(t *testing.T) {
	provider := credSettings.ProviderPath()
	slicedPath := strings.Split(provider, string(os.PathSeparator))
	providersJson := slicedPath[len(slicedPath)-1]

	assert.Equal(t, "providers.json", providersJson)
}

func TestCredentialsPath(t *testing.T) {
	credentials := credSettings.CredentialsPath()
	slicedPath := strings.Split(credentials, string(os.PathSeparator))
	fmt.Println(slicedPath)
	providersDir := slicedPath[len(slicedPath)-1]

	assert.Equal(t, "credentials", providersDir)
}

func TestProvidersArr(t *testing.T) {
	credentials := NewDefaultCredentials()
	providersArray := NewProviderArr(credentials)

	assert.Equal(t, providersArray[len(providersArray)-1], AddNew)
	assert.NotNil(t, providersArray)
	assert.Equal(t, len(providersArray), 0)
}
