package session

const (
	sessionFilePattern    = "%s/session"
	passphraseFilePattern = "%s/passphrase"
)

// Session represents a security session of the user.
type Session struct {
	AccessToken  string `json:"access_token"`
	Organization string `json:"organization"`
	Username     string `json:"username"`
	Secret       string `json:"password"`
	TTL          int64  `json:"ttl"`
}

type Manager interface {
	Create(s Session) error
	Current() (Session, error)
	Destroy() error
}

type Validator interface {
	Validate() error
}
