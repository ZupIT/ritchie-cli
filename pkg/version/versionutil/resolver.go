package versionutil

type Resolver interface {
	GetStableVersion() (string, error)
}
