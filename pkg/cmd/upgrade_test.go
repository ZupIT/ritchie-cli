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

package cmd

import (
	"errors"
	"testing"

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
	url func(resolver version.Resolver) string
}

func (uf stubUrlFinder) Url(resolver version.Resolver) string {
	return uf.url(resolver)
}

type stubVersionResolver struct {
	stableVersion    func() (string, error)
	updateCache      func() error
	verifyNewVersion func(current, installed string) string
}

func (vr stubVersionResolver) StableVersion() (string, error) {
	return vr.stableVersion()
}

func (vr stubVersionResolver) UpdateCache() error {
	return vr.updateCache()
}

func (vr stubVersionResolver) VerifyNewVersion(current, installed string) string {
	return vr.VerifyNewVersion(current, installed)
}

func TestUpgradeCmd_runFunc(t *testing.T) {
	type fields struct {
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
				resolver: stubVersionResolver{
					func() (string, error) {
						return "1.0.0", nil
					},
					func() error {
						return nil
					},
					func(current, installed string) string {
						return ""
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return nil
					},
				},
				UrlFinder: stubUrlFinder{
					func(resolver version.Resolver) string {
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
					func(current, installed string) string {
						return ""
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return errors.New("some error")
					},
				},
				UrlFinder: stubUrlFinder{
					func(resolver version.Resolver) string {
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
					func(current, installed string) string {
						return ""
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return errors.New("some error")
					},
				},
				UrlFinder: stubUrlFinder{
					func(resolver version.Resolver) string {
						return "any url"
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpgradeCmd(tt.fields.resolver, tt.fields.Manager, tt.fields.UrlFinder)
			if err := u.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("runFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
