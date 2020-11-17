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
	"context"
	"encoding/json"
	"net/http"
)

var _ Sender = SendManagerHttp{}

var (
	BasicUser = ""
	BasicPass = ""
)

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

func (sm SendManagerHttp) SendUserState(user User) {
	reqBody, err := json.Marshal(&user)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		// TODO endpoint for each send type
		sm.URL,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return
	}

	req.SetBasicAuth(BasicUser, BasicPass)
	req.Header.Add("Content-Type", "application/json")
	resp, err := sm.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func (sm SendManagerHttp) SendCommandData(command Command) {
	reqBody, err := json.Marshal(&command)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, sm.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	req.SetBasicAuth(BasicUser, BasicPass)
	req.Header.Add("Content-Type", "application/json")
	resp, err := sm.client.Do(req)
	if err != nil {
		return
	}
	if err := resp.Body.Close(); err != nil {
		return
	}
}
