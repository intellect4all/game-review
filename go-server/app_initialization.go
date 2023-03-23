package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"

	"log"
	"time"
)

const (
	// ExecutionMode is the environment variable that determines the execution mode of the application.
	// Possible values are "local" and "lambda".
	// if the value is "local", the application will run locally and not in lambda.
	ExecutionMode = "EXECUTION_MODE"
	MongoDBURI    = "MONGODB_URI"
)

type InitializationResponse struct {
	MongoDbClient *mongo.Client
}

func InitializationHandler() (*InitializationResponse, error) {
	loadEnvVariables()

	response := &InitializationResponse{}

	mongoURI := os.Getenv(MongoDBURI)

	mongoClient, err := getMongoClient(mongoURI)
	if err != nil {
		return nil, err
	}

	response.MongoDbClient = mongoClient
	return response, nil

}

func getMongoClient(uri string) (*mongo.Client, error) {
	log.Println("Connecting to MongoDB: ")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	log.Println("mongo connected")

	return client, nil
}

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("unable to load env")
	}

	log.Println("env loaded")

}
