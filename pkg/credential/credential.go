package credential

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

type ListCredDatas []ListCredData

type ListCredData struct {
	Provider string
	Name string
	Value string
	Context string
}

// Fields are used to represents providers.json
type Fields map[string][]Field

type Setter interface {
	Set(d Detail) error
}

type CredFinder interface {
	Find(service string) (Detail, error)
}

type Operations interface {
	ReadCredentialsFields(path string) (Fields, error)
	ReadCredentialsValue(path string) string
	WriteCredentialsFields(fields Fields, path string) error
	WriteDefaultCredentialsFields(path string) error
}

