package formula

type Repo struct {
	Name    string `json:"name"`
	ZipUrl  string `json:"zipUrl"`
	Version string `json:"version"`
	Current bool   `json:"current"`
	Token   string `json:"token,omitempty"`
}

type Repos map[string]Repo

type Tag struct {
	Name   string `json:"name"`
	ZipUrl string `json:"zipball_url"`
}

type Tags map[string]string

type RepositoryAdder interface {
	Add(d Repo) error
}

type RepositoryLister interface {
	List() (Repos, error)
}

type RepositoryUpdater interface {
	Update() error
}

type RepositoryDeleter interface {
	Delete(name string) error
}

type RepositoryCleaner interface {
	Clean(name string) error
}

type RepositoryLoader interface {
	Load() error
}

type RepositoryAddLister interface {
	RepositoryAdder
	RepositoryLister
}

type RepositoryDelLister interface {
	RepositoryDeleter
	RepositoryLister
}
