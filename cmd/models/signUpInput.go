package models

// SignUpInput is a struct which defines fields for signup input
type SignUpInput struct {
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Admin    bool      `json:"admin"`
}
