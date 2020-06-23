package server

const (
	serverFilePattern = "%s/server.json"
)

type Config struct {
	Organization string `json:"organization"`
	URL          string `json:"url"`
	Otp          bool   `json:"otp"`
	PinningKey   string `json:"pinningKey"`
	PinningAddr  string `json:"pinningAddr"`
}

type Setter interface {
	Set(*Config) error
}

type Finder interface {
	Find() (Config, error)
}

type FindSetter interface {
	Finder
	Setter
}
