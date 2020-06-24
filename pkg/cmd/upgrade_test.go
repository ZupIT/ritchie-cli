package cmd

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
)

type stubUpgradeManager struct {
	run func(upgradeUrl string) error
}

func (m stubUpgradeManager) Run(upgradeUrl string) error {
	return m.run(upgradeUrl)
}

func TestUpgradeCmd_runFunc(t *testing.T) {
	type fields struct {
		upgradeUrl string
		Manager    upgrade.Manager
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Run with success",
			fields: fields{
				upgradeUrl: "any url",
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should return err",
			fields: fields{
				upgradeUrl: "any url",
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return errors.New("some error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpgradeCmd(tt.fields.upgradeUrl, tt.fields.Manager)
			if err := u.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("runFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
