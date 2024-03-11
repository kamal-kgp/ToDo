package main

import (
    "api/models"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

var todosCollection *mongo.Collection

func init() {
	// Connect to MongoDB
	// put password in place of <password> in mongoUri
	mongoUri := "mongodb+srv://kamal-healthifyForUs:<password>@cluster0.oz3s0sq.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}

	// Set up the todos collection
	todosCollection = client.Database("dummy").Collection("ToDo")
}

func main() {
	app := fiber.New()

	//routes
	app.Get("/getTodos", getTodos)
	app.Post("/addTodo", addTodo)

	app.Listen(":8080")
}

func getTodos(c *fiber.Ctx) error {
	cursor, err := todosCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return c.Status(500).SendString("Error fetching todos from MongoDB")
	}
	defer cursor.Close(context.Background())

	var todos []models.Todo
	if err := cursor.All(context.Background(), &todos); err != nil {
		return c.Status(500).SendString("Error decoding todos from MongoDB")
	}

	return c.JSON(todos)
}

func addTodo(c *fiber.Ctx) error {
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString("Error parsing request body")
	}

	// Insert new todo into MongoDB 
	result, err := todosCollection.InsertOne(context.Background(), todo)
	if err != nil {
		return c.Status(500).SendString("Error inserting todo into MongoDB")
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	return c.JSON(todo)
}