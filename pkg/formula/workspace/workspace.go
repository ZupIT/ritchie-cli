package workspace

const workspacesPattern = "%s/formula_workspaces.json"

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
