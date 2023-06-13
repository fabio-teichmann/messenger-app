package models

import "go.mongodb.org/mongo-driver/mongo"

type AppControler struct {
	DB *mongo.Client
}

func NewAppControler(client *mongo.Client) AppControler {
	return AppControler{DB: client}
}
