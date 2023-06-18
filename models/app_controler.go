package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppControler struct {
	DB  *mongo.Client
	ESB *EventSubjectBroker
	// channels
	ControlChan chan ControlMsg
	MsgSent     *EventSubjectNew
	MsgRcvd     *EventSubjectNew
	NewUser     *EventSubjectNew // auto-subscription
}

func NewAppControler(client *mongo.Client, esb *EventSubjectBroker) AppControler {
	return AppControler{DB: client, ESB: esb}
}

func InitializeAppControler(client *mongo.Client) *AppControler {
	msgSent := NewEventSubject_(MSG_SENT)
	msgRcvd := NewEventSubject_(MSG_RECEIVED)
	newUser := NewEventSubject_(NEW_USER)
	return &AppControler{
		DB:          client,
		ControlChan: make(chan ControlMsg),
		MsgSent:     msgSent,
		MsgRcvd:     msgRcvd,
		NewUser:     newUser,
	}
}

func (ac *AppControler) AcceptEvent(event *Event) {
	if event.SubjectID == MSG_SENT {
		go func() { ac.MsgSent.Queue <- *event }()

	} else if event.SubjectID == MSG_RECEIVED {
		go func() { ac.MsgRcvd.Queue <- *event }()

	} else if event.SubjectID == NEW_USER {
		go func() { ac.NewUser.Queue <- *event }()

	} else {
		fmt.Printf("unknown event subject %v\n", event.SubjectID)
	}
}

func (ac *AppControler) ReadEventMessages(ctx context.Context) {
	fmt.Println("Listening to events...")
	for {
		select {
		case msg := <-ac.ControlChan:
			switch msg {
			case DoExit:
				fmt.Printf("exit read events\n")
				ac.ESB.ControlChan <- ExitOK
				return
			}
		case event := <-ac.MsgSent.Queue:
			// event, err := ac.GetEventById(ctx, eventMsg.ID)
			// if err != nil {
			// 	fmt.Printf("error ocurred reading MsgSent events: %s", err)
			// 	continue
			// }
			// notify
			ac.MsgSent.NotifySubscriber(ctx, ac, &event)
			// ac.ESB.EventSubject.NotifySubscriber(event)
			// fmt.Println(msg)

		case event := <-ac.MsgRcvd.Queue:
			// notify
			ac.MsgRcvd.NotifySubscriber(ctx, ac, &event)

		case event := <-ac.NewUser.Queue:
			fmt.Printf("Subscribing user %s to channels...\n", event.Sender.Name)
			// subscribe user to all required channels
			ac.MsgSent.AddSubscriber(&event.Sender)
			ac.MsgRcvd.AddSubscriber(&event.Sender)
		}

	}
}
