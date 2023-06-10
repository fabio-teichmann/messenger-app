package tests

import (
	"fmt"
	"messenger-app/models"
	"time"
)

func CreateTestEvent(message string) *models.Event {
	var user1 = models.EventSubscriber{User: models.User{ID: 1, Name: "user1"}}
	var user2 = models.EventSubscriber{User: models.User{ID: 2, Name: "user2"}}

	var esbU1 = models.NewEventSubjectBroker(models.NewEventSubject(1))
	esbU1.EventSubject.AddSubscriber(&user2)

	var esbU2 = models.NewEventSubjectBroker(models.NewEventSubject(2))
	esbU2.EventSubject.AddSubscriber(&user1)

	msg := models.NewMessage(message)
	event, err := user1.CreateEvent(models.MSG_SENT, &msg, &user2)
	if err != nil {
		panic(err)
	}
	return event
}

func RunV1() {
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
			event, err := user3.CreateEvent(models.MSG_SENT, &message, &user2)
			if err != nil {
				fmt.Println(err)
			}
			event.Data.Sent = true

			err = esbU3.AcceptEvent(event)
			fmt.Printf("Event.Message sent: %v\n", event.Data)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println()
			// TODO: update chat
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
