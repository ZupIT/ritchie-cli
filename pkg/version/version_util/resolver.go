package version_util

type Resolver interface {
	GetStableVersion() (string, error)
}
