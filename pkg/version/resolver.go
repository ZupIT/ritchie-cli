package version

type Resolver interface {
	StableVersion() (string, error)
}
