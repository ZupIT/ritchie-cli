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
		//TODO ritLogin
		coffeeInputs Inputs
		//cmdSetCtx, cmdSetCred, cmdCreateUsr,
		//cmdDelUsr, cmdDelCtx, cmdShowCtx, cmdHelp,
		//cmdCompBash, cmdCompZsh, cmdVersion, cmdAddRepo,
		//cmdListRepo, cmdUpdateRepo, cmdCleanRepo Inputs
	)
	ginkgo.BeforeEach(func() {
		ginkgo.GinkgoWriter.Write([]byte("============================== Setting inputs ==============================\n"))
		coffeeInputs = Inputs{
			args:   []string{"scaffold", "generate", "coffee-go"},
			prompt: []string{"Dennis.Ritchie\n", "j", "j", "\n", "\n", "\n"},
		}
		//cmdSetCtx = Inputs{
		//	args:   []string{"set", "context"},
		//	prompt: []string{"\n"},
		//}
		//cmdSetCred = Inputs{
		//	args:   []string{"set", "credential"},
		//	prompt: []string{"\n"},
		//}
		//cmdCreateUsr = Inputs{
		//	args:   []string{"create", "user"},
		//	prompt: []string{"\n"},
		//}
		//cmdDelUsr = Inputs{
		//	args:   []string{"delete", "user"},
		//	prompt: []string{"\n"},
		//}
		//cmdDelCtx = Inputs{
		//	args:   []string{"delete", "user"},
		//	prompt: []string{"\n"},
		//}
		//cmdShowCtx = Inputs{
		//	args:   []string{"show", "context"},
		//	prompt: []string{"\n"},
		//}
		//cmdHelp = Inputs{
		//	args:   []string{"--help"},
		//	prompt: []string{"\n"},
		//}
		//cmdCompBash = Inputs{
		//	args:   []string{"completion", "bash"},
		//	prompt: []string{"\n"},
		//}
		//cmdCompZsh = Inputs{
		//	args:   []string{"completion", "zsh"},
		//	prompt: []string{"\n"},
		//}
		//cmdVersion = Inputs{
		//	args:   []string{"--version"},
		//	prompt: []string{"\n"},
		//}
		//cmdAddRepo = Inputs{
		//	args:   []string{"add", "repo"},
		//	prompt: []string{"\n"},
		//}
		//cmdListRepo = Inputs{
		//	args:   []string{"list", "repo"},
		//	prompt: []string{"\n"},
		//}
		//cmdUpdateRepo = Inputs{
		//	args:   []string{"update", "repo"},
		//	prompt: []string{"\n"},
		//}
		//cmdCleanRepo = Inputs{
		//	args:   []string{"clean", "repo"},
		//	prompt: []string{"\n"},
		//}
	})

	ginkgo.Describe("Choosing and setting inputs using formulas", func() {
		ginkgo.Context("When running coffee formula", func() {
			ginkgo.It("Should run without errors", func() {
				out, err := coffeeInputs.RunRit()
				fmt.Println("OUT >>>> ", out)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(out).To(gomega.ContainSubstring("Your macchiato coffee is ready, have a seat and enjoy your drink"))
			})
		})
		//ginkgo.Context("When running core commands[set context]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdSetCtx.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
		//	})
		//})
		//ginkgo.Context("When running core commands[set credential]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdSetCred.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
		//	})
		//})
		//ginkgo.Context("When running core commands[create user]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdCreateUsr.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
		//	})
		//})
		//ginkgo.Context("When running core commands[delete user]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdDelUsr.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
		//	})
		//})
		//ginkgo.Context("When running core commands[delete context]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdDelCtx.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Set context successful!"))
		//	})
		//})
		//ginkgo.Context("When running core commands[show context]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdShowCtx.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Current context:"))
		//	})
		//})
		//ginkgo.Context("When running core commands[--help]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdHelp.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Use \"rit [command] --help\" for more information about a command."))
		//	})
		//})
		//ginkgo.Context("When running core commands[completion bash]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdCompBash.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("complete -F __start_rit rit"))
		//	})
		//})
		//ginkgo.Context("When running core commands[completion zsh]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdCompZsh.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("_complete rit  2>/dev/null"))
		//	})
		//})
		//ginkgo.Context("When running core commands[--version]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdVersion.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("rit version"))
		//	})
		//})
		//ginkgo.Context("When running core commands[add repo]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdAddRepo.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("has been added to your repositories."))
		//	})
		//})
		//ginkgo.Context("When running core commands[list repo]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdListRepo.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("NAME"))
		//	})
		//})
		//ginkgo.Context("When running core commands[update repo]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdUpdateRepo.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("Done."))
		//	})
		//})
		//ginkgo.Context("When running core commands[clean repo]", func() {
		//	ginkgo.It("Should run without errors", func() {
		//		err, out := cmdCleanRepo.RunRit()
		//		fmt.Println(out)
		//		gomega.Expect(err).To(gomega.Succeed())
		//		gomega.Expect(out).To(gomega.ContainSubstring("has been cleaned successfully"))
		//	})
		//})
	})
})
