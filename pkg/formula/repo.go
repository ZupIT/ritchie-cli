package formula

type Repo struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Version  string `json:"version"`
	Token    string `json:"token,omitempty"`
	Priority int    `json:"priority"`
}

type Repos []Repo

func (r Repos) Len() int {
	return len(r)
}

func (r Repos) Less(i, j int) bool {
	return r[i].Priority < r[j].Priority
}

func (r Repos) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

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

type RepositoryPrioritySetter interface {
	SetPriority(repoName string, priority int) error
}

type RepositoryAddLister interface {
	RepositoryAdder
	RepositoryLister
}

type RepositoryDelLister interface {
	RepositoryDeleter
	RepositoryLister
}
