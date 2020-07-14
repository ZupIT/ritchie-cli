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

type RepoInfo struct {
	Owner string
	Repo  string
	Token string
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

	return RepoInfo{
		Owner: owner,
		Repo:  repo,
		Token: token,
	}
}

// ZipUrl returns the GitHub API URL for download zipball repository
// e.g. https://api.github.com/repos/{{owner}}/{{repo}}/zipball/{{tag-version}}
func (in RepoInfo) ZipUrl(version string) string {
	return fmt.Sprintf(ZipUrlPattern, in.Owner, in.Repo, version)
}

// TagsUrl returns the GitHub API URL for get all tags
// e.g. https://api.github.com/repos/{{owner}}/{{repo}}/tags
func (in RepoInfo) TagsUrl() string {
	return fmt.Sprintf(TagsUrlPattern, in.Owner, in.Repo)
}

// LatestTagUrl returns the GitHub API URL for get latest tag release
// https://api.github.com/repos/:owner/:repo/releases/latest
func (in RepoInfo) LatestTagUrl() string {
	return fmt.Sprintf(LatestTagUrlPattern, in.Owner, in.Repo)
}

// TokenHeader returns the Authorization value formatted for Github API integration
// e.g. "token f39c5aca-858f-4a04-9ca3-5104d02b9c56"
func (in RepoInfo) TokenHeader() string {
	return fmt.Sprintf("token %s", in.Token)
}
