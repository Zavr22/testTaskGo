package utils

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

var redisClient *redis.Client

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func IsUserValid(username, password string) (bool, error) {
	val, err := redisClient.Keys(redisClient.Context(), "*").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Key does not exist in Redis:", username)
			return false, nil
		}
		log.Println("Redis error while retrieving value:", err)
		return false, err
	}
	for _, value := range val {
		user, err := redisClient.HGetAll(redisClient.Context(), value).Result()
		if err != nil {
			return false, fmt.Errorf("error while getting, %s", err)
		}
		if user["username"] == username {
			err = bcrypt.CompareHashAndPassword([]byte(user["password"]), []byte(password))
			if err != nil {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}

func IsAdmin(username, password string) (bool, error) {
	val, err := redisClient.Keys(redisClient.Context(), "*").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Key does not exist in Redis:", username)
			return false, nil
		}
		log.Println("Redis error while retrieving value:", err)
		return false, err
	}
	for _, value := range val {
		user, err := redisClient.HGetAll(redisClient.Context(), value).Result()
		if err != nil {
			return false, err
		}
		adminValue, _ := strconv.Atoi(user["admin"])
		if user["username"] == username {
			err = bcrypt.CompareHashAndPassword([]byte(user["password"]), []byte(password))
			if err != nil {
				continue
			}
			if adminValue == 1 {
				return true, nil
			}
		}
	}
	return false, nil
}
