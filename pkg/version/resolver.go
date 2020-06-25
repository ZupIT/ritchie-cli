package version

type Resolver interface {
	StableVersion(fromCmd bool) (string, error)
}
