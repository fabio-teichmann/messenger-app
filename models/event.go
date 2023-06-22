package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID        primitive.ObjectID `bson:"_id"`
	SubjectID EventType          `bson:"subject_id"`      // on which queue to publish the message
	Sender    EventSubscriber    `bson:"sender"`          // event origin
	Target    EventSubscriber    `bson:"target"`          // event destination
	Data      Message            `bson:"data, omitempty"` // contains payload
	Time      time.Time          `bson:"time"`            // when the event was published
	// EventType EventType          `bson:"event_type"` // to classify events
}

// func (e *Event) SendToChat(chat Chat) {
// 	chat.Chat <- *e
// 	chat.History = append(chat.History, e.Data)
// }

func (e *Event) SentToRcvd() {
	// change Subject
	e.SubjectID = MSG_RECEIVED
	// switch sender and target
	temp := e.Sender
	e.Sender = e.Target
	e.Target = temp
	e.Time = time.Now()
}

func (ac *AppControler) GetEventByMessageId(ctx context.Context, messageId uint32) (*Event, error) {
	results := []Event{}

	client := ac.DB
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)

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

func (ac *AppControler) GetEventById(ctx context.Context, eventId primitive.ObjectID) (*Event, error) {
	results := []Event{}

	client := ac.DB
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))
	fmt.Println(eventId)
	coll := client.Database(dbKey).Collection(collKey)

	// filter := bson.D{{"subject_id"}}
	filter := bson.M{"_id": bson.M{"$eq": eventId}}
	cursor, err := coll.Find(ctx, filter)
	fmt.Println(cursor)
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
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)

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
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)

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
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)

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

	filter := bson.D{{Key: "data.message_id", Value: msgId}}

	result, err := coll.DeleteOne(ctx, filter)
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

func (ac *AppControler) UpdateEventMessageToSentById(ctx context.Context, eventId primitive.ObjectID) error {
	client := ac.DB
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)
	fmt.Println(eventId)
	// event, err := ac.GetEventById(ctx, eventId)
	// if err != nil {
	// 	return err
	// }

	filter := bson.D{{Key: "_id", Value: eventId}}
	update := bson.D{{"$set", bson.D{{"data.sent", true}, {"data.time_sent", time.Now()}}}} // replace the message within event

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("updated %v event(s)\n", result.ModifiedCount)
	return nil
}

func (ac *AppControler) UpdateEventMessageToRcvdById(ctx context.Context, eventId primitive.ObjectID) error {
	client := ac.DB
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollectionKey))

	coll := client.Database(dbKey).Collection(collKey)

	event, err := ac.GetEventById(ctx, eventId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: event.ID}}
	update := bson.D{{"$set", bson.D{{"data.received", true}, {"data.time_rcvd", time.Now()}}}} // replace the message within event

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("updated %v event(s)\n", result.ModifiedCount)
	return nil
}
