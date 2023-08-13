package models

import "github.com/google/uuid"

type UserResponse struct {
	ID       uuid.UUID `json:"ID"`
	Username string    `json:"user_name"`
	Email    string    `json:"email"`
}
