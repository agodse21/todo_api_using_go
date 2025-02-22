package services

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Task      string    `bson:"task,omitempty" json:"task,omitempty" validate:"required"`
	Completed bool      `bson:"completed" json:"completed"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

var Client *mongo.Client

func New(mongo *mongo.Client) Todo {

	Client = mongo

	return Todo{}
}

func returnCollectionPointer() *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}
	return Client.Database(os.Getenv("MONGO_DB_NAME")).Collection("todos")
}

func (t *Todo) InsertTodo(body Todo) error {
	collection := returnCollectionPointer()

	_, err := collection.InsertOne(context.TODO(), Todo{
		Task:      body.Task,
		Completed: body.Completed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		slog.Error("Error inserting todo")
		return err
	}

	return nil

}

func (t *Todo) GetAllTodos() ([]Todo, error) {
	collection := returnCollectionPointer()

	var todos []Todo
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		slog.Error("Error finding todos")
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		err := cursor.Decode(&todo)
		if err != nil {
			slog.Error("Error decoding todo")
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil

}

func (t *Todo) GetTodoByID(id string) (Todo, error) {
	collection := returnCollectionPointer()

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("Invalid ID format")
		return Todo{}, err
	}

	var todo Todo
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		slog.Error("Error finding todo")
		return Todo{}, err
	}

	return todo, nil

}

func (t *Todo) UpdateTodoByID(id string, body Todo) error {

	collection := returnCollectionPointer()

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("Invalid ID format")
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"task":       body.Task,
			"completed":  body.Completed,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		slog.Error("Error updating todo")
		return err
	}

	return nil
}

func (t *Todo) DeleteTodoByID(id string) error {

	collection := returnCollectionPointer()

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("Invalid ID format")
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		slog.Error("Error deleting todo")
		return err
	}

	return nil
}
