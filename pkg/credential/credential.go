package credential

const (
	// Admin role
	Admin = "admin"
	// Me credential
	Me = "me"
	// Other credential
	Other = "other"
)

// Info represents a credential information of the user.
type Detail struct {
	Username   string     `json:"username"`
	Credential Credential `json:"credential"`
	Service    string     `json:"service"`
}

// A Credential represents the key-value pairs for the Service (User/Pass, Github, Jenkins, etc).
type Credential map[string]string

// Field represents a credential field associated with your type.
type Field struct {
	Name string `json:"field"`
	Type string `json:"type"`
}

// Fields represents a collection of credential fields that can be returned by the Server (Team Edition).
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
