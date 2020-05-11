
package functional

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestRitScaffold(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite")
}

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
