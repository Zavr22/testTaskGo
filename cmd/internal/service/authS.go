package service

import (
	"context"
	"github.com/google/uuid"
	"testTask/cmd/models"
)

type Authorization interface {
	SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	SignIn(ctx context.Context, user *models.SignInInput) (string, error)
}

type AuthService struct {
	repo Authorization
}

func NewAuthService(repo Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
	return s.repo.SignUp(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, user *models.SignInInput) (string, error) {
	return s.repo.SignIn(ctx, user)
}
