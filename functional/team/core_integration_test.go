package team

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

var _ = Describe("RitCore", func() {
	scenariosCore := functional.LoadScenarios("core_feature.json")

	DescribeTable("When running core command",
		func(scenario functional.Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},

		Entry("Set", scenariosCore[0]),
		Entry("Add", scenariosCore[1]),
		Entry("Completion", scenariosCore[2]),
		Entry("Clean", scenariosCore[3]),
		Entry("Create", scenariosCore[4]),
		Entry("Delete", scenariosCore[5]),
		Entry("List", scenariosCore[6]),
		Entry("Show", scenariosCore[7]),
		Entry("Update", scenariosCore[8]),
	)
})
