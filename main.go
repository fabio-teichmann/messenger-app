package main

import (
	"fmt"
	"messenger-app/models"
	"time"
)

func main() {
	// initiate control channel for graceful shutdown
	controlChan := make(chan models.ControlMsg, 5)

	var user1 = models.EventSubscriber{User: models.User{ID: 1, Name: "user1"}}
	var user2 = models.EventSubscriber{User: models.User{ID: 2, Name: "user2"}}
	var user3 = models.EventSubscriber{User: models.User{ID: 3, Name: "user3"}}

	var chat12 = user1.CreateChat([]*models.EventSubscriber{&user2})
	var chat13 = user1.CreateChat([]*models.EventSubscriber{&user3})

	// listen for messages
	go func() {
		chat13.ReadMessages(controlChan)
		chat12.ReadMessages(controlChan)
	}()

	// send test messages
	for _, i := range []int{1, 2, 3, 1, 2, 3} {
		message := models.NewMessage(fmt.Sprintf("%s_%v", "test_message", i))
		go func() {
			event, err := user1.User.CreateEventMessage(&message, &user2.User)
			if err != nil {
				fmt.Println(err)
			}

			event.SendToChat(chat12)
		}()

		if i != 2 {
			event2, err := user1.User.CreateEventMessage(&message, &user3.User)
			if err != nil {
				fmt.Println(err)
			}
			event2.SendToChat(chat13)
		}
		// event.SendToChat(chat12)
		time.Sleep(1 * time.Second)
	}

	for {
		select {
		case <-time.After(7 * time.Second):
			fmt.Println("Timed out...")
			controlChan <- models.DoExit
			<-controlChan
			fmt.Println("Exit program")
			return
		}
	}
}
