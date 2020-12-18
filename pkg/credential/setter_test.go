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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SetterTestSuite struct {
	suite.Suite

	HomePath string

	envHolderNil     *env.Holder
	envHolderDefault *env.Holder
	envHolderProd    *env.Holder

	DetailCredentialInfo *Detail
}

func (suite *SetterTestSuite) SetupSuite() {
	nameSuite := "SetterTestSuite"
	tempDir := os.TempDir()
	detailExample := &Detail{
		Service:    "github",
		Username:   "ritchie",
		Credential: Credential{"token": "123", "username": "hackerman"},
		Type:       "",
	}

	suite.HomePath = filepath.Join(tempDir, nameSuite)
	suite.envHolderNil = &env.Holder{Current: ""}
	suite.envHolderDefault = &env.Holder{Current: "default"}
	suite.envHolderProd = &env.Holder{Current: "prod", All: []string{"defauld", "prod"}}
	suite.DetailCredentialInfo = detailExample
}

func (suite *SetterTestSuite) fileInfo(path string) (string, error) {
	fileManager := stream.NewFileManager()
	b, err := fileManager.Read(path)
	return string(b), err
}

func (suite *SetterTestSuite) TestSetCredentialToDefault() {
	for _, t := range []struct {
		testName string
		env      env.Holder
	}{
		{"env informed", *suite.envHolderDefault},
		{"env not informed", *suite.envHolderNil},
	} {
		suite.Run(t.testName, func() {
			defer os.RemoveAll(suite.HomePath)
			envFinderMock := new(mocks.EnvFinderMock)
			filePathExpectedCreated := filepath.Join(suite.HomePath, credentialDir, suite.envHolderDefault.Current, suite.DetailCredentialInfo.Service)

			envFinderMock.On("Find").Return(t.env, nil)
			setter := NewSetter(suite.HomePath, envFinderMock, dirManager)

			suite.NoFileExists(filePathExpectedCreated)

			response := setter.Set(*suite.DetailCredentialInfo)

			suite.Nil(response)
			suite.FileExists(filePathExpectedCreated)

			nameExpected := fmt.Sprintf("\"username\":\"%s\"", suite.DetailCredentialInfo.Username)
			data, err := suite.fileInfo(filePathExpectedCreated)
			suite.Nil(err)
			suite.Contains(data, nameExpected)
		})
	}
}

func TestSetterTestSuite(t *testing.T) {
	suite.Run(t, new(SetterTestSuite))
}
