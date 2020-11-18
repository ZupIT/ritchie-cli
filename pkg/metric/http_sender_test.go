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
)

func TestSendManagerHttp_Send(t *testing.T) {
	type in struct {
	}
	tests := []struct {
		name string
		in   in
	}{
		{
			name: "success",
			in:   in{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := server()
			defer server.Close()
			data := DataCollectorMock{
				CollectCommandDataMock: func() Command {
					return Command{}
				},
				CollectUserStateMock: func() User {
					return User{}
				},
			}
			checker := CheckerMock{
				func() bool {
					return true
				},
			}
			httpSender := NewHttpSender(server.URL, server.Client(), data, checker)
			httpSender.SendUserState("tt.in.APIData")
		})
	}
}

func server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

type DataCollectorMock struct {
	CollectCommandDataMock func() Command
	CollectUserStateMock   func() User
}

func (dc DataCollectorMock) CollectCommandData(commandExecutionTime float64, commandError ...string) Command {
	return dc.CollectCommandDataMock()
}

func (dc DataCollectorMock) CollectUserState(ritVersion string) User {
	return dc.CollectUserStateMock()
}

type CheckerMock struct {
	CheckMock func() bool
}

func (c CheckerMock) Check() bool {
	return true
}
