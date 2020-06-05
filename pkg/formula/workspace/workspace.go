package workspace

const (
	workspacesPattern          = "%s/formula_workspaces.json"
	DefaultWorkspaceName       = "Default"
	DefaultWorkspaceDirPattern = "%s/ritchie-formulas-local"
)

type Workspaces map[string]string

type Workspace struct {
	Name string `json:"name"`
	Dir  string `json:"dir"`
}

type Adder interface {
	Add(workspace Workspace) error
}

type Lister interface {
	List() (Workspaces, error)
}

type AddLister interface {
	Adder
	Lister
}
