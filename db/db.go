package db

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

func ConnectToMongo() (*mongo.Client, error) {

	// MongoDB connection string
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	mongoURI := os.Getenv("MONGO_DB_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_DB_URI is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		slog.Error("Error connecting to MongoDB")
		return nil, err
	}

	// Check the connection
	log.Println("Connected to MongoDB")

	return client, nil

}
