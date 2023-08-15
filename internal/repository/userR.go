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

// UserRepo has an internal db object
type UserRepo struct {
	client *redis.Client
}

// NewUserRepo used to init UserAP
func NewUserRepo(client *redis.Client) *UserRepo {
	return &UserRepo{client: client}
}

// CreateUser is used to create user by admin using redis
func (r *UserRepo) CreateUser(ctx context.Context, user *models.SignUpInput) (uuid.UUID, error) {
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

// GetAllUsers is used to get users using redis
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
			return nil, fmt.Errorf("Error occurred while retrieving user data for ID %s: %s", id, err)
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

// GetUser is used to get user by id using redis
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

// UpdateProfile is used to update user profile by admin  using redis
func (r *UserRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) error {
	tx := r.client.TxPipeline()
	if input.NewEmail != "" {
		tx.HSet(ctx, userID.String(), "email", input.NewEmail)
	}
	if input.NewUsername != "" {
		tx.HSet(ctx, userID.String(), "username", input.NewUsername)
	}
	if input.NewPassword != "" {
		tx.HSet(ctx, userID.String(), "password", input.NewPassword)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("error while hashing password, %s", err)
			return fmt.Errorf("error while hashing password, %s", err)
		}
		tx.HSet(ctx, "users", userID.String(), fmt.Sprintf("%s:%s", input.NewUsername, hashedPassword))
	}
	tx.HSet(ctx, userID.String(), "admin", input.Admin)
	_, err := tx.Exec(ctx)
	if err != nil {
		if err == redis.TxFailedErr {
			log.Printf("transaction failed due to concurrent modification, %s", err)
			return fmt.Errorf("transaction failed due to concurrent modificationб, %s", err)
		}
		log.Printf("error in transaction, %s", err)
		return fmt.Errorf("error in transaction, %s", err)
	}
	return nil
}

// DeleteProfile is used to delete user profile by admin  using redis
func (r *UserRepo) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	tx := r.client.TxPipeline()
	tx.Del(ctx, userID.String())
	tx.Del(ctx, userID.String(), "*:*")
	_, err := tx.Exec(ctx)
	if err != nil {
		if err == redis.TxFailedErr {
			log.Println("transaction failed due to concurrent modification")
			return fmt.Errorf("transaction failed due to concurrent modification")
		}
		log.Printf("error in transaction, %s", err)
		return fmt.Errorf("error in transaction, %s", err)
	}

	return nil
}
