package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"strings"
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

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*models.UserResponse, error) {
	var userProfiles []*models.UserResponse
	users, err := r.client.HGetAll(ctx, "users").Result()
	if err != nil {
		log.Printf("Error occurred while retrieving users: %s", err)
		return nil, fmt.Errorf("error occurred: %s", err)
	}
	for id := range users {
		userData, err := r.client.HGetAll(ctx, id).Result()
		if err != nil {
			log.Printf("Error occurred while retrieving user data for ID %s: %s", id, err)
			continue
		}
		userProfile := &models.UserResponse{
			ID:       uuid.MustParse(strings.Split(id, ":")[0]),
			Email:    userData["email"],
			Username: userData["username"],
		}
		userProfiles = append(userProfiles, userProfile)
		log.Println(userData)
	}
	return userProfiles, nil
}

func (r *UserRepo) GetUser(ctx context.Context, userID uuid.UUID) (models.UserResponse, error) {
	userData, err := r.client.HGetAll(ctx, userID.String()).Result()
	users := r.client.HGetAll(ctx, "users")
	log.Println(users)
	if err != nil {
		log.Printf("Error occurred while retrieving user data for ID %s: %s", userID, err)
		return models.UserResponse{}, fmt.Errorf("Error occurred while retrieving user data for ID %s: %s", userID, err)
	}
	userProfile := models.UserResponse{
		ID:       userID,
		Email:    userData["email"],
		Username: userData["username"],
	}
	log.Println(userData)

	return userProfile, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error {
	return nil
}

func (r *UserRepo) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	err := r.client.Del(ctx, userID.String()).Err()
	if err != nil {
		log.Printf("error while delete user, %s", err)
		return err
	}
	return nil
}
