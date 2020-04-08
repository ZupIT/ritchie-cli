package security

// LogoutManager perform user logout
type LogoutManager interface {
	Logout() error
}
