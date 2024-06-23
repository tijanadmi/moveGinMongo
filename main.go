package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tijanadmi/moveginmongo/cmd/api"
	db "github.com/tijanadmi/moveginmongo/repository"
	"github.com/tijanadmi/moveginmongo/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	mongoURL = "mongodb://localhost:27017"
)

var client *mongo.Client

func main(){
	config, err := util.LoadConfig(".")
	if err != nil{
		log.Fatal().Err(err).Msg("cannot load config")
	}

	
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	client, err := connectToMongo(config.MongoURL, config.Username, config.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	
	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	mongoClient := db.NewMongoClient(client)
	runGinServer(config, mongoClient)
	
}

func connectToMongo(mongoURL string, username string, password string) (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	// clientOptions.SetAuth(options.Credential{
	// 	Username: username,
	// 	Password: password,
	// })

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		//log.Println("Error connecting:", err)
		log.Fatal().Err(err).Msg("error connecting")
		return nil, err
	}

	return c, nil
}

func runGinServer(config util.Config, store *db.MongoClient) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}