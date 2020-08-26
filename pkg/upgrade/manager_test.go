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
	type fields struct {
		Updater Updater
	}
	type args struct {
		upgradeUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Run with success",
			args: args{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
			},
			fields: fields{
				Updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
			},
			wantErr: false,
		},
		{
			name: "Should return err when url is empty",
			args: args{
				upgradeUrl: "",
			},
			fields: fields{
				Updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
			},
			wantErr: true,
		},
		{
			name: "Should return err when happening err when perform get",
			args: args{
				upgradeUrl: "some url",
			},
			fields: fields{
				Updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
			},
			wantErr: true,
		},
		{
			name: "Should return err when get return not 200 code",
			args: args{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
				})).URL,
			},
			fields: fields{
				Updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return nil
				}},
			},
			wantErr: true,
		},
		{
			name: "Should return err when fail to apply",
			args: args{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
			},
			fields: fields{
				Updater: stubUpdater{apply: func(reader io.Reader, opts update.Options) error {
					return errors.New("some error")
				}},
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := NewDefaultManager(tt.fields.Updater)
			if err := m.Run(tt.args.upgradeUrl); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
