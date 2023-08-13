package utils

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"strings"
)

var redisClient *redis.Client

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func IsUserValid(authToken []string) (bool, error) {
	val, err := redisClient.HGet(redisClient.Context(), "users", authToken[0]).Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Key does not exist in Redis:", authToken[0])
			return false, nil
		}
		log.Println("Redis error while retrieving value:", err)
		return false, err
	}
	hashedPassword := val[strings.Index(val, ":")+1:]
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(authToken[1]))
	if err == nil {
		return true, nil
	} else if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Println("Value from Redis does not match expected password:", val)
		return false, nil
	} else {
		log.Println("Error comparing passwords:", err)
	}
	return false, err
}

func IsAdmin(authToken []string) (bool, error) {
	user, err := redisClient.HGetAll(redisClient.Context(), authToken[0]).Result()
	adminValue, _ := strconv.Atoi(user["admin"])
	if adminValue == 1 {
		return true, nil
	}
	return false, err
}
