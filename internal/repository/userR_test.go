package repository

import (
	"context"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/alicebob/miniredis"
	"github.com/google/uuid"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_CreateUser(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := NewUserRepo(client)
	user := &models.SignUpInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
		Admin:    false,
	}
	userID, err := repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	userData, err := client.HGetAll(context.Background(), userID.String()).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve user data from Redis: %v", err)
	}

	assert.Equal(t, user.Email, userData["email"], "Email does not match")
	assert.Equal(t, user.Username, userData["username"], "Username does not match")
}

func TestUserRepo_GetAllUsers(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := NewUserRepo(client)
	user1 := &models.SignUpInput{
		Email:    "user1@example.com",
		Username: "user1",
		Password: "password123",
		Admin:    false,
	}
	user2 := &models.SignUpInput{
		Email:    "user2@example.com",
		Username: "user2",
		Password: "password456",
		Admin:    true,
	}

	_, err = repo.CreateUser(context.Background(), user1)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	_, err = repo.CreateUser(context.Background(), user2)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}
	var users []*models.UserResponse
	users, err = repo.GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}
	assert.Len(t, users, 2, "Unexpected number of users")
	assert.IsType(t, []*models.UserResponse{}, users)
}

func TestUserRepo_UpdateProfile(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := NewUserRepo(client)
	input := models.UpdateProfileInput{
		NewEmail:    "updated@example.com",
		NewUsername: "updateduser",
		NewPassword: "newpassword",
		Admin:       true,
	}
	user := &models.SignUpInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
		Admin:    false,
	}
	userID, err := repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	err = repo.UpdateProfile(context.Background(), userID, input)
	if err != nil {
		t.Fatalf("Failed to update profile: %v", err)
	}

	// Retrieve the updated user data from Redis
	userData, err := client.HGetAll(context.Background(), userID.String()).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve updated user data from Redis: %v", err)
	}

	assert.Equal(t, input.NewEmail, userData["email"], "Email does not match")
	assert.Equal(t, input.NewUsername, userData["username"], "Username does not match")
}

func TestUserRepo_DeleteProfile(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := NewUserRepo(client)

	user := &models.SignUpInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
		Admin:    false,
	}
	userID, err := repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	err = repo.DeleteProfile(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to delete profile: %v", err)
	}
	exists, err := client.Exists(context.Background(), userID.String()).Result()
	if err != nil {
		t.Fatalf("Failed to check if user data exists in Redis: %v", err)
	}

	assert.Equal(t, int64(0), exists, "User data still exists in Redis after deletion")
}

func TestUserRepo_GetUser(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := NewUserRepo(client)
	userID := uuid.New()
	userData := map[string]string{
		"email":    "user1@example.com",
		"username": "user1",
	}
	err = client.HMSet(context.Background(), userID.String(), userData).Err()
	if err != nil {
		t.Fatalf("Failed to store test user data in Redis: %v", err)
	}
	userProfile, err := repo.GetUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("Failed to get user profile: %v", err)
	}
	expectedProfile := models.UserResponse{
		ID:       userID,
		Email:    userData["email"],
		Username: userData["username"],
	}
	assert.Equal(t, expectedProfile, userProfile, "User profile does not match")
}
