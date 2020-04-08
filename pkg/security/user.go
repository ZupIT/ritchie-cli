package security

// User represents the user of the system.
type User struct {
	Organization string `json:"organization"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// Manager manages user creation and deletion
type UserManager interface {
	Create(User) error
	Delete(User) error
}
