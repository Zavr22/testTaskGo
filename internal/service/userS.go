// Package service
package service

import (
	"context"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=userS.go -destination=mocks/user_mock.go

// User interface consists of user repo methods
type User interface {
	CreateUser(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	GetAllUsers(ctx context.Context) ([]*models.UserResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

// UserService contains of User repo interface
type UserService struct {
	userRepo User
}

// NewUserService init repo obj
func NewUserService(userRepo User) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser is service method that call repo func
func (s *UserService) CreateUser(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
	return s.userRepo.CreateUser(ctx, user)
}

// GetAllUsers is service method that call repo func
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.UserResponse, error) {
	return s.userRepo.GetAllUsers(ctx)
}

// GetUser is service method that call repo func
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error) {
	return s.userRepo.GetUser(ctx, userID)
}

// UpdateProfile is service method that call repo func
func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error {
	return s.userRepo.UpdateProfile(ctx, userID, input)
}

// DeleteProfile is service method that call repo func
func (s *UserService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}
