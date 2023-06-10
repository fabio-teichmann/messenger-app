package main

import (
	"context"
	"fmt"
	"log"
	"messenger-app/storage"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// tests.RunV1()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.MongoConfig{
		// Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
	}
	client, close, err := storage.NewMongoConnection(config)

	if err != nil {
		panic(err)
	}
	defer close()

	// make sure connection is established
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("connection to mongodb not established")
		panic(err)
	}

	// r := gin.Default()

	// fmt.Print(r)

}
