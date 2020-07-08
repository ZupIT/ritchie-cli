package formula

type Repo struct {
	Name     string `json:"name"`
	ZipUrl   string `json:"zipUrl"`
	Version  string `json:"version"`
	Priority int    `json:"priority"`
	Token    string `json:"token,omitempty"`
}

type RepoFile struct {
	Values []Repo`json:"repositories,omitempty"`
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
	List() ([]Repo, error)
}

type RepositoryPrioritySetter interface {
	SetPriority(repo Repo, priority int) error
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
