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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

type StubResolverVersions struct {
	stableVersion    func() (string, error)
	updateCache      func() error
	verifyNewVersion func(resolve Resolver, currentVersion string) string
}

func (r StubResolverVersions) StableVersion() (string, error) {
	return r.stableVersion()
}

func (r StubResolverVersions) UpdateCache() error {
	return r.updateCache()
}

func (r StubResolverVersions) VerifyNewVersion(resolve Resolver, currentVersion string) string {
	return r.verifyNewVersion(resolve, currentVersion)
}

func TestDefaultVersionResolver_StableVersion(t *testing.T) {

	// case 1 Should get stableVersion
	// expectedResultCase1 := "1.0.0"
	//
	// mockHttpCase1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(200)
	// 	_, _ = w.Write([]byte(expectedResultCase1 + "\n"))
	// }))
	//

	// case 4 Should not get stableVersion from expired cache
	// expectedResultCase4 := "2.5.6"
	okStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("2.0.4" + "\n"))
	}))

	internalErrorStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	// mockHttpCase4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(500)
	// }))

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
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"stableVersion\":\"2.0.4\",\"expiresAt\":1598396366}"), nil
					},
				},
				HttpClient: &http.Client{},
			},
			want:    "2.0.4",
			wantErr: false,
		},
		{
			name: "error reading from cache",
			in: in{
				StableVersionUrl: "any value",
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), errors.New("reading error")
					},
				},
				HttpClient: &http.Client{},
			},
			wantErr: true,
		},
		{
			name: "error on unmarshal cache data",
			in: in{
				StableVersionUrl: "any value",
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"expiresAt\":\"2.0.4\"}"), nil
					},
				},
				HttpClient: &http.Client{},
			},
			wantErr: true,
		},
		{
			name: "error on request stable version with no url",
			in: in{
				StableVersionUrl: "",
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"stableVersion\":\"2.0.4\",\"expiresAt\":1598396}"), nil
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
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"stableVersion\":\"2.0.4\",\"expiresAt\":1598396}"), nil
					},
				},
				HttpClient: internalErrorStub.Client(),
			},
			wantErr: true,
		},
		{
			name: "success case when cache expired",
			in: in{
				StableVersionUrl: okStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"stableVersion\":\"2.0.4\",\"expiresAt\":1598396}"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				HttpClient: okStub.Client(),
			},
			want:    "2.0.4",
			wantErr: false,
		},
		{
			name: "error on saving cache",
			in: in{
				StableVersionUrl: okStub.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"stableVersion\":\"2.0.4\",\"expiresAt\":1598396}"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return errors.New("error on saving cache")
					},
				},
				HttpClient: okStub.Client(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewManager(tt.in.StableVersionUrl, tt.in.file, tt.in.HttpClient)
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

// func TestVerifyNewVersion(t *testing.T) {
// 	type args struct {
// 		resolve        Resolver
// 		currentVersion string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "Should return empty when current version equals to stableVersion",
// 			args: args{
// 				resolve: StubResolverVersions{
// 					stableVersion: func() (string, error) {
// 						return "1.0.0", nil
// 					},
// 				},
// 				currentVersion: "1.0.0",
// 			},
// 			want: "",
// 		},
// 		{
// 			name: "Should return msg when current version it not equals to stableVersion",
// 			args: args{
// 				resolve: StubResolverVersions{
// 					stableVersion: func() (string, error) {
// 						return "1.0.1", nil
// 					},
// 				},
// 				currentVersion: "1.0.0",
// 			},
// 			want: MsgRitUpgrade,
// 		},
// 		{
// 			name: "Should return empty on error in StableVersion ",
// 			args: args{
// 				resolve: StubResolverVersions{
// 					stableVersion: func() (string, error) {
// 						return "", errors.New("any error")
// 					},
// 				},
// 				currentVersion: "1.0.0",
// 			},
// 			want: "",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := NewManager("version url",
// 				sMocks.FileWriteReadExisterCustomMock{},
// 				http.DefaultClient)
// 			if got := r.VerifyNewVersion(tt.args.currentVersion); got != tt.want {
// 				t.Errorf("VerifyNewVersion() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestDefaultVersionResolver_UpdateCache(t *testing.T) {

	expectedResultCase1 := "1.0.0"

	mockHttpCase1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(expectedResultCase1 + "\n"))
	}))

	type fields struct {
		StableVersionUrl string
		file             stream.FileWriteReadExister
		HttpClient       *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error on requestStableVersion",
			fields: fields{
				StableVersionUrl: "any url",
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				HttpClient: &http.Client{},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				StableVersionUrl: mockHttpCase1.URL,
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					ReadMock: func(path string) ([]byte, error) {
						return []byte("some data"), nil
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				HttpClient: mockHttpCase1.Client(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewManager(tt.fields.StableVersionUrl, tt.fields.file, tt.fields.HttpClient)
			if err := r.UpdateCache(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
