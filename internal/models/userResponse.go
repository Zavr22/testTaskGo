package models

import "github.com/google/uuid"

// UserResponse is a struct that defines fields for user response
type UserResponse struct {
	ID       uuid.UUID `json:"ID"`
	Username string    `json:"user_name"`
	Email    string    `json:"email"`
}
