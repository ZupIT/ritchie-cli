
package scaffold

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/ZupIT/ritchie-cli/functional"
)

func TestRitSingleScaffold(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite Scaffold")
}

var _ = Describe("RitScaffold", func() {
	BeforeSuite(func() {
		functional.RitSingleInit()
	})

	scenariosScaffold := functional.LoadScenarios("scaffold_feature.json")

	DescribeTable("When running core command",
		func(scenario functional.Scenario) {
			out, err := scenario.RunSteps()
			Expect(err).To(Succeed())
			Expect(out).To(ContainSubstring(scenario.Result))
		},
		Entry("Run scaffold coffee-go", scenariosScaffold[0]),
	)
})
