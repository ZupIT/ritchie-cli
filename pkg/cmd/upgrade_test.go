package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/inconshreveable/go-update"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version/version_util"
)

type stubResolver struct {
	getCurrentVersion func() (string, error)
	getStableVersion  func() (string, error)
}

func (r stubResolver) GetCurrentVersion() (string, error) {
	return r.getCurrentVersion()
}

func (r stubResolver) GetStableVersion() (string, error) {
	return r.getStableVersion()
}

var stubUpgradeApplyExecutions = 0

type StubUpgrade struct {
	apply func(reader io.Reader, opts update.Options) error
}

func (u StubUpgrade) Apply(reader io.Reader, opts update.Options) error {
	stubUpgradeApplyExecutions++
	return u.apply(reader, opts)
}

func TestUpgradeCmd(t *testing.T) {
	type fields struct {
		upgradeUrl string
		upgrade    upgrade.Upgrade
	}
	tests := []struct {
		name                     string
		fields                   fields
		executionWantedOfUpgrade int
		wantErr                  bool
	}{
		{
			name: "Run with success",
			fields: fields{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
				upgrade: StubUpgrade{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgrade: 1,
			wantErr:                  false,
		},
		{
			name: "Should return err when url is empty",
			fields: fields{
				upgradeUrl: "",
				upgrade: StubUpgrade{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgrade: 0,
			wantErr:                  true,
		},
		{
			name: "Should return err when happening err when perform get",
			fields: fields{
				upgradeUrl: "some url",
				upgrade: StubUpgrade{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgrade: 0,
			wantErr:                  true,
		},
		{
			name: "Should return err when get return not 200 code",
			fields: fields{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
				})).URL,
				upgrade: StubUpgrade{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgrade: 0,
			wantErr:                  true,
		},
		{
			name: "Should return err when fail to apply",
			fields: fields{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
				upgrade: StubUpgrade{
					apply: func(reader io.Reader, opts update.Options) error {
						return errors.New("some error")
					},
				},
			},
			executionWantedOfUpgrade: 1,
			wantErr:                  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubUpgradeApplyExecutions = 0
			err := NewUpgradeCmd(tt.fields.upgradeUrl, tt.fields.upgrade).Execute()
			if stubUpgradeApplyExecutions != tt.executionWantedOfUpgrade {
				t.Errorf("Expected %d executions of StubUpgrade.apply", tt.executionWantedOfUpgrade)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("UpgradeCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpgradeUrl(t *testing.T) {
	type args struct {
		edition  api.Edition
		resolver version_util.Resolver
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
					getStableVersion: func() (string, error) {
						return "1.0.0", nil
					},
					getCurrentVersion: func() (string, error) {
						return "1.1.0", nil
					},
				},
			},
			want: func() string {
				expected := fmt.Sprintf(upgradeUrlFormat, "1.0.0", runtime.GOOS, api.Single)
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
					getStableVersion: func() (string, error) {
						return "1.0.0", nil
					},
					getCurrentVersion: func() (string, error) {
						return "1.1.0", nil
					},
				},
			},
			want: func() string {
				expected := fmt.Sprintf(upgradeUrlFormat, "1.0.0", runtime.GOOS, api.Team)
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
					getStableVersion: func() (string, error) {
						return "1.0.0", errors.New("some error")
					},
					getCurrentVersion: func() (string, error) {
						return "1.1.0", nil
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
