package server

const (
	serverFilePattern = "%s/server.json"
)

type Config struct {
	Organization string `json:"organization"`
	URL          string `json:"url"`
}

type Setter interface {
	Set(Config) error
}

type Finder interface {
	Find() (Config, error)
}

type FindSetter interface {
	Finder
	Setter
}
