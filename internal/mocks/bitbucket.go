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

package mocks

import "github.com/stretchr/testify/mock"

type RepoInfo struct {
	mock.Mock
}

func (ri *RepoInfo) ZipUrl(version string) string {
	args := ri.Called(version)
	return args.String(0)
}

func (ri *RepoInfo) TagsUrl() string {
	args := ri.Called()
	return args.String(0)
}

func (ri *RepoInfo) LatestTagUrl() string {
	args := ri.Called()
	return args.String(0)
}

func (ri *RepoInfo) TokenHeader() string {
	args := ri.Called()
	return args.String(0)
}

func (ri *RepoInfo) Token() string {
	args := ri.Called()
	return args.String(0)
}
