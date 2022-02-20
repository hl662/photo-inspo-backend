package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

var client *mongo.Client

func SignUpEndpoint(c *gin.Context) {
	type responseJSON struct {
		UserToken primitive.ObjectID `json:"userToken,omitempty" validate:"required"`
	}
	var newUser User

	err := c.BindJSON(&newUser)
	if err != nil {
		// Send back internal server error with a message
		return
	}
	userCollection := client.Database("photoInspo").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		// return bad json.
		return
	}
	c.JSON(200, responseJSON{UserToken: result.InsertedID.(primitive.ObjectID)})
}
func main() {
	atlasURI := fmt.Sprintf("mongodb+srv://%s:%s@photoinspo.pa7r9.mongodb.net/%s?retryWrites=true&w=majority", os.Getenv("mongoUsr"), os.Getenv("mongoPwd"), os.Getenv("mongoDBName"))
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(atlasURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "World",
		})
	})
	router.POST("/signup", SignUpEndpoint)
	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
