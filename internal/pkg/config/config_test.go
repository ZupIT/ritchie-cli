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

package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	type in struct {
		ritHome string
		configs Configs
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				ritHome: os.TempDir(),
				configs: Configs{
					Language: "English",
				},
			},
		},
		{
			name: "invalid path",
			in: in{
				ritHome: "/invalid",
				configs: Configs{
					Language: "English",
				},
			},
			want: errors.New("open /invalid/configs.toml: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewManager(tt.in.ritHome)

			got := config.Write(tt.in.configs)
			if tt.want != nil {
				assert.Equal(t, tt.want.Error(), got.Error())
			} else {
				filePath := filepath.Join(os.TempDir(), File)
				assert.FileExists(t, filePath)

				var configs Configs
				_, _ = toml.DecodeFile(filePath, &configs)
				assert.Equal(t, tt.in.configs, configs)
			}

			_ = os.Remove(filepath.Join(tt.in.ritHome, File))
		})
	}
}

func TestRead(t *testing.T) {
	type in struct {
		ritHome string
		configs *Configs
	}

	type out struct {
		configs Configs
		err     error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success",
			in: in{
				ritHome: os.TempDir(),
				configs: &Configs{
					Language: "English",
				},
			},
			out: out{
				configs: Configs{
					Language: "English",
					Tutorial: "",
					Metrics:  "",
				},
			},
		},
		{
			name: "config file not found",
			in: in{
				ritHome: "/invalid",
			},
			out: out{
				err: errors.New("open /invalid/configs.toml: no such file or directory"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewManager(tt.in.ritHome)
			if tt.in.configs != nil {
				_ = config.Write(*tt.in.configs)
			}

			got, err := config.Read()

			assert.Equal(t, tt.out.configs, got)
			if tt.out.err != nil {
				assert.Equal(t, tt.out.err.Error(), err.Error())
			}

			_ = os.Remove(filepath.Join(tt.in.ritHome, File))
		})
	}
}
