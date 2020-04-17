package server

type Setter interface {
	Set(url string) error
}

type Validator interface {
	Validate() error
}