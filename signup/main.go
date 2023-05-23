package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `bson:"id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
}

func dbConn() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://Shubhu123:<Shubhangam@jha123>@cluster0.phue6zi.mongodb.net/?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func getAllUsers(c *gin.Context) {
	client := dbConn()
	defer client.Disconnect(context.Background())

	col := client.Database("clusteralpha").Collection("users")

	cursor, err := col.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	users := []User{}
	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	client := dbConn()
	defer client.Disconnect(context.Background())

	var user User
	err := c.BindJSON(&user)
	if err != nil {
		log.Fatal(err)
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hashedPassword)

	col := client.Database("clusteralpha").Collection("users")
	_, err = col.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Created Successfully"})
}

func main() {
	r := gin.Default()

	r.GET("/users", getAllUsers)
	r.POST("/users", createUser)

	r.Run(":8080")
}
