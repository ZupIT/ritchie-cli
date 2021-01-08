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
package bitbucket

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
					Name:        "0.0.1",
					Description: "",
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
				version: "0.0.1",
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
			want:    git.Tag{Name: "0.0.1", Description: ""},
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

const (
	PayloadListAllTags = `{
		"pagelen": 10,
		"values": [
			{
				"name": "0.0.1",
				"links": {
					"commits": {
						"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commits/0.0.1"
					},
					"self": {
						"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/refs/tags/0.0.1"
					},
					"html": {
						"href": "https://bitbucket.org/username/ritchie-formulas/commits/tag/0.0.1"
					}
				},
				"tagger": {
					"raw": "Username <user@user.com>",
					"type": "author",
					"user": {
						"display_name": "Username",
						"uuid": "{b796025c-834e-475a-b669-7c6649156f64}",
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/users/%7Bb796025c-834e-475a-b669-7c6649156f64%7D"
							},
							"html": {
								"href": "https://bitbucket.org/%7Bb796025c-834e-475a-b669-7c6649156f64%7D/"
							},
							"avatar": {
								"href": "https://secure.gravatar.com/avatar/24215d1e57db99c8de7178e8d6f8e29d?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FLD-5.png"
							}
						},
						"nickname": "Username",
						"type": "user",
						"account_id": "5e67b3292a0bb00ce033d413"
					}
				},
				"date": "2020-12-22T18:33:52+00:00",
				"message": "Added tag 0.0.1 for changeset 224421a2680b",
				"type": "tag",
				"target": {
					"hash": "224421a2680b648a02a9feca6ef025fddadcaeff",
					"repository": {
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas"
							},
							"html": {
								"href": "https://bitbucket.org/username/ritchie-formulas"
							},
							"avatar": {
								"href": "https://bytebucket.org/ravatar/%7Bb2004d9e-f4b2-4784-99cf-b9c7d4322a5f%7D?ts=default"
							}
						},
						"type": "repository",
						"name": "ritchie-formulas",
						"full_name": "username/ritchie-formulas",
						"uuid": "{b2004d9e-f4b2-4784-99cf-b9c7d4322a5f}"
					},
					"links": {
						"self": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"comments": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff/comments"
						},
						"patch": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/patch/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"html": {
							"href": "https://bitbucket.org/username/ritchie-formulas/commits/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"diff": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/diff/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"approve": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff/approve"
						},
						"statuses": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff/statuses"
						}
					},
					"author": {
						"raw": "Username <user@user.com>",
						"type": "author",
						"user": {
							"display_name": "Username",
							"uuid": "{b796025c-834e-475a-b669-7c6649156f64}",
							"links": {
								"self": {
									"href": "https://api.bitbucket.org/2.0/users/%7Bb796025c-834e-475a-b669-7c6649156f64%7D"
								},
								"html": {
									"href": "https://bitbucket.org/%7Bb796025c-834e-475a-b669-7c6649156f64%7D/"
								},
								"avatar": {
									"href": "https://secure.gravatar.com/avatar/24215d1e57db99c8de7178e8d6f8e29d?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FLD-5.png"
								}
							},
							"nickname": "Username",
							"type": "user",
							"account_id": "5e67b3292a0bb00ce033d413"
						}
					},
					"parents": [],
					"date": "2020-12-22T18:32:19+00:00",
					"message": "Initial commit",
					"type": "commit"
				}
			}
		],
		"page": 1,
		"size": 1
	}`

	PayloadListLastTags = `{
		"pagelen": 10,
		"values": [
			{
				"name": "0.0.1",
				"links": {
					"commits": {
						"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commits/0.0.1"
					},
					"self": {
						"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/refs/tags/0.0.1"
					},
					"html": {
						"href": "https://bitbucket.org/username/ritchie-formulas/commits/tag/0.0.1"
					}
				},
				"tagger": {
					"raw": "Username <user@user.com>",
					"type": "author",
					"user": {
						"display_name": "Username",
						"uuid": "{b796025c-834e-475a-b669-7c6649156f64}",
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/users/%7Bb796025c-834e-475a-b669-7c6649156f64%7D"
							},
							"html": {
								"href": "https://bitbucket.org/%7Bb796025c-834e-475a-b669-7c6649156f64%7D/"
							},
							"avatar": {
								"href": "https://secure.gravatar.com/avatar/24215d1e57db99c8de7178e8d6f8e29d?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FLD-5.png"
							}
						},
						"nickname": "Username",
						"type": "user",
						"account_id": "5e67b3292a0bb00ce033d413"
					}
				},
				"date": "2020-12-22T18:33:52+00:00",
				"message": "Added tag 0.0.1 for changeset 224421a2680b",
				"type": "tag",
				"target": {
					"hash": "224421a2680b648a02a9feca6ef025fddadcaeff",
					"repository": {
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas"
							},
							"html": {
								"href": "https://bitbucket.org/username/ritchie-formulas"
							},
							"avatar": {
								"href": "https://bytebucket.org/ravatar/%7Bb2004d9e-f4b2-4784-99cf-b9c7d4322a5f%7D?ts=default"
							}
						},
						"type": "repository",
						"name": "ritchie-formulas",
						"full_name": "username/ritchie-formulas",
						"uuid": "{b2004d9e-f4b2-4784-99cf-b9c7d4322a5f}"
					},
					"links": {
						"self": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"comments": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff/comments"
						},
						"patch": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/patch/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"html": {
							"href": "https://bitbucket.org/username/ritchie-formulas/commits/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"diff": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/diff/224421a2680b648a02a9feca6ef025fddadcaeff"
						},
						"approve": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff/approve"
						},
						"statuses": {
							"href": "https://api.bitbucket.org/2.0/repositories/username/ritchie-formulas/commit/224421a2680b648a02a9feca6ef025fddadcaeff/statuses"
						}
					},
					"author": {
						"raw": "Username <user@user.com>",
						"type": "author",
						"user": {
							"display_name": "Username",
							"uuid": "{b796025c-834e-475a-b669-7c6649156f64}",
							"links": {
								"self": {
									"href": "https://api.bitbucket.org/2.0/users/%7Bb796025c-834e-475a-b669-7c6649156f64%7D"
								},
								"html": {
									"href": "https://bitbucket.org/%7Bb796025c-834e-475a-b669-7c6649156f64%7D/"
								},
								"avatar": {
									"href": "https://secure.gravatar.com/avatar/24215d1e57db99c8de7178e8d6f8e29d?d=https%3A%2F%2Favatar-management--avatars.us-west-2.prod.public.atl-paas.net%2Finitials%2FLD-5.png"
								}
							},
							"nickname": "Username",
							"type": "user",
							"account_id": "5e67b3292a0bb00ce033d413"
						}
					},
					"parents": [],
					"date": "2020-12-22T18:32:19+00:00",
					"message": "Initial commit",
					"type": "commit"
				}
			}
		],
		"page": 1,
		"size": 1
	}`
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
