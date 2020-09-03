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

package version

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

var defaultVersion = "2.0.4"

func buildStableBody(expiresAt int64) []byte {
	cache := stableVersionCache{
		Stable:    defaultVersion,
		ExpiresAt: expiresAt,
	}
	b, _ := json.Marshal(cache)
	return b
}

var okStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(defaultVersion + "\n"))
}))

var internalErrorStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}))

var wrongContentLengthStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "1")
}))

func TestManager_StableVersion(t *testing.T) {
	notExpiredCache := time.Now().Add(time.Hour).Unix()

	type in struct {
		StableVersionUrl string
		file             stream.FileWriteReadExister
		HttpClient       *http.Client
	}
	tests := []struct {
		name    string
		in      in
		want    string
		wantErr bool
	}{
		{
			name: "success get stableVersion from cache",
			in: in{
				StableVersionUrl: "any value",
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return buildStableBody(notExpiredCache), nil
					},
				},
			},
			want:    defaultVersion,
			wantErr: false,
		},
		{
			name: "error reading from cache",
			in: in{
				StableVersionUrl: "any value",
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), errors.New("reading error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error on unmarshal cache data",
			in: in{
				StableVersionUrl: "any value",
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"expiresAt\":\"2.0.4\"}"), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error on request stable version with no url",
			in: in{
				StableVersionUrl: "",
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return buildStableBody(1), nil
					},
				},
				HttpClient: okStub.Client(),
			},
			wantErr: true,
		},
		{
			name: "error on request stable status not 200",
			in: in{
				StableVersionUrl: internalErrorStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return buildStableBody(1), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error on saving cache",
			in: in{
				StableVersionUrl: okStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
					ReadMock: func(path string) ([]byte, error) {
						return buildStableBody(1), nil
					},
					WriteMock: func(path string, content []byte) error {
						return errors.New("error on saving cache")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error reading bytes from response body",
			in: in{
				StableVersionUrl: wrongContentLengthStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return buildStableBody(1), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "success case when cache expired",
			in: in{
				StableVersionUrl: okStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return buildStableBody(1), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
			},
			want:    defaultVersion,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewManager(tt.in.StableVersionUrl, tt.in.file)
			got, err := r.StableVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("StableVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("StableVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_UpdateCache(t *testing.T) {
	type in struct {
		StableVersionUrl string
		file             stream.FileWriteReadExister
		HttpClient       *http.Client
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "success get update cache",
			in: in{
				StableVersionUrl: okStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				HttpClient: okStub.Client(),
			},
			wantErr: false,
		},
		{
			name: "error on request during update cache",
			in: in{
				StableVersionUrl: internalErrorStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				HttpClient: internalErrorStub.Client(),
			},
			wantErr: true,
		},
		{
			name: "error on writing file",
			in: in{
				StableVersionUrl: okStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					WriteMock: func(path string, content []byte) error {
						return errors.New("writing file error")
					},
				},
				HttpClient: internalErrorStub.Client(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewManager(tt.in.StableVersionUrl, tt.in.file)
			err := r.UpdateCache()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_VerifyNewVersion(t *testing.T) {
	type in struct {
		current   string
		installed string
	}
	tests := []struct {
		name string
		in   in
		want string
	}{
		{
			name: "Should return empty when current version equals to current version",
			in: in{
				current:   defaultVersion,
				installed: defaultVersion,
			},
			want: "",
		},
		{
			name: "Should message when version is outdated",
			in: in{
				current:   defaultVersion,
				installed: "1.0.0",
			},
			want: MsgRitUpgrade,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewManager("version url",
				sMocks.FileWriteReadExisterMock{})
			if got := r.VerifyNewVersion(tt.in.current, tt.in.installed); got != tt.want {
				t.Errorf("VerifyNewVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
