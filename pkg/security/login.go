package security

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginManager perform user login
type LoginManager interface {
	Login(User) error
}
