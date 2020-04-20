package server

type Validator interface {
	Validate() error
}

type Setter interface {
	Set(url string) error
}

type Finder interface {
	Find() (string, error)
}