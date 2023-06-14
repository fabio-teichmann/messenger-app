package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"messenger-app/models"
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

	ac := models.NewAppControler(client)

	event, err := ac.GetEventByMessageId(ctx, 0)
	if err != nil {
		fmt.Println(err)
	}
	res, _ := json.Marshal(event)
	fmt.Println("MessageID 0:", string(res))

	event1, err := ac.GetEventByMessageId(ctx, 1)
	if err != nil {
		fmt.Println(err)
	}
	res, _ = json.Marshal(event1)
	fmt.Println("MessageID 1:", string(res))

	event2, err := ac.GetEventByMessageId(ctx, 2)
	if err != nil {
		fmt.Println(err)
	}
	res, _ = json.Marshal(event2)
	fmt.Println("MessageID 2:", string(res))

	count, err := ac.CountMessagesBySubjectId(ctx, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("# messages id 0:", count)

	count, err = ac.CountMessagesBySubjectId(ctx, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("# messages id 1:", count)

	user1 := models.User{ID: 1, Name: "user1"}
	user2 := models.User{ID: 2, Name: "user2"}
	user3 := models.User{ID: 3, Name: "user3"}

	count, err = ac.CountMessagesSentByUser(ctx, &user1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("# messages sender id 1:", count)

	count, err = ac.CountMessagesSentByUser(ctx, &user2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("# messages sender id 2:", count)

	count, err = ac.CountMessagesSentByUser(ctx, &user3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("# messages sender id 3:", count)

	// event := tests.CreateTestEvent("message string")

	// coll := client.Database("messenger-test").Collection("events")
	// // _, err = coll.InsertOne(ctx, event)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// fmt.Println("Collection:", *coll)

	// var results []models.Event
	// cursor, err := coll.Find(ctx, bson.D{})
	// if err != nil {
	// 	panic(err)
	// }

	// if err = cursor.All(context.TODO(), &results); err != nil {
	// 	panic(err)
	// }
	// for _, result := range results {
	// 	res, _ := json.Marshal(result)
	// 	fmt.Println(string(res))
	// }

	// fmt.Println("Results:", res)

	// r := gin.Default()

	// fmt.Print(r)

}
