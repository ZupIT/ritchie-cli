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
		Entry("Add", scenariosCore[2]),
		Entry("Completion", scenariosCore[4]),
		Entry("Clean", scenariosCore[5]),
		Entry("Create", scenariosCore[7]),
		Entry("Delete", scenariosCore[8]),
		Entry("List", scenariosCore[9]),
		Entry("Show", scenariosCore[10]),
		Entry("Update", scenariosCore[11]),

		Entry("Set Server", scenariosCore[12]),
	)
})
