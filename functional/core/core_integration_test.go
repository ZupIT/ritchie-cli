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

package core

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

func TestRitSingleCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite Core")
}

var _ = Describe("RitCore", func() {
	BeforeSuite(func() {
		functional.RitSingleInit()
	})

	scenariosCore := functional.LoadScenarios("core_feature.json")

	DescribeTable("When running core command",
		func(scenario functional.Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},

		Entry("Show context", scenariosCore[0]),
		// Entry("Set context", scenariosCore[1]),
		// Entry("Delete context", scenariosCore[2]),

		Entry("Add", scenariosCore[3]),
		// Entry("Add new repo", scenariosCore[4]),

		Entry("List", scenariosCore[5]),
		Entry("List repo", scenariosCore[6]),
		// Entry("Delete repo", scenariosCore[7]),

		// Entry("Set Credential", scenariosCore[8]),
		Entry("Set", scenariosCore[9]),
		// Entry("Update repo", scenariosCore[10]),

		Entry("Completion", scenariosCore[11]),
		Entry("Completion bash", scenariosCore[12]),
		Entry("Completion zsh", scenariosCore[13]),

		Entry("Version", scenariosCore[14]),
		Entry("Delete", scenariosCore[15]),
		Entry("Show", scenariosCore[16]),
		Entry("Help", scenariosCore[17]),
	)

})
