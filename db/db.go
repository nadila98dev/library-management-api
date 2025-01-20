package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func DBSession() {
	
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}


	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: "", 
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("Could not connect to Redis")
		return
	}
	fmt.Println("Connected to Redis")


	err := client.Set(crx, "key", "value", 0).Err()
	if err != nil {
		fmt.Println(err)
		return
	}

	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("key", val)

	defer client.Close()
		


}