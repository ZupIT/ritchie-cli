package server

const (
	serverFilePattern = "%s/server.json"
)

type Config struct {
	Organization string
	URL          string
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
