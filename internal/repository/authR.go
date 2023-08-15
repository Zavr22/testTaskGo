package repository

import (
	"context"
	"fmt"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// AuthRepo has an internal db object
type AuthRepo struct {
	client *redis.Client
}

// NewAuthRepo is used to init auth repo
func NewAuthRepo(client *redis.Client) *AuthRepo {
	return &AuthRepo{client: client}
}

// SignUp func is used to sign user up using redis
func (r *AuthRepo) SignUp(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
	userID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}
	userData := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
		"password": string(hashedPassword),
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

// SignIn func is used to sign user in using redis
func (r *AuthRepo) SignIn(ctx context.Context, user *models.SignInInput) error {
	log.Printf("Attempting to sign in user with username: %s", user.Username)
	users, err := r.client.HGetAll(ctx, "users").Result()
	if err != nil {
		log.Printf("Error occurred while retrieving users: %s", err)
		return fmt.Errorf("error occurred: %s", err)
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
		return fmt.Errorf("user not found")
	}
	userData, err := r.client.HGetAll(ctx, userID).Result()
	if err != nil {
		log.Printf("Error occurred while retrieving user data: %s", err)
		return fmt.Errorf("error occurred: %s", err)
	}
	storedPassword := userData["password"]
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		log.Printf("Incorrect password for user with username %s", user.Username)
		return fmt.Errorf("incorrect password")
	}
	log.Printf("User with username %s successfully signed in", user.Username)
	log.Printf("User data: %+v", userData)

	return nil
}

// GetUsername is used to get usernames from db
func (r *AuthRepo) GetUsername(ctx context.Context) ([]string, error) {
	userNameArr := make([]string, 0)
	val, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Key does not exist in Redis:")
			return nil, fmt.Errorf("key does not exist in Redis:%s", err)
		}
		log.Println("Redis error while retrieving value:", err)
		return nil, fmt.Errorf("redis error while retrieving value: %s", err)
	}
	for _, value := range val {
		user, err := r.client.HGetAll(ctx, value).Result()
		if err != nil {
			return nil, fmt.Errorf("error while getting username")
		}
		userNameArr = append(userNameArr, user["username"])
	}
	return userNameArr, nil
}
