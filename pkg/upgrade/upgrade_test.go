package upgrade

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/resource"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type stubResolver struct {
	stableVersion func() (string, error)
}

func (r stubResolver) StableVersion() (string, error) {
	return r.stableVersion()
}

func TestUpgradeUrl(t *testing.T) {
	type args struct {
		edition  api.Edition
		resolver version.Resolver
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get url for single edition",
			args: args{
				edition: api.Single,
				resolver: stubResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", nil
					},
				},
			},
			want: func() string {
				expected := fmt.Sprintf(resource.UpgradeUrlFormat, "1.0.0", runtime.GOOS, api.Single)
				if runtime.GOOS == "windows" {
					expected += ".exe"
				}
				return expected
			}(),
		},
		{
			name: "Get url for team edition",
			args: args{
				edition: api.Team,
				resolver: stubResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", nil
					},
				},
			},
			want: func() string {
				expected := fmt.Sprintf(resource.UpgradeUrlFormat, "1.0.0", runtime.GOOS, api.Team)
				if runtime.GOOS == "windows" {
					expected += ".exe"
				}
				return expected
			}(),
		},
		{
			name: "Get url for when happening a error",
			args: args{
				edition: api.Team,
				resolver: stubResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", errors.New("some error")
					},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpgradeUrl(tt.args.edition, tt.args.resolver); got != tt.want {
				t.Errorf("UpgradeUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
