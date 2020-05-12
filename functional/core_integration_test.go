package functional

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("RitCore", func() {
	scenariosCore := LoadScenarios("core_feature.json")

	DescribeTable("When running core command",
		func(scenario Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},
		Entry("When running show context", scenariosCore[0]),
		Entry("When running set context", scenariosCore[1]),
		Entry("Add new repo without error", scenariosCore[2]),
		Entry("Clean repo without error", scenariosCore[3]),
		Entry("List repo without error", scenariosCore[4]),
		Entry("Delete repo without error", scenariosCore[5]),
		Entry("Set Credential", scenariosCore[6]),
	)

})
