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

package upgrade

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/inconshreveable/go-update"
)

type Manager interface {
	Run(upgradeUrl string) error
}

type DefaultManager struct {
	Updater
}

func NewDefaultManager(Updater Updater) DefaultManager {
	return DefaultManager{Updater:Updater}
}

func (m DefaultManager) Run(upgradeUrl string) error {
	if upgradeUrl == "" {
		return errors.New("fail to resolve upgrade url")
	}

	resp, err := http.Get(upgradeUrl)
	if err != nil {
		return errors.New("fail to download stable version")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("fail to download stable version status:%d", resp.StatusCode)
	}

	err = m.Updater.Apply(resp.Body, update.Options{})
	if err != nil {
		return errors.New(
			"Fail to upgrade\n" +
				"Please try running this command again as root/Administrator\n" +
				"Example: sudo rit upgrade",
		)
	}
	return nil
}
