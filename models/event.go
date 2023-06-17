package models

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID        primitive.ObjectID `bson:"_id"`
	SubjectID int                `bson:"subject_id"` // on which queue to publish the message
	Sender    EventSubscriber    `bson:"sender"`     // event origin
	Target    EventSubscriber    `bson:"target"`     // event destination
	Data      Message            `bson:"data"`       // contains payload
	EventType EventType          `bson:"event_type"` // to classify events
}

func (e *Event) SendToChat(chat Chat) {
	chat.Chat <- *e
	chat.History = append(chat.History, e.Data)
}

func (ac *AppControler) GetEventByMessageId(ctx context.Context, messageId uint32) (*Event, error) {
	results := []Event{}

	client := ac.DB

	coll := client.Database(ctx.Value(TestDBKey).(string)).Collection(ctx.Value(TestCollectionKey).(string))

	// filter := bson.D{{"subject_id"}}
	cursor, err := coll.Find(ctx, bson.D{{Key: "data.message_id", Value: messageId}})
	if err != nil {
		fmt.Println("no message for given id")
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		fmt.Print("unable to read results")
		return nil, err
	}

	if len(results) > 1 {
		fmt.Println("ambiguous id")
		return nil, errors.New("ambiguous message id for events")
	}

	return &results[0], nil
}

func (ac *AppControler) CountMessagesBySubjectId(ctx context.Context, subjectId int) (int, error) {

	client := ac.DB

	coll := client.Database(ctx.Value(TestDBKey).(string)).Collection(ctx.Value(TestCollectionKey).(string))

	filter := bson.D{{Key: "subject_id", Value: subjectId}}
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return -1, err
	}
	return int(count), nil
}

func (ac *AppControler) CountMessagesSentByUser(ctx context.Context, user *User) (int, error) {
	results := []Event{}

	client := ac.DB

	coll := client.Database(ctx.Value(TestDBKey).(string)).Collection(ctx.Value(TestCollectionKey).(string))

	filter := bson.D{{Key: "sender.user.id", Value: user.ID}}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return -1, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return -1, err
	}

	return len(results), nil
}

func (ac *AppControler) AddEvent(ctx context.Context, event *Event) error {

	client := ac.DB
	coll := client.Database(ctx.Value(TestDBKey).(string)).Collection(ctx.Value(TestCollectionKey).(string))

	result, err := coll.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}

func (ac *AppControler) RemoveEventByMessageId(ctx context.Context, msgId uint32) error {
	client := ac.DB
	coll := client.Database(ctx.Value(TestDBKey).(string)).Collection(ctx.Value(TestCollectionKey).(string))

	result, err := coll.DeleteOne(ctx, bson.D{{Key: "data.message_id", Value: msgId}})
	if err != nil {
		return err
	}
	fmt.Printf("Removed %v documents\n", result.DeletedCount)
	return nil
}

func (ac *AppControler) UpdateEventMessageByMessageId(ctx context.Context, msgId uint32, msg Message) error {
	client := ac.DB
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)

	event, err := ac.GetEventByMessageId(ctx, msgId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: event.ID}}
	update := bson.D{{"$set", bson.D{{"data", msg}}}} // replace the message within event

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("updated %v event(s)", result.ModifiedCount)
	return nil
}
