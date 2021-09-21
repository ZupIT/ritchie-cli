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
package zipremote

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(PayloadListAllTags))
	}))
	mockServerThatFail := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	}))
	mockServerNotFound := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	}))
	mockServerForbidden := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
	}))

	tests := []struct {
		name    string
		client  *http.Client
		info    info
		want    git.Tags
		wantErr string
	}{
		{
			name:   "Run with success",
			client: mockServer.Client(),
			info: info{
				tagsUrl: mockServer.URL,
				token:   "some_token",
			},
			want: git.Tags{},
		},
		{
			name:   "Return err when request fail",
			client: mockServerThatFail.Client(),
			info: info{
				tagsUrl: mockServerThatFail.URL,
			},
			want:    git.Tags{},
			wantErr: "400 Bad Request-",
		},
		{
			name:   "Return err when the protocol is invalid",
			client: mockServerThatFail.Client(),
			info: info{
				zipUrl: "htttp://provider.com/username/repo/",
				token:  "some_token",
			},
			want:    git.Tags{},
			wantErr: "Get \"\": unsupported protocol scheme \"\"",
		},
		{
			name:   "Return err when repo not found",
			client: mockServerNotFound.Client(),
			info: info{
				tagsUrl: mockServerNotFound.URL,
			},
			want:    git.Tags{},
			wantErr: git.ErrRepoNotFound.Error(),
		},
		{
			name:   "Return err when repo access denied",
			client: mockServerForbidden.Client(),
			info: info{
				tagsUrl: mockServerForbidden.URL,
			},
			want:    git.Tags{},
			wantErr: git.ErrRepoNotFound.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("TagsUrl").Return(tt.info.tagsUrl)
			re := NewRepoManager(tt.client)
			got, err := re.Tags(git)
			if err == nil {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want, got)
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestZipball(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		data := `zipValue`
		_, _ = writer.Write([]byte(data))
	}))
	mockServerThatFail := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	}))

	tests := []struct {
		name    string
		client  *http.Client
		info    info
		version string
		want    string
	}{
		{
			name:   "Run with success",
			client: mockServer.Client(),
			info: info{
				zipUrl: mockServer.URL,
				token:  "some_token",
			},
			version: "0.0.1",
			want:    "zipValue",
		},
		{
			name:   "Return err when request fail",
			client: mockServerThatFail.Client(),
			info: info{
				zipUrl: mockServerThatFail.URL,
				token:  "some_token",
			},
			version: "0.0.1",
			want:    "400 Bad Request-",
		},
		{
			name:   "Return err when the protocol is invalid",
			client: mockServerThatFail.Client(),
			info: info{
				zipUrl: "htttp://provider.com/username/repo/",
				token:  "some_token",
			},
			version: "0.0.1",
			want:    "Get \"htttp://provider.com/username/repo/\": unsupported protocol scheme \"htttp\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("ZipUrl", tt.version).Return(tt.info.zipUrl)
			re := RepoManager{client: tt.client}
			got, err := re.Zipball(git, tt.version)
			if err != nil {
				assert.EqualError(t, err, tt.want)
			} else {
				result, err := ioutil.ReadAll(got)
				if assert.Nil(t, err) {
					assert.Equal(t, tt.want, string(result))
				}
			}
		})
	}
}

func TestLatestTag(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(PayloadListLastTags))
	}))
	mockServerThatFail := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	}))
	mockServerNotFound := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(PayloadEmptyTags))
	}))

	tests := []struct {
		name    string
		client  *http.Client
		info    info
		want    git.Tag
		wantErr string
	}{
		{
			name:   "Run with success",
			client: mockServer.Client(),
			info: info{
				latestTagUrl: mockServer.URL,
				token:        "some_token",
			},
			want: git.Tag{},
		},
		{
			name:   "Return err when request fail",
			client: mockServerThatFail.Client(),
			info: info{
				latestTagUrl: mockServerThatFail.URL,
			},
			want:    git.Tag{},
			wantErr: "400 Bad Request-",
		},
		{
			name:   "Return err when not finding tags",
			client: mockServerNotFound.Client(),
			info: info{
				latestTagUrl: mockServerNotFound.URL,
			},
			want:    git.Tag{},
			wantErr: "release not found",
		},
		{
			name:   "Return err when the protocol is invalid",
			client: mockServerThatFail.Client(),
			info: info{
				zipUrl: "htttp://provider.com/username/repo/",
				token:  "some_token",
			},
			want:    git.Tag{},
			wantErr: "Get \"\": unsupported protocol scheme \"\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("LatestTagUrl").Return(tt.info.latestTagUrl)
			re := NewRepoManager(tt.client)
			got, err := re.LatestTag(git)
			if err == nil {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want, got)
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

const (
	PayloadListAllTags = `{}`

	PayloadListLastTags = `{}`

	PayloadEmptyTags = `{}`
)

type info struct {
	zipUrl       string
	tagsUrl      string
	latestTagUrl string
	token        string
}
