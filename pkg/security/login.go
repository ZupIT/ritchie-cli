package security

// LoginManager perform user login
type LoginManager interface {
	Login() error
}
