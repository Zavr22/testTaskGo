package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"testTask/cmd/models"
)

type Authorization interface {
	SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	SignIn(ctx context.Context, user *models.SignInInput) error
}

type User interface {
	CreateUser(ctx context.Context, email, username, password string, admin bool) (uuid.UUID, error)
	GetAllUsers(ctx context.Context) ([]*models.UserProfile, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.UserProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type Handler struct {
	userS User
	authS Authorization
}

func NewHandler(userS User, authS Authorization) *Handler {
	return &Handler{userS: userS, authS: authS}
}

func (h *Handler) InitRoutes(router *echo.Echo) *echo.Echo {
	return router
}
