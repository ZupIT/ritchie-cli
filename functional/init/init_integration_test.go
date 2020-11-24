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

package init

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

func TestRitSingleInit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite Init")
}

var _ = Describe("RitSingleInit", func() {
	BeforeEach(func() {
		functional.RitClearConfigs()
	})

	scenariosCore := functional.LoadScenarios("init_feature.json")

	DescribeTable("When running command without init",
		func(scenario functional.Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},

		Entry("Show env", scenariosCore[0]),
		Entry("Set env", scenariosCore[1]),
		Entry("Delete context", scenariosCore[2]),
		Entry("List repo", scenariosCore[3]),
		Entry("Delete repo", scenariosCore[4]),
		Entry("Create formula", scenariosCore[5]),
		Entry("Do init", scenariosCore[6]),
	)

})
