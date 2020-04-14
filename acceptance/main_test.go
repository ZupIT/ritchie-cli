package main

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestRit(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Rit Suite")
}

var _ = ginkgo.Describe("Rit", func() {
	var (
		coffeeInputs, cmdSetCtx, cmdSetCred, cmdCreateUsr, cmdDelUsr,cmdDelCtx, cmdShowCtx Inputs
	)
	ginkgo.BeforeEach(func() {
		ginkgo.GinkgoWriter.Write([]byte("============================== Setting inputs ==============================\n"))
		coffeeInputs = Inputs{
			args:   []string{"scaffold", "generate", "coffee-go"},
			prompt: []string{"Ritchie\n", "j", "j", "\n", "\n", "\n"},
		}
		cmdSetCtx = Inputs{
			args:   []string{"set", "context"},
			prompt: []string{"\n"},
		}
		cmdSetCred = Inputs{
			args:   []string{"set", "credential"},
			prompt: []string{"\n"},
		}
		cmdCreateUsr = Inputs{
			args:   []string{"create", "user"},
			prompt: []string{"\n"},
		}
		cmdDelUsr = Inputs{
			args:   []string{"delete", "user"},
			prompt: []string{"\n"},
		}
		cmdDelCtx = Inputs{
			args:   []string{"delete", "user"},
			prompt: []string{"\n"},
		}
		cmdShowCtx = Inputs{
			args:   []string{"show", "context"},
			prompt: []string{"\n"},
		}
	})

	ginkgo.Describe("Choosing and setting inputs using formulas", func() {
		ginkgo.Context("When running coffee formula", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := coffeeInputs.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Your macchiato coffee is ready, have a seat and enjoy your drink"))
			})
		})
		ginkgo.Context("When running core commands[set credential]", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := cmdSetCred.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
			})
		})
		ginkgo.Context("When running core commands[set context]", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := cmdSetCtx.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
			})
		})
		ginkgo.Context("When running core commands[create user]", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := cmdCreateUsr.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
			})
		})
		ginkgo.Context("When running core commands[delete context]", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := cmdCreateUsr.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
			})
		})
		ginkgo.Context("When running core commands[delete user]", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := cmdDelUsr.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
			})
		})
		ginkgo.Context("When running core commands[show context]", func() {
			ginkgo.It("Should run without errors", func() {
				err, out := cmdShowCtx.RunRit()
				fmt.Println(out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Current context:"))
			})
		})
	})
})
