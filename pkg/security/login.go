package security

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Totp     string `json:"totp"`
}

// LoginManager perform user login
type LoginManager interface {
	Login(User) error
}
