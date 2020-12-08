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

import (
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type ContextFinderMock struct {
	mock.Mock
}

func (cf *ContextFinderMock) Find() (rcontext.ContextHolder, error) {
	args := cf.Called()

	return args.Get(0).(rcontext.ContextHolder), args.Error(1)
}

type DetailManagerMock struct {
	mock.Mock
}

func (d *DetailManagerMock) LatestTag(repo formula.Repo) string {
	args := d.Called(repo)

	return args.String(0)
}

type SenderMock struct {
	mock.Mock
}

func (s *SenderMock) SendUserState(ritVersion string) {
	s.Called()
}

func (s *SenderMock) SendCommandData(cmd metric.SendCommandDataParams) {
	s.Called()
}
