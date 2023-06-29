package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("connection to mongodb not established")
		panic(err)
	}
	ctx = context.WithValue(ctx, models.TestDBKey, os.Getenv("TEST_MONGODB_NAME"))
	ctx = context.WithValue(ctx, models.TestCollectionKey, os.Getenv("TEST_COLLECTION_EVENTS"))
	ctx = context.WithValue(ctx, models.TestCollKeyMsgs, os.Getenv("TEST_COLLECTION_MSGS"))
	// fmt.Printf("context values - database: %v; collection: %v\n", ctx.Value(models.TestDBKey), ctx.Value(models.TestCollectionKey))

	ac := models.InitializeAppControler(client)

	go ac.ReadEventMessages(ctx)

	// sub1 := models.NewEventSubscriber(models.User{ID: 1, Name: "user1"})

	sub1, event := models.NewEventSubscriberWithEvent("user1")
	if event != nil {
		ac.AcceptEvent(ctx, event)
	}
	sub2, event := models.NewEventSubscriberWithEvent("user2")
	if event != nil {
		ac.AcceptEvent(ctx, event)
	}
	sub3, event := models.NewEventSubscriberWithEvent("user3")
	if event != nil {
		ac.AcceptEvent(ctx, event)
	}

	// add chats between users
	event, err = sub1.CreateEvent(models.CREATE_CHAT, models.NewMessage("Create chat"), sub2)
	if err != nil {
		fmt.Println(err)
	}
	ac.AcceptEvent(ctx, event)
	event, err = sub1.CreateEvent(models.CREATE_CHAT, models.NewMessage("Create chat"), sub3)
	if err != nil {
		fmt.Println(err)
	}
	ac.AcceptEvent(ctx, event)

	// user online event
	for _, user := range sub1.Chats {
		event, err := sub1.CreateEvent(models.USER_ONLINE, models.NewMessage("User online"), &user)
		if err != nil {
			fmt.Println(err)
		}
		ac.AcceptEvent(ctx, event)
	}

	event, err = sub1.CreateEvent(models.MSG_SENT, models.NewMessage("Test message"), sub2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(event)

	// err = ac.AddEvent(ctx, event)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err = ac.SaveNewMessage(ctx, &event.Data)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)

	ac.AcceptEvent(ctx, event)

	for {
		select {
		// case result := <-results:
		// 	fmt.Println(result)
		case <-time.After(5 * time.Second):
			fmt.Println("timed out")
			ac.ControlChan <- models.DoExit
			<-ac.ControlChan // wait for response from go routine (gives the go routine a chance to finish its work)
			fmt.Println("program exit")
			return
		}
	}
	// event, err := ac.GetEventByMessageId(ctx, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// res, _ := json.Marshal(event)
	// fmt.Println("MessageID 0:", string(res))

	// event1, err := ac.GetEventByMessageId(ctx, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// res, _ = json.Marshal(event1)
	// fmt.Println("MessageID 1:", string(res))

	// event2, err := ac.GetEventByMessageId(ctx, 2)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// res, _ = json.Marshal(event2)
	// fmt.Println("MessageID 2:", string(res))

	// count, err := ac.CountMessagesBySubjectId(ctx, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("# messages id 0:", count)

	// count, err = ac.CountMessagesBySubjectId(ctx, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("# messages id 1:", count)

	// user1 := models.User{ID: 1, Name: "user1"}
	// user2 := models.User{ID: 2, Name: "user2"}
	// user3 := models.User{ID: 3, Name: "user3"}

	// count, err := ac.CountMessagesSentByUser(ctx, &user1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("# messages sender id 1:", count)

	// count, err = ac.CountMessagesSentByUser(ctx, &user2)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("# messages sender id 2:", count)

	// count, err = ac.CountMessagesSentByUser(ctx, &user3)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("# messages sender id 3:", count)

	// // Create message
	// msg := models.NewMessage("Test Insert")

	// event, err := models.NewEventSubscriber(user3).CreateEvent(models.MSG_SENT, &msg, models.NewEventSubscriber(user1))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = ac.RemoveEventByMessageId(ctx, msg.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = ac.AddEvent(ctx, event)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// msgUpdate := models.NewMessage("Updated Insert")
	// err = ac.UpdateEventMessageByMessageId(ctx, msg.ID, msgUpdate)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = ac.RemoveEventByMessageId(ctx, msg.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
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
