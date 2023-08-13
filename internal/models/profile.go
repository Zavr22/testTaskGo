package models

import "github.com/google/uuid"

// UserProfile struct is a user profile models that defines fields for profile
type UserProfile struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Admin    bool      `json:"admin"`
}
