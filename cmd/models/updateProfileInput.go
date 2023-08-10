package models

// UpdateProfileInput is struct which defines input for changing profile data
type UpdateProfileInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}
