package github

import (
	"reflect"
	"testing"
)

func TestLatestTagUrl(t *testing.T) {
	type in struct {
		Url   string
		Token string
	}
	tests := []struct {
		name string
		in   in
		want string
	}{
		{
			name: "Generate LatestTagUrlWithSuccess",
			in: in{
				Url: "http://github.com/zupIt/ritchie-cli",
			},
			want: "https://api.github.com/repos/zupIt/ritchie-cli/releases/latest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewRepoInfo(tt.in.Url, tt.in.Token)

			if got := in.LatestTagUrl(); got != tt.want {
				t.Errorf("LatestTagUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagsUrl(t *testing.T) {
	type in struct {
		Url   string
		Token string
	}
	tests := []struct {
		name string
		in   in
		want string
	}{
		{
			name: "Generate LatestTagUrlWithSuccess",
			in: in{
				Url: "http://github.com/zupIt/ritchie-cli",
			},
			want: "https://api.github.com/repos/zupIt/ritchie-cli/releases",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewRepoInfo(tt.in.Url, tt.in.Token)

			if got := in.TagsUrl(); got != tt.want {
				t.Errorf("TagsUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoInfo_TokenHeader(t *testing.T) {
	type in struct {
		Url   string
		Token string
	}
	tests := []struct {
		name string
		in   in
		want string
	}{
		{
			name: "Generate LatestTagUrlWithSuccess",
			in: in{
				Url:   "http://github.com/zupIt/ritchie-cli",
				Token: "any_token",
			},
			want: "token any_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewRepoInfo(tt.in.Url, tt.in.Token)

			if got := in.TokenHeader(); got != tt.want {
				t.Errorf("TokenHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipUrl(t *testing.T) {
	type in struct {
		Url     string
		Token   string
		version string
	}

	tests := []struct {
		name string
		in   in
		want string
	}{
		{
			name: "Generate LatestTagUrlWithSuccess",
			in: in{
				Url:     "http://github.com/zupIt/ritchie-cli",
				version: "0.0.3",
			},
			want: "https://api.github.com/repos/zupIt/ritchie-cli/zipball/0.0.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewRepoInfo(tt.in.Url, tt.in.Token)
			if got := in.ZipUrl(tt.in.version); got != tt.want {
				t.Errorf("ZipUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTags_Names(t *testing.T) {
	tests := []struct {
		name string
		t    Tags
		want []string
	}{
		{
			name: "Return tags name",
			t: Tags{
				{
					Name: "tag1",
				},
				{
					Name: "tag2",
				},
			},
			want: []string{"tag1", "tag2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Names(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Names() = %v, want %v", got, tt.want)
			}
		})
	}
}
