package autocomplete

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
	Generate(s ShellName) (string, error)
}

type ShellName string

func (s ShellName) String() string {
	return string(s)
}
