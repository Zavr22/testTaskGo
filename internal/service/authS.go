// Package service
package service

import (
	"context"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=authS.go -destination=mocks/auth_mock.go

// Authorization interface consists of methods of auth repo
type Authorization interface {
	SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	SignIn(ctx context.Context, user *models.SignInInput) error
}

// AuthService struct contains of auth repo interface
type AuthService struct {
	repo Authorization
}

// NewAuthService is used to init auth repo obj
func NewAuthService(repo Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// SignUp is service method that call repo func
func (s *AuthService) SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
	return s.repo.SignUp(ctx, user)
}

// SignIn is service method that call repo func
func (s *AuthService) SignIn(ctx context.Context, user *models.SignInInput) error {
	return s.repo.SignIn(ctx, user)
}
