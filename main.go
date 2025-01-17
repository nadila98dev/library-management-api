package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main(){
	godotenv.Load()
	fmt.Println(("Test Redis Coneect"))


	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: "", 
		DB:       0,
	})

	ping, err := client.Ping(context.Background()).Result();
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println((ping))

	type Person struct {
	ID string
	Name  string             `json:"name"`
	Age  int                `json:"age"`
	Occupation string		`json:"occupation"`
}


	keyId := uuid.NewString()

	jsonString, err := json.Marshal(Person{
		ID: keyId,
		Name:  "Mark",
		Age:   25,
		Occupation: "Musician",
	})

	if err != nil {
		fmt.Println("Failed to marshall", err.Error())
		return
	}

	keyIds := fmt.Sprintf("person", keyId)
	err = client.Set((context.Background()), keyIds, jsonString, 0).Err()
	if err != nil {
		fmt.Println("Failed to set a value ")
		return
	}

	val, err := client.Get((context.Background()), keyIds).Result()
	if err != nil {
		fmt.Println("Failed to Get a value")
		return
	}

	fmt.Println(("test redis retrieved"), val)


}