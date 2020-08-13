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

package stdin

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

func TestRitSingleStdin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite Stdin")
}

var _ = Describe("RitStdin", func() {
	BeforeSuite(func() {
		functional.RitSingleInit()
	})

	scenariosStdin := functional.LoadScenarios("stdin_feature.json")

	DescribeTable("When running core command",
		func(scenario functional.Scenario) {
			out, err := scenario.RunStdin()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},

		Entry("Set context STDIN", scenariosStdin[0]),
		Entry("Delete context STDIN", scenariosStdin[1]),
		Entry("Add new repo STDIN", scenariosStdin[2]),
		Entry("Delete repo STDIN", scenariosStdin[3]),
		Entry("Set credentials STDIN", scenariosStdin[4]),
	)

})
