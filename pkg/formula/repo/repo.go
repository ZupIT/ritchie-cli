package repo

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

type Adder interface {
	Add(d Repository) error
}

type Lister interface {
	List() ([]Repository, error)
}

type Updater interface {
	Update() error
}

type Deleter interface {
	Delete(name string) error
}

type Cleaner interface {
	Clean(name string) error
}

type Loader interface {
	Load() error
}
