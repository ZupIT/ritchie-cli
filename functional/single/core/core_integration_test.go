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
		functional.RitInit()
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
		// Entry("Clean repo", scenariosCore[5]),
		Entry("List repo", scenariosCore[6]),
		// Entry("Delete repo", scenariosCore[7]),

		Entry("Set", scenariosCore[8]),
		// Entry("Set Credential", scenariosCore[9]),

		Entry("Completion", scenariosCore[10]),
		Entry("Completion bash", scenariosCore[11]),
		Entry("Completion zsh", scenariosCore[12]),

		Entry("Version", scenariosCore[13]),
		Entry("Delete", scenariosCore[14]),
		Entry("Show", scenariosCore[15]),
		Entry("Help", scenariosCore[16]),
	)

})
