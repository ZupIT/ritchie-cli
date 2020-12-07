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

const (
	commandsURL  = "/commands"
	userStateURL = "/user"
)

type SendManagerHttp struct {
	URL       string
	client    *http.Client
	collector Collector
	checker   Checker
}

func NewHttpSender(
	url string,
	client *http.Client,
	collector Collector,
	checker Checker,
) SendManagerHttp {
	return SendManagerHttp{
		URL:       url,
		client:    client,
		checker:   checker,
		collector: collector,
	}
}

func (sm SendManagerHttp) SendUserState(ritVersion string) {
	if !sm.checker.Check() {
		return
	}
	userState := sm.collector.CollectUserState(ritVersion)
	reqBody, err := json.Marshal(&userState)
	if err != nil {
		return
	}
	sm.doRequest(reqBody, userStateURL)
}

func (sm SendManagerHttp) SendCommandData(cmd SendCommandDataParams) {
	if !sm.checker.Check() {
		return
	}
	command := sm.collector.CollectCommandData(cmd.ExecutionTime, cmd.Error)
	reqBody, err := json.Marshal(&command)
	if err != nil {
		return
	}
	sm.doRequest(reqBody, commandsURL)
}

func (sm SendManagerHttp) doRequest(reqBody []byte, URLsuffix string) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		sm.URL+URLsuffix,
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
