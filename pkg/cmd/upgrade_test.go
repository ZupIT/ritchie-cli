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
	stableVersion func() (string, error)
	updateCache func() error
}

func (vr stubVersionResolver) StableVersion() (string, error) {
	return vr.stableVersion()
}

func (vr stubVersionResolver) UpdateCache() error {
	return vr.updateCache()
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
					func() (string, error) {
						return "1.0.0", nil
					},
					func() error {
						return nil
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
			name: "Should return err on UpdateCache",
			fields: fields{
				resolver: stubVersionResolver{
					func() (string, error) {
						return "", nil
					},
					func() error {
						return errors.New("some error")
					},
				},
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
		{
			name: "Should return err on Run",
			fields: fields{
				resolver: stubVersionResolver{
					func() (string, error) {
						return "", nil
					},
					func() error {
						return nil
					},
				},
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
