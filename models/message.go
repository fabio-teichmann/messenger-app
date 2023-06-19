package models

import (
	"context"
	"errors"
	"fmt"
	"messenger-app/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	ID         uint32           `bson:"message_id"` // to identify messages in history
	Body       string           `json:"body"`
	Attachment []byte           `json:"attachmet" bson:"omitempty"`
	Sender     *EventSubscriber `bson:"sender"`
	Target     *EventSubscriber `bson:"target"`
	Sent       bool             `bson:"sent"`
	TimeSent   time.Time        `bson:"time_sent"`
	Received   bool             `bson:"received"`
	TimeRcvd   time.Time        `bson:"time_rcvd"`
	Seen       bool             `bson:"seen"`
	TimeSeen   time.Time        `bson:"time_seen"`
}

func NewMessage(body string) *Message {
	hash := util.CreateHash([]byte(body))

	message := Message{
		ID:   hash,
		Body: body,
	}
	return &message
}

func NewMessageWithSender(body string, sender *EventSubscriber) *Message {
	msg := NewMessage(body)
	msg.Sender = sender
	return msg
}

func getCollection(ctx context.Context, ac *AppControler) *mongo.Collection {
	client := ac.DB
	dbKey, collKey := fmt.Sprintf("%v", ctx.Value(TestDBKey)), fmt.Sprintf("%v", ctx.Value(TestCollKeyMsgs))

	return client.Database(dbKey).Collection(collKey)
}

func (ac *AppControler) GetMessageById(ctx context.Context, msgId primitive.ObjectID) (*Message, error) {
	coll := getCollection(ctx, ac)
	results := []Message{}

	filter := bson.M{"_id": msgId}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println("unable to read results")
		return nil, err
	}

	if len(results) > 1 {
		return nil, errors.New("ambiguous id")
	}

	return &results[0], nil

}

func (ac *AppControler) SaveNewMessage(ctx context.Context, msg *Message) error {
	coll := getCollection(ctx, ac)

	result, err := coll.InsertOne(ctx, msg)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with _id: %v", result.InsertedID)
	return nil
}

func (ac *AppControler) SetMessageToSent(ctx context.Context, msgId primitive.ObjectID) error {
	coll := getCollection(ctx, ac)

	filter := bson.M{"_id": msgId}
	update := bson.D{{"$set", bson.D{{"sent", true}, {"time_sent", time.Now()}}}} // update message

	result, err := coll.UpdateByID(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("updated %v event(s)", result.ModifiedCount)
	return nil
}

func (ac *AppControler) SetMessageToRcvd(ctx context.Context, msgId primitive.ObjectID) error {
	coll := getCollection(ctx, ac)

	filter := bson.M{"_id": msgId}
	update := bson.D{{"$set", bson.D{{"received", true}, {"time_rcvd", time.Now()}}}} // update message

	result, err := coll.UpdateByID(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("updated %v message(s)", result.ModifiedCount)
	return nil
}

func (ac *AppControler) SetMessageToSeen(ctx context.Context, msgId primitive.ObjectID) error {
	coll := getCollection(ctx, ac)

	filter := bson.M{"_id": msgId}
	update := bson.D{{"$set", bson.D{{"seen", true}, {"time_seen", time.Now()}}}} // update message

	result, err := coll.UpdateByID(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("updated %v message(s)", result.ModifiedCount)
	return nil
}

func (ac *AppControler) DeleteMessageById(ctx context.Context, msgId primitive.ObjectID) error {
	coll := getCollection(ctx, ac)

	filter := bson.M{"_id": msgId}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Println("deleted message(s):", result.DeletedCount)
	return nil
}
