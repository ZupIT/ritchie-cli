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
		Entry("Delete repo STDIN", scenariosStdin[4]),
		Entry("Set credentials STDIN", scenariosStdin[5]),

	)

})
