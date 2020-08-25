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

package upgrade

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type stubResolver struct {
	stableVersion func() (string, error)
	updateCache   func() error
	verifyNewVersion func(current, installed string) string
}

func (r stubResolver) StableVersion() (string, error) {
	return r.stableVersion()
}

func (r stubResolver) UpdateCache() error {
	return r.updateCache()
}

func (r stubResolver) VerifyNewVersion(current, installed string) string {
	return r.verifyNewVersion(current, installed)
}

func TestUpgradeUrl(t *testing.T) {
	type args struct {
		resolver version.Resolver
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get url with success",
			args: args{
				resolver: stubResolver{
					stableVersion: func() (string, error) {
						return "1.0.0", nil
					},
				},
			},
			want: func() string {
				expected := fmt.Sprintf(upgradeUrlFormat, "1.0.0", runtime.GOOS)
				if runtime.GOOS == "windows" {
					expected += ".exe"
				}
				return expected
			}(),
		},
		{
			name: "Get url for when happening a error",
			args: args{
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
			duf := DefaultUrlFinder{}
			if got := duf.Url(tt.args.resolver); got != tt.want {
				t.Errorf("UpgradeUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
