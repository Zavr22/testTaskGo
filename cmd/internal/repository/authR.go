package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
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
	userData := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
		"password": user.Password,
		"admin":    user.Admin,
	}
	err = r.client.HMSet(ctx, userID.String(), userData).Err()
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
	log.Printf("Attempting to sign in user with username: %s", user.Username)
	users, err := r.client.HGetAll(ctx, "users").Result()
	if err != nil {
		log.Printf("Error occurred while retrieving users: %s", err)
		return "", fmt.Errorf("error occurred: %s", err)
	}
	var userID string
	for id, userData := range users {
		username := strings.Split(userData, ":")[0]
		if username == user.Username {
			userID = id
			break
		}
	}
	if userID == "" {
		log.Printf("User with username %s not found", user.Username)
		return "", fmt.Errorf("user not found")
	}
	userData, err := r.client.HGetAll(ctx, userID).Result()
	if err != nil {
		log.Printf("Error occurred while retrieving user data: %s", err)
		return "", fmt.Errorf("error occurred: %s", err)
	}
	storedPassword, ok := userData["password"]
	if !ok {
		log.Printf("Password not found for user with username %s", user.Username)
		return "", fmt.Errorf("user password not found")
	}

	if storedPassword != user.Password {
		log.Printf("Incorrect password for user with username %s", user.Username)
		return "", fmt.Errorf("incorrect password")
	}

	log.Printf("User with username %s successfully signed in", user.Username)
	log.Printf("User data: %+v", userData)
	encodedUserID := base64.StdEncoding.EncodeToString([]byte(userID))
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(user.Password))
	token := encodedUserID + ":" + encodedPassword
	return token, nil
}
