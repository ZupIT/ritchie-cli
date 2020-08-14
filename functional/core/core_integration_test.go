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

		Entry(scenariosCore[0].Entry, scenariosCore[0]),
		// Entry(scenariosCore[1].Entry, scenariosCore[1]),
		// Entry(scenariosCore[2].Entry, scenariosCore[2]),

		Entry(scenariosCore[3].Entry, scenariosCore[3]),
		// Entry(scenariosCore[4].Entry, scenariosCore[4]),

		Entry(scenariosCore[5].Entry, scenariosCore[5]),
		Entry(scenariosCore[6].Entry, scenariosCore[6]),
		// Entry(scenariosCore[7].Entry, scenariosCore[7]),

		// Entry(scenariosCore[8].Entry, scenariosCore[8]),
		Entry(scenariosCore[9].Entry, scenariosCore[9]),
		// Entry(scenariosCore[10].Entry, scenariosCore[10]),

		Entry(scenariosCore[11].Entry, scenariosCore[11]),
		Entry(scenariosCore[12].Entry, scenariosCore[12]),
		Entry(scenariosCore[13].Entry, scenariosCore[13]),

		Entry(scenariosCore[14].Entry, scenariosCore[14]),
		Entry(scenariosCore[15].Entry, scenariosCore[15]),
		Entry(scenariosCore[16].Entry, scenariosCore[16]),
		Entry(scenariosCore[17].Entry, scenariosCore[17]),
	)
})
