package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func DBSession() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	dbAddress := os.Getenv("DB_ADDRESS")
	if dbAddress == "" {
		fmt.Println("DB_ADDRESS environment variable is not set")
		return
	}

	client = redis.NewClient(&redis.Options{
		Addr:     dbAddress,
		Password: "", 
		DB:       0,  
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("Could not connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis")
}

// GetClient returns the Redis client for use in other parts of the application
func GetClient() *redis.Client {
	return client
}