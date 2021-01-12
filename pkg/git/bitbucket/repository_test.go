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
					Name:        "0.0.1",
					Description: "",
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
			wantErr: "There was an error adding the repository, status: 400 - Bad Request.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("TokenHeader").Return(tt.info.tokenHeader)
			git.On("TagsUrl").Return(tt.info.tagsUrl)
			re := NewRepoManager(tt.client)
			got, err := re.Tags(git)
			if err == nil {
				assert.Equal(t, tt.want, got)
			} else {
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
			want:    "400 Bad Request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("TokenHeader").Return(tt.info.tokenHeader)
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
			want: git.Tag{Name: "0.0.1", Description: ""},
		},
		{
			name:   "Return err when request fail",
			client: mockServerThatFail.Client(),
			info: info{
				latestTagUrl: mockServerThatFail.URL,
			},
			want: git.Tag{},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &mocks.RepoInfo{}
			git.On("Token").Return(tt.info.token)
			git.On("TokenHeader").Return(tt.info.tokenHeader)
			git.On("LatestTagUrl").Return(tt.info.latestTagUrl)
			re := NewRepoManager(tt.client)
			got, err := re.LatestTag(git)
			if err == nil {
				assert.Equal(t, tt.want, got)
			} else {
				assert.EqualError(t, err, tt.wantErr)
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

	PayloadEmptyTags = `{
		"pagelen": 10,
		"values": [],
		"page": 1,
		"size": 0
	}`
)

type info struct {
	zipUrl       string
	tagsUrl      string
	latestTagUrl string
	tokenHeader  string
	token        string
}
