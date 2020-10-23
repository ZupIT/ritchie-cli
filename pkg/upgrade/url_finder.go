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
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type UrlFinder interface {
	Url(os string) string
}

type DefaultUrlFinder struct {
	version version.Resolver
}

func NewDefaultUrlFinder(version version.Resolver) DefaultUrlFinder {
	return DefaultUrlFinder{version: version}
}

func (duf DefaultUrlFinder) Url(os string) string { //nolint:stylecheck
	//stableVersion, err := resolver.StableVersion()
	stableVersion, err := duf.version.StableVersion()

	if err != nil {
		return ""
	}

	upgradeURL := fmt.Sprintf(upgradeUrlFormat, stableVersion, os)

	if os == osutil.Windows {
		upgradeURL += ".exe"
	}

	return upgradeURL
}
