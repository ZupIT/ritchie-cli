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

package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestSendManagerHttp_Send(t *testing.T) {
	type in struct {
		checkerReturn bool
	}

	tests := []struct {
		name string
		in   in
	}{
		{
			name: "success run",
			in:   in{checkerReturn: true},
		},
		{
			name: "success run with false checker",
			in:   in{checkerReturn: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := server()
			defer server.Close()
			checkerMock := &CheckerMock{}
			checkerMock.On("Check").Return(tt.in.checkerReturn)

			m := &DataCollectorMock{}
			m.On("CollectUserState", mock.Anything).Return(User{})
			m.On("CollectCommandData", mock.Anything, mock.Anything).Return(Command{})

			httpSender := NewHttpSender(server.URL, server.Client(), m, checkerMock)
			httpSender.SendUserState("2.0.1")
			httpSender.SendCommandData(SendCommandDataParams{})
		})
	}
}

func server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

type DataCollectorMock struct {
	mock.Mock
}

func (dc *DataCollectorMock) CollectCommandData(
	commandExecutionTime float64,
	commandError ...string,
) Command {
	args := dc.Called(commandExecutionTime, commandError)
	return args.Get(0).(Command)
}

func (dc *DataCollectorMock) CollectUserState(ritVersion string) User {
	args := dc.Called(ritVersion)
	return args.Get(0).(User)
}

type CheckerMock struct {
	mock.Mock
}

func (c *CheckerMock) Check() bool {
	return c.Called().Get(0).(bool)
}
