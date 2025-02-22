package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-mongo-todos/db"
	"github.com/go-mongo-todos/handlers"
	"github.com/go-mongo-todos/services"
)

type Application struct {
	Models services.Models
}

func main() {

	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		slog.Error("Error connecting to MongoDB, in main.go")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			slog.Error("Error disconnecting from MongoDB, in main.go")
			panic(err)
		}
	}()

	services.New(mongoClient)

	slog.Info("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CreateRouter()))

}
