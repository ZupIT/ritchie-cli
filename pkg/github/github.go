package github

import (
	"fmt"
	"io"
	"strings"
)

const (
	ZipUrlPattern       = "https://api.github.com/repos/%s/%s/zipball/%s"
	TagsUrlPattern      = "https://api.github.com/repos/%s/%s/releases"
	LatestTagUrlPattern = "https://api.github.com/repos/%s/%s/releases/latest"
)

type Tag struct {
	Name string `json:"tag_name"`
}

type Tags []Tag

func (t Tags) Names() []string {
	var tags []string
	for i := range t {
		tags = append(tags, t[i].Name)
	}

	return tags
}

type RepoInfo interface {
	ZipUrl(version string) string
	TagsUrl() string
	LatestTagUrl() string
	TokenHeader() string
	Token() string
}

type DefaultRepoInfo struct {
	owner string
	repo  string
	token string
}

type Repositories interface {
	Zipball(info RepoInfo, version string) (io.ReadCloser, error)
	Tags(info RepoInfo) (Tags, error)
	LatestTag(info RepoInfo) (Tag, error)
}

// NewRepoInfo returns the RepoInfo built by repository url
// Repository url e.g. https://github.com/{{owner}}/{{repo}}
func NewRepoInfo(url string, token string) RepoInfo {
	split := strings.Split(url, "/")
	repo := split[len(split)-1]
	owner := split[len(split)-2]

	return DefaultRepoInfo{
		owner: owner,
		repo:  repo,
		token: token,
	}
}

// ZipUrl returns the GitHub API URL for download zipball repository
// e.g. https://api.github.com/repos/{{owner}}/{{repo}}/zipball/{{tag-version}}
func (in DefaultRepoInfo) ZipUrl(version string) string {
	return fmt.Sprintf(ZipUrlPattern, in.owner, in.repo, version)
}

// TagsUrl returns the GitHub API URL for get all tags
// e.g. https://api.github.com/repos/{{owner}}/{{repo}}/tags
func (in DefaultRepoInfo) TagsUrl() string {
	return fmt.Sprintf(TagsUrlPattern, in.owner, in.repo)
}

// LatestTagUrl returns the GitHub API URL for get latest tag release
// https://api.github.com/repos/:owner/:repo/releases/latest
func (in DefaultRepoInfo) LatestTagUrl() string {
	return fmt.Sprintf(LatestTagUrlPattern, in.owner, in.repo)
}

// TokenHeader returns the Authorization value formatted for Github API integration
// e.g. "token f39c5aca-858f-4a04-9ca3-5104d02b9c56"
func (in DefaultRepoInfo) TokenHeader() string {
	return fmt.Sprintf("token %s", in.token)
}

func (in DefaultRepoInfo) Token() string {
	return in.token
}
