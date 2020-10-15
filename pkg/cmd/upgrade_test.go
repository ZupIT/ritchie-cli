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

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
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
	url func() string
}

func (uf stubUrlFinder) Url(os string) string {
	return uf.url()
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
	return vr.verifyNewVersion(current, installed)
}

func TestUpgradeCmd_runFunc(t *testing.T) {
	type in struct {
		resolver  version.Resolver
		Manager   upgrade.Manager
		UrlFinder upgrade.UrlFinder
		input     prompt.InputList
		file      stream.FileWriteReadExister
	}
	tests := []struct {
		name    string
		in      in
		wantErr bool
	}{
		{
			name: "Run with success",
			in: in{
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
					func() string {
						return "any url"
					},
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should return err on UpdateCache",
			in: in{
				resolver: stubVersionResolver{
					func() (string, error) {
						return "", nil
					},
					func() error {
						return errors.New("update cache error")
					},
					func(current, installed string) string {
						return ""
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return errors.New("upgrade url error")
					},
				},
				UrlFinder: stubUrlFinder{
					func() string {
						return "any url"
					},
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should return err on Run",
			in: in{
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
						return errors.New("upgrade url error")
					},
				},
				UrlFinder: stubUrlFinder{
					func() string {
						return "any url"
					},
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			wantErr: true,
		},
		{
			name: "success with no metrics file",
			in: in{
				resolver: stubVersionResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", nil
					},
					updateCache: func() error {
						return nil
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return nil
					},
				},
				UrlFinder: stubUrlFinder{
					url: func() string {
						return "any url"
					},
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
				input: inputListCustomMock{func(name string, items []string, defaultValue string) (string, error) {
					return DoNotAcceptMetrics, nil
				}},
			},
			wantErr: false,
		},
		{
			name: "fail on list with no metrics file",
			in: in{
				resolver: stubVersionResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", nil
					},
					updateCache: func() error {
						return nil
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return nil
					},
				},
				UrlFinder: stubUrlFinder{
					func() string {
						return "any url"
					},
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
				},
				input: inputListErrorMock{},
			},
			wantErr: true,
		},
		{
			name: "fail on write with no metrics file",
			in: in{
				resolver: stubVersionResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", nil
					},
					updateCache: func() error {
						return nil
					},
				},
				Manager: stubUpgradeManager{
					func(upgradeUrl string) error {
						return nil
					},
				},
				UrlFinder: stubUrlFinder{
					func() string {
						return "any url"
					},
				},
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					WriteMock: func(path string, content []byte) error {
						return errors.New("error writing file")
					},
				},
				input: inputListMock{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpgradeCmd(
				tt.in.resolver,
				tt.in.Manager,
				tt.in.UrlFinder,
				tt.in.input,
				tt.in.file)
			if err := u.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("runFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
