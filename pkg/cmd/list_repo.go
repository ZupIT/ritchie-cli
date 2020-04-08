package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

// listRepoCmd type for list repo command
type listRepoCmd struct {
	formula.RepoLister
}

// NewListRepoCmd creates a new cmd instance
func NewListRepoCmd(ls formula.RepoLister) *cobra.Command {
	l := &listRepoCmd{ls}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "List all repositories.",
		Example: "rit list repo",
		RunE:    l.RunFunc(),
	}

	return cmd
}

func (l listRepoCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		rr, err := l.List()
		if err != nil {
			return err
		}

		printList(rr)

		return nil
	}
}

func printList(rr []formula.Repository) {
	table := uitable.New()
	table.AddRow("NAME", "URL")
	for _, re := range rr {
		table.AddRow(re.Name, re.TreePath)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))
}
