package models

import "github.com/google/uuid"

type CreateUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}
