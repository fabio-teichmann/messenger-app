package main

import (
	"fmt"
	"messenger-app/models"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// initiate control channel for graceful shutdown
	// controlChan := make(chan models.ControlMsg, 2)

	var user1 = models.EventSubscriber{User: models.User{ID: 1, Name: "user1"}}
	var user2 = models.EventSubscriber{User: models.User{ID: 2, Name: "user2"}}
	var user3 = models.EventSubscriber{User: models.User{ID: 3, Name: "user3"}}

	var esbU1 = models.NewEventSubjectBroker(models.NewEventSubject(1))
	esbU1.EventSubject.AddSubscriber(&user2)

	var esbU2 = models.NewEventSubjectBroker(models.NewEventSubject(2))
	esbU2.EventSubject.AddSubscriber(&user1)
	esbU2.EventSubject.AddSubscriber(&user3)

	var esbU3 = models.NewEventSubjectBroker(models.NewEventSubject(3))
	esbU3.EventSubject.AddSubscriber(&user2)

	streams := []models.EventSubjectBroker{esbU1, esbU2, esbU3}
	for i := range streams {
		// wg.Add(1)
		go func(stream models.EventSubjectBroker) {

			stream.ReadEvents()
			// wg.Done()
		}(streams[i])
	}

	for _, i := range []int{1, 2, 3, 1, 2, 3} {
		message := models.NewMessage(fmt.Sprintf("%s_%v", "test_message", i))
		go func() {
			event, err := user3.User.CreateEvent(models.MsgSent, &message, &user2.User)
			if err != nil {
				fmt.Println(err)
			}
			esbU3.AcceptEvent(event)
			// event.SendToChat(chat12)
		}()

		// go func() {
		// 	event, err := user2.User.CreateEventMessage(&message, &user1.User)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	esbU2.AcceptEvent(event)
		// }()

		// if i != 2 {
		// 	event2, err := user1.User.CreateEventMessage(&message, &user3.User)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	event2.SendToChat(chat13)
		// }
		// event.SendToChat(chat12)
		time.Sleep(1 * time.Second)
	}

	for {
		select {
		case <-time.After(7 * time.Second):
			fmt.Println("Timed out...")
			for _, stream := range streams {
				stream.ControlChan <- models.DoExit
				<-stream.ControlChan
			}
			// controlChan <- models.DoExit
			// <-controlChan
			fmt.Println("Exit program")
			return
		}
	}

}

// func main() {
// 	// initiate control channel for graceful shutdown
// 	controlChan := make(chan models.ControlMsg, 5)

// 	var user1 = models.EventSubscriber{User: models.User{ID: 1, Name: "user1"}}
// 	var user2 = models.EventSubscriber{User: models.User{ID: 2, Name: "user2"}}
// 	var user3 = models.EventSubscriber{User: models.User{ID: 3, Name: "user3"}}

// 	var chat12 = user1.CreateChat([]*models.EventSubscriber{&user2})
// 	var chat13 = user1.CreateChat([]*models.EventSubscriber{&user3})

// 	// listen for messages
// 	go func() {
// 		chat13.ReadMessages(controlChan)
// 		chat12.ReadMessages(controlChan)
// 	}()

// 	// send test messages
// 	for _, i := range []int{1, 2, 3, 1, 2, 3} {
// 		message := models.NewMessage(fmt.Sprintf("%s_%v", "test_message", i))
// 		go func() {
// 			event, err := user1.User.CreateEventMessage(&message, &user2.User)
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			event.SendToChat(chat12)
// 		}()

// 		if i != 2 {
// 			event2, err := user1.User.CreateEventMessage(&message, &user3.User)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			event2.SendToChat(chat13)
// 		}
// 		// event.SendToChat(chat12)
// 		time.Sleep(1 * time.Second)
// 	}

// 	for {
// 		select {
// 		case <-time.After(7 * time.Second):
// 			fmt.Println("Timed out...")
// 			controlChan <- models.DoExit
// 			<-controlChan
// 			fmt.Println("Exit program")
// 			return
// 		}
// 	}
// }
