package single

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

var _ = Describe("RitStdin", func() {
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
		Entry("Clean repo STDIN", scenariosStdin[3]),
		Entry("Delete repo STDIN", scenariosStdin[4]),
		Entry("Set credentials STDIN", scenariosStdin[5]),

	)

})
