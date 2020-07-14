package formula

type Repository struct {
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	TreePath string `json:"treePath"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RepositoryFile struct {
	Values []Repository `json:"repositories,omitempty"`
}

type RepoAdder interface {
	Add(d Repository) error
}

type RepoLister interface {
	List() ([]Repository, error)
}

type RepoUpdater interface {
	Update() error
}

type RepoDeleter interface {
	Delete(name string) error
}

type RepoLoader interface {
	Load() error
}

type RepoAddLister interface {
	RepoAdder
	RepoLister
}

type RepoDelLister interface {
	RepoDeleter
	RepoLister
}
