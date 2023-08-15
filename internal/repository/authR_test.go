package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthRepo_SignUp(t *testing.T) {
	// Create a new mini Redis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()

	// Create a new Redis client using the mini Redis server's address
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Create a new instance of AuthRepo with the Redis client
	repo := NewAuthRepo(client)

	// Define the test input
	user := &models.SignUpInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
		Admin:    false,
	}

	// Call the SignUp method
	userID, err := repo.SignUp(context.Background(), user)
	if err != nil {
		t.Errorf("SignUp returned an error: %v", err)
	}

	// Verify the user ID
	if userID == uuid.Nil {
		t.Error("SignUp returned a nil user ID")
	}

	// Verify that the user data is stored in Redis
	userData, err := client.HGetAll(context.Background(), userID.String()).Result()
	if err != nil {
		t.Errorf("Failed to retrieve user data from Redis: %v", err)
	}

	// Verify the stored email
	if userData["email"] != user.Email {
		t.Errorf("Stored email does not match: expected=%s, got=%s", user.Email, userData["email"])
	}

	// Verify the stored username
	if userData["username"] != user.Username {
		t.Errorf("Stored username does not match: expected=%s, got=%s", user.Username, userData["username"])
	}

	// Verify the stored password (using bcrypt comparison)
	err = bcrypt.CompareHashAndPassword([]byte(userData["password"]), []byte(user.Password))
	if err != nil {
		t.Errorf("Stored password is not correct: %v", err)
	}
}

func TestAuthRepo_SignIn(t *testing.T) {
	// Create a new mini Redis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start mini Redis server: %v", err)
	}
	defer s.Close()

	// Create a new Redis client using the mini Redis server's address
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Create a new instance of AuthRepo with the Redis client
	repo := NewAuthRepo(client)

	// Define the test user data
	userData := map[string]interface{}{
		"email":    "test@example.com",
		"username": "testuser",
		"password": "$2a$10$Wz9QlLrFfKZbZxz7PQ9/Su2z6sYJZk2Bm8v9hBkY0c4O2mJswu8Xi", // bcrypt hash of "password123"
		"admin":    false,
	}

	// Store the user data in Redis
	err = client.HSet(context.Background(), "users", uuid.New().String(), fmt.Sprintf("%s:%s", userData["username"], userData["password"])).Err()
	if err != nil {
		t.Fatalf("Failed to store user data in Redis: %v", err)
	}
	err = client.HMSet(context.Background(), "user_id", userData).Err()
	if err != nil {
		t.Fatalf("Failed to store user data in Redis: %v", err)
	}

	// Define the test input
	user := &models.SignInInput{
		Username: "testuser",
		Password: "password123",
	}

	// Call the SignIn method
	err = repo.SignIn(context.Background(), user)
	if err != nil {
		t.Errorf("SignIn returned an error: %v", err)
	}
}
