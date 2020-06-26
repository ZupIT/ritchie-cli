package version

type Resolver interface {
	StableVersion(fromCache bool) (string, error)
}
