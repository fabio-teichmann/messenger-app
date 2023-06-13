package models

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// EventType
type EventType int

// event types
const (
	USER_ONLINE = iota
	USER_TYPING
	MSG_SENT
	MSG_RECEIVED
)

type Event struct {
	SubjectID int             `bson:"subject_id"`  // on which queue to publish the message
	Sender    EventSubscriber `bson:"sender"`      // event origin
	Target    EventSubscriber `bson:"target"`      // event destination
	Data      Message         `bson:"data,inline"` // contains payload
	EventType EventType       `bson:"event_type"`  // to classify events
}

func (e *Event) SendToChat(chat Chat) {
	chat.Chat <- *e
	chat.History = append(chat.History, e.Data)
}

func (ac AppControler) GetEventByMessageID(ctx context.Context, messageId int) (*Event, error) {
	results := []Event{}

	client := ac.DB

	coll := client.Database("messenger-test").Collection("events")

	// filter := bson.D{{"subject_id"}}
	cur, err := coll.Find(ctx, bson.D{{"message_id", messageId}})
	if err != nil {
		fmt.Println("no message for given id")
		return nil, err
	}

	if err = cur.All(ctx, &results); err != nil {
		fmt.Print("unable to read results")
		return nil, err
	}

	if len(results) > 1 {
		fmt.Println("ambiguous id")
		return nil, errors.New("ambiguous message id for events")
	}

	return &results[0], nil
}
