package service

import (
	"context"
	"github.com/google/uuid"
	"testTask/cmd/models"
)

type User interface {
	CreateUser(ctx context.Context, email, username, password string, admin bool) (uuid.UUID, error)
	GetAllUsers(ctx context.Context) ([]*models.UserProfile, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.UserProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type UserService struct {
	userRepo User
}

func NewUserService(userRepo User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, email, username, password string, admin bool) (uuid.UUID, error) {
	return s.userRepo.CreateUser(ctx, email, username, password, admin)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.UserProfile, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (models.UserProfile, error) {
	return s.userRepo.GetUser(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error {
	return s.userRepo.UpdateProfile(ctx, userID, input)
}

func (s *UserService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}
