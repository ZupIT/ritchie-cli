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
	"github.com/ZupIT/ritchie-cli/pkg/version/versionutil"
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

var stubUpgradeUtilApplyExecutions = 0

type StubUpgradeUtil struct {
	apply func(reader io.Reader, opts update.Options) error
}

func (u StubUpgradeUtil) Apply(reader io.Reader, opts update.Options) error {
	stubUpgradeUtilApplyExecutions++
	return u.apply(reader, opts)
}

func TestUpgradeCmd(t *testing.T) {
	type fields struct {
		upgradeUrl  string
		upgradeUtil UpgradeUtil
	}
	tests := []struct {
		name                         string
		fields                       fields
		executionWantedOfUpgradeUtil int
		wantErr                      bool
	}{
		{
			name: "Run with success",
			fields: fields{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
				upgradeUtil: StubUpgradeUtil{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgradeUtil: 1,
			wantErr:                      false,
		},
		{
			name: "Should return err when url is empty",
			fields: fields{
				upgradeUrl: "",
				upgradeUtil: StubUpgradeUtil{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgradeUtil: 0,
			wantErr:                      true,
		},
		{
			name: "Should return err when happening err when perform get",
			fields: fields{
				upgradeUrl: "some url",
				upgradeUtil: StubUpgradeUtil{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgradeUtil: 0,
			wantErr:                      true,
		},
		{
			name: "Should return err when get return not 200 code",
			fields: fields{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
				})).URL,
				upgradeUtil: StubUpgradeUtil{
					apply: func(reader io.Reader, opts update.Options) error {
						return nil
					},
				},
			},
			executionWantedOfUpgradeUtil: 0,
			wantErr:                      true,
		},
		{
			name: "Should return err when fail to apply",
			fields: fields{
				upgradeUrl: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).URL,
				upgradeUtil: StubUpgradeUtil{
					apply: func(reader io.Reader, opts update.Options) error {
						return errors.New("some error")
					},
				},
			},
			executionWantedOfUpgradeUtil: 1,
			wantErr:                      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubUpgradeUtilApplyExecutions = 0
			err := NewUpgradeCmd(tt.fields.upgradeUrl, tt.fields.upgradeUtil).Execute()
			if stubUpgradeUtilApplyExecutions != tt.executionWantedOfUpgradeUtil {
				t.Errorf("Expected %d executions of StubUpgradeUtil.apply", tt.executionWantedOfUpgradeUtil)
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
		resolver versionutil.Resolver
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
