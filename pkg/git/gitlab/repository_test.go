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
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

func TestTags(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(PayloadListAllTags))
	}))

	mockServerThatFail := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
	}))

	type in struct {
		client *http.Client
		info   git.RepoInfo
	}
	tests := []struct {
		name    string
		in      in
		want    git.Tags
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				client: mockServer.Client(),
				info: RepoInfoCustomMock{
					tagsUrl: mockServer.URL,
					token:   "some_token",
				},
			},
			want: git.Tags{
				{
					Name: "v1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Return err when request fail",
			in: in{
				client: mockServerThatFail.Client(),
				info: RepoInfoCustomMock{
					tagsUrl: mockServerThatFail.URL,
				},
			},
			want:    git.Tags{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := NewRepoManager(tt.in.client)

			got, err := re.Tags(tt.in.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tags() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipball(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		data := `zipValue`
		_, _ = writer.Write([]byte(data))
	}))

	type in struct {
		client  *http.Client
		info    git.RepoInfo
		version string
	}
	tests := []struct {
		name    string
		in      in
		want    string
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				client: mockServer.Client(),
				info: RepoInfoCustomMock{
					zipUrl: func(version string) string {
						return mockServer.URL
					},
					token: "some_token",
				},
				version: "v1.0.0",
			},
			want:    "zipValue",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := RepoManager{
				client: tt.in.client,
			}
			got, err := re.Zipball(tt.in.info, tt.in.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zipball() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result, err := ioutil.ReadAll(got)
			if err != nil {
				t.Errorf("fail to parse result")
			}

			if string(result) != tt.want {
				t.Errorf("Zipball() got = %v, want %v", got, tt.want)
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
)

type RepoInfoCustomMock struct {
	zipUrl       func(version string) string
	tagsUrl      string
	latestTagUrl string
	tokenHeader  string
	token        string
}

func (m RepoInfoCustomMock) ZipUrl(version string) string {
	return m.zipUrl(version)
}

func (m RepoInfoCustomMock) TagsUrl() string {
	return m.tagsUrl
}

func (m RepoInfoCustomMock) LatestTagUrl() string {
	return m.latestTagUrl
}

func (m RepoInfoCustomMock) TokenHeader() string {
	return m.tokenHeader
}

func (m RepoInfoCustomMock) Token() string {
	return m.token
}
