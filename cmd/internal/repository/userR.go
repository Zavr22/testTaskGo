package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"testTask/cmd/models"
)

// UserRepo has an internal db object
type UserRepo struct {
	client *redis.Client
}

// NewUserRepo used to init UsesAP
func NewUserRepo(client *redis.Client) *UserRepo {
	return &UserRepo{client: client}
}

func (r *UserRepo) CreateUser(ctx context.Context, email, username, password string, admin bool) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*models.UserProfile, error) {
	return make([]*models.UserProfile, 0), nil
}

func (r *UserRepo) GetUser(ctx context.Context, userID uuid.UUID) (models.UserProfile, error) {
	return models.UserProfile{}, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error {
	return nil
}

func (r *UserRepo) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return nil
}
