package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: Error loading .env file")
	}

	// Set up MongoDB connection
	credentials := options.Credential{
		Username:  os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		AuthSource: os.Getenv("MONGO_DB"),
	}
	mongoURL :=  "mongodb://" + os.Getenv("MONGO_URL")
	clientOptions := options.Client().ApplyURI(mongoURL).SetAuth(credentials)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(os.Getenv("MONGO_DB")).Collection("items")

	// Set up Gin router
	r := gin.Default()

	// Define routes
	r.POST("/item", insertItem)
	r.GET("/items", getItems)
	r.DELETE("/item/:id", deleteItem)

	// Run the server
	r.Run(":8080")
}

func insertItem(c *gin.Context) {
	var item map[string]interface{}
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item inserted", "id": result.InsertedID})
}

func getItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching items"})
		return
	}
	defer cursor.Close(ctx)

	var items []bson.M
	if err = cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func deleteItem(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting item"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}