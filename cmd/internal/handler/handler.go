package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"testTask/cmd/internal/middleware"
	"testTask/cmd/models"
)

type Authorization interface {
	SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	SignIn(ctx context.Context, user *models.SignInInput) (string, error)
}

type User interface {
	CreateUser(ctx context.Context, email, username, password string, admin bool) (uuid.UUID, error)
	GetAllUsers(ctx context.Context) ([]*models.UserResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error)
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

	auth := router.Group("/auth")
	auth.POST("/sign_up", h.SignUp)
	auth.POST("/sign_in", h.SignIn)

	api := router.Group("api")
	api.POST("/users", h.CreateUser, middleware.AdminMiddleware())
	api.GET("/users", h.GetUsers)
	api.GET("/users/:id", h.GetUserByID)
	api.PUT("/users/:id", h.UpdateUser, middleware.AdminMiddleware())
	api.DELETE("/users/:id", h.DeleteUser, middleware.AdminMiddleware())

	router.Logger.Fatal(router.Start(":9000"))
	return router
}
