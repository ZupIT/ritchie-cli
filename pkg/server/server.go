package server

type Setter interface {
	Set(url string) error
}