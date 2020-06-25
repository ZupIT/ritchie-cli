package version

type Resolver interface {
	StableVersionForCmd() (string, error)
	StableVersion() (string, error)
}
