
package functional

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("RitScaffold", func() {
	scenariosScaffold := LoadScenarios("scaffold_feature.json")

	DescribeTable("Running entry for Scaffolds",
		func(scenario Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},
		Entry("Run scaffold coffee-go", scenariosScaffold[0]),
	)
})
