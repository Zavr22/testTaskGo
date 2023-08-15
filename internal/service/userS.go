package service

import (
	"context"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=userS.go -destination=mocks/user_mock.go

type User interface {
	CreateUser(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	GetAllUsers(ctx context.Context) ([]*models.UserResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type UserService struct {
	userRepo User
}

func NewUserService(userRepo User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.UserResponse, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error) {
	return s.userRepo.GetUser(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error {
	return s.userRepo.UpdateProfile(ctx, userID, input)
}

func (s *UserService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}
