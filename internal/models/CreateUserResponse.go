package models

import "github.com/google/uuid"

// CreateUserResponse is a struct that defines field for create user and sign up response
type CreateUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}
