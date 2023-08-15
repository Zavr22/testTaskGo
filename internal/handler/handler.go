package handler

import (
	"context"
	"github.com/Zavr22/testTaskGo/internal/middleware"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Authorization interface consists of methods of auth service
type Authorization interface {
	SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	SignIn(ctx context.Context, user *models.SignInInput) error
}

// User interface consists of user service methods
type User interface {
	CreateUser(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	GetAllUsers(ctx context.Context) ([]*models.UserResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

// Handler struct contains of interfaces of user and auth service
type Handler struct {
	userS User
	authS Authorization
}

// NewHandler is used to init handler obj
func NewHandler(userS User, authS Authorization) *Handler {
	return &Handler{userS: userS, authS: authS}
}

// InitRoutes is used to init routes for web service
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
