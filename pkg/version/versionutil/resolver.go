package versionutil

type Resolver interface {
	GetCurrentVersion() (string, error)
	GetStableVersion() (string, error)
}