package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	DBName   string
	User     string
	Password string
	Port     string
}

func NewMongoConnection(config *MongoConfig) (*mongo.Client, func(), error) {
	fmt.Println("Connecting to MongoDB...")
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
		config.User, config.Password, config.DBName,
	)
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// instanciate MongoDB
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		fmt.Println("unable to connect to db")
		return nil, nil, err
	}

	// defer closing
	closer := func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, closer, nil
}
