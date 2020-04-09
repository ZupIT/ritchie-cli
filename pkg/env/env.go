package env

const (
	// Credential resolver
	Credential = "CREDENTIAL"
)

type Resolvers map[string]Resolver

// Resolver is an interface that we can use to resolve reserved envs
type Resolver interface {
	Resolve(name string) (string, error)
}
