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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/internal/pkg/config"
	"github.com/ZupIT/ritchie-cli/pkg/api"
)

func TestTranslation(t *testing.T) {
	ritHome := api.RitchieHomeDir()
	c := config.NewManager(ritHome)
	defaultConfig, _ := c.Read()

	type in struct {
		msgID string
		lang  string
	}

	tests := []struct {
		name string
		in   in
		want string
	}{
		{
			name: "success with default language",
			in: in{
				msgID: "init.cmd.description",
			},
			want: "Initialize rit configuration",
		},
		{
			name: "success with selected language",
			in: in{
				msgID: "init.cmd.description",
				lang:  "Portuguese",
			},
			want: "Inicializar as configurações do rit",
		},
		{
			name: "use default language when message not found for current lang",
			in: in{
				msgID: "init.successful",
				lang:  "Portuguese",
			},
			want: "\n✅ Initialization successful!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			if in.lang != "" {
				configs := config.Configs{
					Language: in.lang,
				}
				_ = c.Write(configs)
			}

			message := T(in.msgID)

			assert.Equal(t, tt.want, message)

			_ = c.Write(defaultConfig)
		})
	}

}
