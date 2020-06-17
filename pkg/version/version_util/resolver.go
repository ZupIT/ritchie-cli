package version_util

type Resolver interface {
	StableVersion() (string, error)
}
