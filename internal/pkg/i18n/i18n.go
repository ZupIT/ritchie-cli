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

package i18n

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/ZupIT/ritchie-cli/internal/pkg/config"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
)

var Langs = map[string]string{
	"English":    "en",
	"Portuguese": "pt_BR",
}

type Translation struct {
	bundle *i18n.Bundle
	config config.Reader
}

// T receives msgID to identify your message and parameters
// to format your message, you can use parameters like fmt.Sprintf.
// Examples:
// message := i18n.T("my.message.without.params")
// message := i18n.T("my.message.with.params", "name", "test")
var T = func(msgID string, params ...interface{}) string {
	return t.resolver(msgID, params...)
}

var t *Translation

func init() {
	ritHome := api.RitchieHomeDir()
	c := config.NewManager(ritHome)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for i, v := range Langs {
		file := fmt.Sprintf("resources/i18n/%s.toml", v)

		bytes, err := Asset(file)
		if err != nil {
			panic(fmt.Sprintf("Error to load %q translation", i))
		}

		bundle.MustParseMessageFileBytes(bytes, file)
	}

	t = &Translation{
		bundle: bundle,
		config: c,
	}
}

func (t *Translation) resolver(msgID string, params ...interface{}) string {
	var msg string
	c, _ := t.config.Read()
	loc := i18n.NewLocalizer(t.bundle, Langs[c.Language])
	locConfig := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{},
	}

	// If OS is Windows we will try
	// to find the specific message,
	// case not found message for
	// Windows we will try to find
	// the default message.
	if osutil.IsWindows() {
		id := fmt.Sprintf("%s_windows", msgID)
		locConfig.DefaultMessage.ID = id

		msg, _ = loc.Localize(locConfig)
		if msg != "" {
			msg = fmt.Sprintf(msg, params...)
			return msg
		}
	}

	locConfig.DefaultMessage.ID = msgID
	msg, _ = loc.Localize(locConfig)
	msg = fmt.Sprintf(msg, params...)

	return msg
}
