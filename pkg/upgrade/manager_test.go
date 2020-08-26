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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inconshreveable/go-update"
)

type stubUpdater struct {
	apply func(reader io.Reader, opts update.Options) error
}

func (u stubUpdater) Apply(reader io.Reader, opts update.Options) error {
	return u.apply(reader, opts)
}

func TestDefaultManager_Run(t *testing.T) {
	type in struct {
		updater    updater
		upgradeUrl string
	}
	type args struct {
	}
	tests := []struct {
		name    string
		in      in
		args    args
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				updater:    NewDefaultUpdater(),
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
			},
			wantErr: false,
		},
		{
			name: "Should return err when url is empty",
			in: in{
				updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
				upgradeUrl: "",
			},
			wantErr: true,
		},
		{
			name: "Should return err when happening err when perform get",
			in: in{
				updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
				upgradeUrl: "some url",
			},
			wantErr: true,
		},
		{
			name: "Should return err when get return not 200 code",
			in: in{
				updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
				})).URL,
			},
			wantErr: true,
		},
		{
			name: "Should return err when fail to apply",
			in: in{
				updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return errors.New("some error")
				}},
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewDefaultManager(tt.in.updater)
			if err := m.Run(tt.in.upgradeUrl); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
