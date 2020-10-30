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
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ---------- > SETUP <----------
type SetterTestSuite struct {
	suite.Suite

	HomePath          string
	TestContextFinder *contextFinderMock

	contextHolderNil     *rcontext.ContextHolder
	contextHolderDefault *rcontext.ContextHolder
	contextHolderProd    *rcontext.ContextHolder

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
	suite.TestContextFinder = new(contextFinderMock)
	suite.contextHolderNil = &rcontext.ContextHolder{Current: ""}
	suite.contextHolderDefault = &rcontext.ContextHolder{Current: "default"}
	suite.contextHolderProd = &rcontext.ContextHolder{Current: "prod", All: []string{"defauld", "prod"}}
	suite.DetailCredentialInfo = detailExample
}

func (suite *SetterTestSuite) SetupTest() {
	os.RemoveAll(suite.HomePath)
}

// func (suite *SetterTestSuite) BeforeTest(suiteName, testName string) {
// 	os.Remove(suite.HomePath)
// }

func (suite *SetterTestSuite) AfterTest(suiteName, testName string) {
	os.RemoveAll(suite.HomePath)
}

// <---------------------------->

// ---------- > TESTS <----------

func (suite *SetterTestSuite) TestSetCredentialToContextDefault() {
	filePathExpectedCreated := File(suite.HomePath, suite.contextHolderDefault.Current, suite.DetailCredentialInfo.Service)

	suite.TestContextFinder.On("Find").Return(*suite.contextHolderDefault, nil)
	setter := NewSetter(suite.HomePath, suite.TestContextFinder)

	suite.NoFileExists(filePathExpectedCreated)

	response := setter.Set(*suite.DetailCredentialInfo)

	suite.Nil(response)
	suite.FileExists(filePathExpectedCreated)
}

func (suite *SetterTestSuite) TestSetCredentialToContextDefaultWhenContextNotInformed() {
	filePathExpectedCreated := File(suite.HomePath, suite.contextHolderDefault.Current, suite.DetailCredentialInfo.Service)

	suite.TestContextFinder.On("Find").Return(*suite.contextHolderNil, nil)
	setter := NewSetter(suite.HomePath, suite.TestContextFinder)

	suite.NoFileExists(filePathExpectedCreated)

	response := setter.Set(*suite.DetailCredentialInfo)

	suite.Nil(response)
	suite.FileExists(filePathExpectedCreated)
}

func TestSetterTestSuite(t *testing.T) {
	suite.Run(t, new(SetterTestSuite))
}

// <---------------------------->

// ---------- > MOCKS <----------
type contextFinderMock struct {
	mock.Mock
}

func (cf *contextFinderMock) Find() (rcontext.ContextHolder, error) {

	args := cf.Called()
	return args.Get(0).(rcontext.ContextHolder), args.Error(1)

}

// <---------------------------->
