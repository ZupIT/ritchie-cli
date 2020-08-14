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
	"runtime"
	"testing"
	"time"
)

func TestSendManagerHttp_Send(t *testing.T) {
	type in struct {
		dataset Dataset
	}
	tests := []struct {
		name string
		in   in
	}{
		{
			name: "success",
			in: in{
				dataset: Dataset{
					Id:         "metric-id",
					UserId:     "user-id",
					Timestamp:  time.Now(),
					So:         runtime.GOOS,
					RitVersion: "2.0.0",
					Data:       nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := server()
			defer server.Close()

			httpSender := NewHttpSender(server.URL, server.Client())
			httpSender.Send(tt.in.dataset)
		})
	}
}

func server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
