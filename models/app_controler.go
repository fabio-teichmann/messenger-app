package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppControler struct {
	DB *mongo.Client
	// ESB *EventSubjectBroker
	// channels
	ControlChan chan ControlMsg
	MsgSent     *EventSubject
	MsgRcvd     *EventSubject
	NewUser     *EventSubject // for auto-subscription
	UserLogIn   *EventSubject
	UserOnl     *EventSubject
	UserLogOut  *EventSubject
}

// func NewAppControler(client *mongo.Client, esb *EventSubjectBroker) AppControler {
// 	return AppControler{DB: client, ESB: esb}
// }

func InitializeAppControler(client *mongo.Client) *AppControler {
	msgSent := NewEventSubject(MSG_SENT)
	msgRcvd := NewEventSubject(MSG_RECEIVED)
	newUser := NewEventSubject(NEW_USER)
	userLogIn := NewEventSubject(USER_LOGIN)
	userOnl := NewEventSubject(USER_ONLINE)
	userLogOut := NewEventSubject(USER_LOGOUT)

	return &AppControler{
		DB:          client,
		ControlChan: make(chan ControlMsg),
		MsgSent:     msgSent,
		MsgRcvd:     msgRcvd,
		NewUser:     newUser,
		UserLogIn:   userLogIn,
		UserOnl:     userOnl,
		UserLogOut:  userLogOut,
	}
}

func (ac *AppControler) AcceptEvent(ctx context.Context, event *Event) {
	// persist incoming event
	err := ac.AddEvent(ctx, event)
	if err != nil {
		fmt.Println(err)
	}
	// to avoid blocking when calling AcceptEvent, use go functions
	if event.SubjectID == MSG_SENT {
		go func() { ac.MsgSent.Queue <- *event }()

	} else if event.SubjectID == MSG_RECEIVED {
		go func() { ac.MsgRcvd.Queue <- *event }()

	} else if event.SubjectID == NEW_USER {
		go func() { ac.NewUser.Queue <- *event }()

	} else {
		fmt.Printf("unknown event subject: %v\n", event.SubjectID)
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
				ac.ControlChan <- ExitOK
				return
			}
		case event := <-ac.MsgSent.Queue:
			// notify
			ac.MsgSent.NotifySubscriber(ctx, ac, &event)

		case event := <-ac.MsgRcvd.Queue:
			// notify
			ac.MsgRcvd.NotifySubscriber(ctx, ac, &event)

		case event := <-ac.NewUser.Queue:
			fmt.Printf("Subscribing user %s to channels...\n", event.Sender.Name)
			// subscribe user to all required channels
			ac.MsgSent.AddSubscriber(&event.Sender)
			ac.MsgRcvd.AddSubscriber(&event.Sender)

		case event := <-ac.UserLogIn.Queue:
			ac.UserLogIn.NotifySubscriber(ctx, ac, &event)

		case event := <-ac.UserOnl.Queue:
			ac.UserOnl.NotifySubscriber(ctx, ac, &event)

		case event := <-ac.UserLogOut.Queue:
			ac.UserLogOut.NotifySubscriber(ctx, ac, &event)
		}
	}
}
