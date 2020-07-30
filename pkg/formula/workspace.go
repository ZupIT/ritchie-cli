package formula

const (
	DefaultWorkspaceName = "Default"
	DefaultWorkspaceDir  = "ritchie-formulas-local"
	WorkspacesFile       = "formula_workspaces.json"
)

type (
	Workspace struct {
		Name string `json:"name"`
		Dir  string `json:"dir"`
	}

	Workspaces map[string]string
)

type WorkspaceAdder interface {
	Add(workspace Workspace) error
}

type WorkspaceLister interface {
	List() (Workspaces, error)
}

type WorkspaceValidator interface {
	Validate(workspace Workspace) error
}

type WorkspaceAddListValidator interface {
	WorkspaceAdder
	WorkspaceLister
	WorkspaceValidator
}
