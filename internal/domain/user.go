package domain

// User represents a user in the system.
// This is a domain entity.
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
