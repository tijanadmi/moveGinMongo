package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoClient combines all collection clients
type MongoClient struct {
	Hall       HallClient
	Movie      MovieClient
	Repertoire RepertoireClient
	Users      UsersClient
}

// NewMongoClient initializes the MongoDB clients and sets up their collections
func NewMongoClient(client *mongo.Client) *MongoClient {
	return &MongoClient{
		Hall: HallClient{
			col: getCollection(client, "halls"),
		},
		Movie: MovieClient{
			col: getCollection(client, "movies"),
		},
		Repertoire: RepertoireClient{
			col: getCollection(client, "repertoires"),
		},
		Users: UsersClient{
			col: getCollection(client, "users"),
		},
	}
}
