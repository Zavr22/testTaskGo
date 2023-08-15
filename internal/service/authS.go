// Package service
package service

import (
	"context"
	"fmt"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=authS.go -destination=mocks/auth_mock.go

// Authorization interface consists of methods of auth repo
type Authorization interface {
	SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error)
	SignIn(ctx context.Context, user *models.SignInInput) error
	GetUsername(ctx context.Context) ([]string, error)
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
	usernames, err := s.repo.GetUsername(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while getting unames to compare")
	}
	for _, username := range usernames {
		if username == user.Username {
			return uuid.Nil, fmt.Errorf("username already exists, %s", err)
		}
	}
	return s.repo.SignUp(ctx, user)
}

// SignIn is service method that call repo func
func (s *AuthService) SignIn(ctx context.Context, user *models.SignInInput) error {
	return s.repo.SignIn(ctx, user)
}
