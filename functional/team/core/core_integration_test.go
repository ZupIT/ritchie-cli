package core

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

func TestRitTeam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite")
}

var _ = Describe("RitCore", func() {
	BeforeSuite(func() {
		functional.RitTeamInit()
	})

	scenariosCore := functional.LoadScenarios("core_feature.json")

	DescribeTable("When running core command",
		func(scenario functional.Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},

		Entry("Set", scenariosCore[0]),
		// Entry("Set Credential", scenariosCore[1]),
		// Entry("Set Context", scenariosCore[2]),
		Entry("Create", scenariosCore[3]),
		Entry("Delete", scenariosCore[4]),
		// Entry("Delete Context", scenariosCore[5]),
		Entry("Show", scenariosCore[6]),
		Entry("Show Context", scenariosCore[7]),
		Entry("Help", scenariosCore[8]),
		Entry("Completion", scenariosCore[9]),
		Entry("Completion bash", scenariosCore[10]),
		Entry("Completion zsh", scenariosCore[11]),
		Entry("Version", scenariosCore[12]),
		Entry("Add", scenariosCore[13]),
		// Entry("Add repo", scenariosCore[14]),
		Entry("List", scenariosCore[15]),
		Entry("List repo", scenariosCore[16]),
		Entry("Update", scenariosCore[17]),
		// Entry("Update repo", scenariosCore[18]),
		// Entry("delete repo", scenariosCore[21]),
	)
})
