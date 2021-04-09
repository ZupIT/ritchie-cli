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
package gitlab

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
			want: git.Tags{
				{
					Name:        "v1.0.0",
					Description: "Test",
				},
			},
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
				tagsUrl: "htttp://yourhost.com/username/repo/",
				token:   "some_token",
			},
			want:    git.Tags{},
			wantErr: "Get \"htttp://yourhost.com/username/repo/\": unsupported protocol scheme \"htttp\"",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("TokenHeader").Return(tt.info.token)
			git.On("TagsUrl").Return(tt.info.tagsUrl)
			re := NewRepoManager(tt.client)
			got, err := re.Tags(git)

			assert.Equal(t, tt.want, got)
			if err != nil {
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
		wantErr bool
	}{
		{
			name:   "Run with success",
			client: mockServer.Client(),
			info: info{
				zipUrl: mockServer.URL,
				token:  "some_token",
			},
			version: "v1.0.0",
			want:    "zipValue",
			wantErr: false,
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
				zipUrl: "htttp://yourhost.com/username/repo/",
				token:  "some_token",
			},
			version: "0.0.1",
			want:    "Get \"htttp://yourhost.com/username/repo/\": unsupported protocol scheme \"htttp\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("TokenHeader").Return(tt.info.token)
			git.On("ZipUrl", tt.version).Return(tt.info.zipUrl)
			re := RepoManager{client: tt.client}
			got, err := re.Zipball(git, tt.version)
			if err != nil {
				assert.EqualError(t, err, tt.want)
			} else {
				result, err := ioutil.ReadAll(got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(result))
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
			want: git.Tag{Name: "1.0.1", Description: "New golang formula"},
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
				zipUrl: "htttp://yourhost.com/username/repo/",
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
			git.On("TokenHeader").Return(tt.info.token)
			git.On("LatestTagUrl").Return(tt.info.latestTagUrl)
			re := NewRepoManager(tt.client)
			got, err := re.LatestTag(git)

			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

const (
	PayloadListAllTags = `[
    {
        "name": "v1.0.0",
        "tag_name": "v1.0.0",
        "description": "Test",
        "description_html": "<p data-sourcepos=\"1:1-1:4\" dir=\"auto\">Test</p>",
        "created_at": "2020-08-04T12:34:36.398Z",
        "released_at": "2020-08-04T12:34:36.398Z",
        "author": {
            "id": 2418203,
            "name": "User name",
            "username": "username",
            "state": "active",
            "avatar_url": "https://secure.gravatar.com/avatar/0ff2eaae88cd45073b3cf17b13151f15?s=80&d=identicon",
            "web_url": "https://gitlab.com/username"
        },
        "commit": {
            "id": "8d5ec4f376066836f21110766314210b5e21bbd2",
            "short_id": "8d5ec4f3",
            "created_at": "2020-08-04T11:56:29.000+00:00",
            "parent_ids": [],
            "title": "Initial commit",
            "message": "Initial commit",
            "author_name": "user name",
            "author_email": "user@user.com",
            "authored_date": "2020-08-04T11:56:29.000+00:00",
            "committer_name": "user name",
            "committer_email": "user@user.com",
            "committed_date": "2020-08-04T11:56:29.000+00:00",
            "web_url": "https://gitlab.com/username/ritchie-formulas/-/commit/8d5ec4f376066836f21110766314210b5e21bbd2"
        },
        "upcoming_release": false,
        "commit_path": "/username/ritchie-formulas/-/commit/8d5ec4f376066836f21110766314210b5e21bbd2",
        "tag_path": "/username/ritchie-formulas/-/tags/1.0.0",
        "assets": {
            "count": 4,
            "sources": [
                {
                    "format": "zip",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.0/ritchie-formulas-1.0.0.zip"
                },
                {
                    "format": "tar.gz",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.0/ritchie-formulas-1.0.0.tar.gz"
                },
                {
                    "format": "tar.bz2",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.0/ritchie-formulas-1.0.0.tar.bz2"
                },
                {
                    "format": "tar",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.0/ritchie-formulas-1.0.0.tar"
                }
            ],
            "links": []
        },
        "evidences": [],
        "_links": {
            "self": "https://gitlab.com/username/ritchie-formulas/-/releases/1.0.0",
            "edit_url": "https://gitlab.com/username/ritchie-formulas/-/releases/1.0.0/edit"
        }
    }
]`

	PayloadListLastTags = `[
    {
        "name": "1.0.1",
        "tag_name": "1.0.1",
        "description": "New golang formula",
        "description_html": "<p data-sourcepos=\"1:1-1:18\" dir=\"auto\">New golang formula</p>",
        "created_at": "2020-08-04T15:49:24.101Z",
        "released_at": "2020-08-04T15:49:24.101Z",
        "author": {
            "id": 2418203,
            "name": "User Name",
            "username": "username",
            "state": "active",
            "avatar_url": "https://secure.gravatar.com/avatar/0ff2eaae88cd45073b3cf17b13151f15?s=80&d=identicon",
            "web_url": "https://gitlab.com/username"
        },
        "commit": {
            "id": "03f883dfa3821672ef74f3bcc12ae2c83b068dd8",
            "short_id": "03f883df",
            "created_at": "2020-08-04T12:48:09.000-03:00",
            "parent_ids": [
                "8d5ec4f376066836f21110766314210b5e21bbd2"
            ],
            "title": "Create formula go",
            "message": "Create formula go\n",
            "author_name": "User Name",
            "author_email": "user@user.com",
            "authored_date": "2020-08-04T12:48:09.000-03:00",
            "committer_name": "User Name",
            "committer_email": "user@user.com",
            "committed_date": "2020-08-04T12:48:09.000-03:00",
            "web_url": "https://gitlab.com/username/ritchie-formulas/-/commit/03f883dfa3821672ef74f3bcc12ae2c83b068dd8"
        },
        "upcoming_release": false,
        "commit_path": "/username/ritchie-formulas/-/commit/03f883dfa3821672ef74f3bcc12ae2c83b068dd8",
        "tag_path": "/username/ritchie-formulas/-/tags/1.0.1",
        "assets": {
            "count": 4,
            "sources": [
                {
                    "format": "zip",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.1/ritchie-formulas-1.0.1.zip"
                },
                {
                    "format": "tar.gz",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.1/ritchie-formulas-1.0.1.tar.gz"
                },
                {
                    "format": "tar.bz2",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.1/ritchie-formulas-1.0.1.tar.bz2"
                },
                {
                    "format": "tar",
                    "url": "https://gitlab.com/username/ritchie-formulas/-/archive/1.0.1/ritchie-formulas-1.0.1.tar"
                }
            ],
            "links": []
        },
        "evidences": [
            {
                "sha": "1ee3430fc63f6fb9e2a3e6ffe62fcbf90e3d273372ad",
                "filepath": "https://gitlab.com/username/ritchie-formulas/-/releases/1.0.1/evidences/285309.json",
                "collected_at": "2020-08-04T15:49:24.158Z"
            }
        ],
        "_links": {
            "self": "https://gitlab.com/username/ritchie-formulas/-/releases/1.0.1",
            "edit_url": "https://gitlab.com/username/ritchie-formulas/-/releases/1.0.1/edit"
        }
    }
]`

	PayloadEmptyTags = `[]`
)

type info struct {
	zipUrl       string
	tagsUrl      string
	latestTagUrl string
	token        string
}
