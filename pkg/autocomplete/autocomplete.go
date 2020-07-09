package autocomplete

import "github.com/spf13/cobra"

type (
	BashCommand struct {
		LastCommand string
		RootCommand string
		Commands    string
		Level       int
	}

	CompletionCommand struct {
		Content []string
		Before  string
	}
)

type Generator interface {
	Generate(s ShellName, cmd *cobra.Command) (string, error)
}

type ShellName string

func (s ShellName) String() string {
	return string(s)
}
