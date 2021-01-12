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
)

var Langs = map[Lang]string{
	"English":   "en",
	"Português": "pt_BR",
}

type Lang string

func (l Lang) String() string {
	return string(l)
}

type Translation struct {
	bundle *i18n.Bundle
}

var T Translation

func NewTranslation() Translation {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for i, v := range Langs {
		file := fmt.Sprintf("resources/i18n/%s.toml", v)

		bytes, err := Asset(file)
		if err != nil {
			fmt.Printf("Error to load %q translation", i)
		}

		bundle.MustParseMessageFileBytes(bytes, v)
	}

	return Translation{
		bundle: bundle,
	}
}
