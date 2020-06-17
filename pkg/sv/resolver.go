package sv

type Resolver interface {
	StableVersion() (string, error)
}
