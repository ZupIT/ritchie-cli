package formula

type Repo struct {
	Name     RepoName    `json:"name"`
	Version  RepoVersion `json:"version"`
	Url      string      `json:"url"`
	Token    string      `json:"token,omitempty"`
	Priority int         `json:"priority"`
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

type RepoName string

func (r RepoName) String() string {
	return string(r)
}

type RepoVersion string

func (r RepoVersion) String() string {
	return string(r)
}

type RepositoryAdder interface {
	Add(repo Repo) error
}

type RepositoryLister interface {
	List() (Repos, error)
}

type RepositoryUpdater interface {
	Update(name RepoName, version RepoVersion) error
}

type RepositoryDeleter interface {
	Delete(name RepoName) error
}

type RepositoryPrioritySetter interface {
	SetPriority(name RepoName, priority int) error
}

type RepositoryCreator interface {
	Create(repo Repo) error
}

type RepositoryAddLister interface {
	RepositoryAdder
	RepositoryLister
}

type RepositoryDelLister interface {
	RepositoryDeleter
	RepositoryLister
}

type RepositoryListCreator interface {
	RepositoryLister
	RepositoryCreator
}

type RepositoryListUpdater interface {
	RepositoryLister
	RepositoryUpdater
}
