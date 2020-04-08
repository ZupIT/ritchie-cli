package env

const (
	// Credential resolver
	Credential = "CREDENTIAL"
	// Team version
	Team = "team"
	// Single version
	Single = "single"
)

var (
	//Url Server
	ServerURL = "https://ritchie-server.zup.io"
	// Edition option
	Edition = Team
)

type Resolvers map[string]Resolver

// Resolver is an interface that we can use to resolve reserved envs
type Resolver interface {
	Resolve(name string) (string, error)
}
