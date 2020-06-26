package cmd

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type stubUpgradeManager struct {
	run func(upgradeUrl string) error
}

func (m stubUpgradeManager) Run(upgradeUrl string) error {
	return m.run(upgradeUrl)
}

type stubUrlFinder struct {
	url func(edition api.Edition, resolver version.Resolver) string
}

func (uf stubUrlFinder) Url(edition api.Edition, resolver version.Resolver) string {
	return uf.url(edition, resolver)
}

type stubVersionResolver struct {
	stableVersion func(fromCache bool) (string, error)
}

func (vr stubVersionResolver) StableVersion(fromCache bool) (string, error) {
	return vr.stableVersion(fromCache)
}

func TestUpgradeCmd_runFunc(t *testing.T) {
	type fields struct {
		edition   api.Edition
		resolver  version.Resolver
		Manager   upgrade.Manager
		UrlFinder upgrade.UrlFinder
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run with success",
			fields: fields{
				edition: "tingle",
				resolver: stubVersionResolver{
					func(fromCache bool) (string, error) {
						return "1.0.0", nil
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return nil
					},
				},
				UrlFinder: stubUrlFinder{
					func(edition api.Edition, resolver version.Resolver) string {
						return "any url"
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should return err",
			fields: fields{
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return errors.New("some error")
					},
				},
				UrlFinder: stubUrlFinder{
					func(edition api.Edition, resolver version.Resolver) string {
						return "any url"
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpgradeCmd(tt.fields.edition, tt.fields.resolver, tt.fields.Manager, tt.fields.UrlFinder)
			if err := u.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("runFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
