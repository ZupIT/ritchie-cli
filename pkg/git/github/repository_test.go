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

package github

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/git"
)

const (
	PayloadListLastTags = `
{
  "url": "https://api.github.com/repos/octocat/Hello-World/releases/1",
  "html_url": "https://github.com/octocat/Hello-World/releases/v1.0.0",
  "assets_url": "https://api.github.com/repos/octocat/Hello-World/releases/1/assets",
  "upload_url": "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
  "tarball_url": "https://api.github.com/repos/octocat/Hello-World/tarball/v1.0.0",
  "zipball_url": "https://api.github.com/repos/octocat/Hello-World/zipball/v1.0.0",
  "id": 1,
  "node_id": "MDc6UmVsZWFzZTE=",
  "tag_name": "v1.0.0",
  "target_commitish": "master",
  "name": "v1.0.0",
  "body": "Description of the release",
  "draft": false,
  "prerelease": false,
  "created_at": "2013-02-27T19:35:32Z",
  "published_at": "2013-02-27T19:35:32Z",
  "author": {
    "login": "octocat",
    "id": 1,
    "node_id": "MDQ6VXNlcjE=",
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/octocat/followers",
    "following_url": "https://api.github.com/users/octocat/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat/orgs",
    "repos_url": "https://api.github.com/users/octocat/repos",
    "events_url": "https://api.github.com/users/octocat/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat/received_events",
    "type": "User",
    "site_admin": false
  },
  "assets": [
    {
      "url": "https://api.github.com/repos/octocat/Hello-World/releases/assets/1",
      "browser_download_url": "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
      "id": 1,
      "node_id": "MDEyOlJlbGVhc2VBc3NldDE=",
      "name": "example.zip",
      "label": "short description",
      "state": "uploaded",
      "content_type": "application/zip",
      "size": 1024,
      "download_count": 42,
      "created_at": "2013-02-27T19:35:32Z",
      "updated_at": "2013-02-27T19:35:32Z",
      "uploader": {
        "login": "octocat",
        "id": 1,
        "node_id": "MDQ6VXNlcjE=",
        "avatar_url": "https://github.com/images/error/octocat_happy.gif",
        "gravatar_id": "",
        "url": "https://api.github.com/users/octocat",
        "html_url": "https://github.com/octocat",
        "followers_url": "https://api.github.com/users/octocat/followers",
        "following_url": "https://api.github.com/users/octocat/following{/other_user}",
        "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
        "organizations_url": "https://api.github.com/users/octocat/orgs",
        "repos_url": "https://api.github.com/users/octocat/repos",
        "events_url": "https://api.github.com/users/octocat/events{/privacy}",
        "received_events_url": "https://api.github.com/users/octocat/received_events",
        "type": "User",
        "site_admin": false
      }
    }
  ]
}`
	PayloadListAllTags = `
[
  {
    "url": "https://api.github.com/repos/octocat/Hello-World/releases/1",
    "html_url": "https://github.com/octocat/Hello-World/releases/v1.0.0",
    "assets_url": "https://api.github.com/repos/octocat/Hello-World/releases/1/assets",
    "upload_url": "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
    "tarball_url": "https://api.github.com/repos/octocat/Hello-World/tarball/v1.0.0",
    "zipball_url": "https://api.github.com/repos/octocat/Hello-World/zipball/v1.0.0",
    "id": 1,
    "node_id": "MDc6UmVsZWFzZTE=",
    "tag_name": "v1.0.0",
    "target_commitish": "master",
    "name": "v1.0.0",
    "body": "Description of the release",
    "draft": false,
    "prerelease": false,
    "created_at": "2013-02-27T19:35:32Z",
    "published_at": "2013-02-27T19:35:32Z",
    "author": {
      "login": "octocat",
      "id": 1,
      "node_id": "MDQ6VXNlcjE=",
      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
      "gravatar_id": "",
      "url": "https://api.github.com/users/octocat",
      "html_url": "https://github.com/octocat",
      "followers_url": "https://api.github.com/users/octocat/followers",
      "following_url": "https://api.github.com/users/octocat/following{/other_user}",
      "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
      "organizations_url": "https://api.github.com/users/octocat/orgs",
      "repos_url": "https://api.github.com/users/octocat/repos",
      "events_url": "https://api.github.com/users/octocat/events{/privacy}",
      "received_events_url": "https://api.github.com/users/octocat/received_events",
      "type": "User",
      "site_admin": false
    },
    "assets": [
      {
        "url": "https://api.github.com/repos/octocat/Hello-World/releases/assets/1",
        "browser_download_url": "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
        "id": 1,
        "node_id": "MDEyOlJlbGVhc2VBc3NldDE=",
        "name": "example.zip",
        "label": "short description",
        "state": "uploaded",
        "content_type": "application/zip",
        "size": 1024,
        "download_count": 42,
        "created_at": "2013-02-27T19:35:32Z",
        "updated_at": "2013-02-27T19:35:32Z",
        "uploader": {
          "login": "octocat",
          "id": 1,
          "node_id": "MDQ6VXNlcjE=",
          "avatar_url": "https://github.com/images/error/octocat_happy.gif",
          "gravatar_id": "",
          "url": "https://api.github.com/users/octocat",
          "html_url": "https://github.com/octocat",
          "followers_url": "https://api.github.com/users/octocat/followers",
          "following_url": "https://api.github.com/users/octocat/following{/other_user}",
          "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
          "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
          "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
          "organizations_url": "https://api.github.com/users/octocat/orgs",
          "repos_url": "https://api.github.com/users/octocat/repos",
          "events_url": "https://api.github.com/users/octocat/events{/privacy}",
          "received_events_url": "https://api.github.com/users/octocat/received_events",
          "type": "User",
          "site_admin": false
        }
      }
    ]
  }
]`
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

func TestLatestTag(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(PayloadListLastTags))
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
		want    git.Tag
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
				client: mockServer.Client(),
				info: RepoInfoCustomMock{
					latestTagUrl: mockServer.URL,
					token:        "some_token",
				},
			},
			want:    git.Tag{Name: "v1.0.0", Description: "Description of the release"},
			wantErr: false,
		},
		{
			name: "Return err when request fail",
			in: in{
				client: mockServerThatFail.Client(),
				info: RepoInfoCustomMock{
					latestTagUrl: mockServerThatFail.URL,
				},
			},
			want:    git.Tag{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := NewRepoManager(tt.in.client)

			got, err := re.LatestTag(tt.in.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("LatestTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LatestTag() got = %v, want %v", got, tt.want)
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
				version: "1.0.0",
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
