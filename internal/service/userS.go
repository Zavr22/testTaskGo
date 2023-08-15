// Package service
package service

import (
	"context"
	"fmt"
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
	GetUsername(ctx context.Context) ([]string, error)
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
	usernames, err := s.userRepo.GetUsername(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while getting unames to compare")
	}
	for _, username := range usernames {
		if username == user.Username {
			return uuid.Nil, fmt.Errorf("username already exists, %s", err)
		}
	}
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
	usernames, err := s.userRepo.GetUsername(ctx)
	if err != nil {
		return fmt.Errorf("error while getting unames to compare")
	}
	for _, username := range usernames {
		if username == input.NewUsername {
			return fmt.Errorf("username already exists, %s", err)
		}
	}
	return s.userRepo.UpdateProfile(ctx, userID, input)
}

// DeleteProfile is service method that call repo func
func (s *UserService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}
