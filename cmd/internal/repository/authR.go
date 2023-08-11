package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"testTask/cmd/models"
)

// AuthRepo has an internal db object
type AuthRepo struct {
	client *redis.Client
}

// NewAuthRepo is used to init auth repo
func NewAuthRepo(client *redis.Client) *AuthRepo {
	return &AuthRepo{client: client}
}

func (r *AuthRepo) SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
	userID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}
	err = r.client.HSet(ctx, "users", userID.String(), fmt.Sprintf("%s:%s", user.Username, string(hashedPassword))).Err()
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func (r *AuthRepo) SignIn(ctx context.Context, user *models.SignInInput) (string, error) {
	usersData, err := r.client.HGetAll(ctx, "users").Result()
	if err != nil {
		return "", err
	}
	var userID string
	var userCredentials string
	for id, credentials := range usersData {
		if credentials == fmt.Sprintf("%s:%s", user.Username, user.Password) {
			userID = id
			userCredentials = credentials
			break
		}
	}
	if userID == "" {
		return "", fmt.Errorf("User not found")
	}
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", userID, userCredentials)))
	return token, nil
}
