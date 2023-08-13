package models

// UpdateProfileInput is struct which defines input for changing profile data
type UpdateProfileInput struct {
	NewEmail    string `json:"email"`
	NewUsername string `json:"username"`
	NewPassword string `json:"password"`
	Admin       bool   `json:"admin"`
}
