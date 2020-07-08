package credential

const (
	// Other credential path /admin
	Other Type = "admin"
	// Me credential path /me
	Me Type = "me"
	// Org credential path /org
	Org Type = "org"
)

// Info represents a credential information of the user.
type Detail struct {
	Username   string     `json:"username"`
	Credential Credential `json:"credential"`
	Service    string     `json:"service"`
	Type       Type       `json:"type"`
}

type Type string

func (t Type) String() string {
	return string(t)
}

// A Credential represents the key-value pairs for the Service (User/Pass, Github, Jenkins, etc).
type Credential map[string]string

// Field represents a credential field associated with your type.
type Field struct {
	Name string `json:"field"`
	Type string `json:"type"`
}

// Fields represents a collection of credential fields returned by the Server (Team).
// Fields are used on single to represents providers.json
type Fields map[string][]Field

type Setter interface {
	Set(d Detail) error
}

type Finder interface {
	Find(service string) (Detail, error)
}

type Settings interface {
	Fields() (Fields, error)
}

type SingleSettings interface {
	ReadCredentials(path string) (Fields, error)
	WriteCredentials(fields Fields, path string) error
	DefaultCredentials() error
}