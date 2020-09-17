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
	"bytes"
	"encoding/json"
	"net/http"
)

var _ Sender = SendManagerHttp{}

type SendManagerHttp struct {
	URL    string
	client *http.Client
}

func NewHttpSender(url string, client *http.Client) SendManagerHttp {
	return SendManagerHttp{
		URL:    url,
		client: client,
	}
}

func (sm SendManagerHttp) Send(APIData APIData) {
	reqBody, err := json.Marshal(&APIData)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, sm.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	_, _ = sm.client.Do(req)
}
