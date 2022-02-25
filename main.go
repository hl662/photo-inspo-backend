package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hl662/photo-inspo-backend/internal/authentication"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var err error
	var client *mongo.Client
	router := gin.Default()
	atlasURI := fmt.Sprintf("mongodb+srv://%s:%s@photoinspo.pa7r9.mongodb.net/%s?retryWrites=true&w=majority", os.Getenv("mongoUsr"), os.Getenv("mongoPwd"), os.Getenv("mongoDBName"))

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
	router.Use(CORSMiddleware())
	apiHandler := new(authentication.APIHandler)
	apiHandler.MongoClient = client
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "World",
		})
	})
	router.POST("/signup", apiHandler.SignupEndpoint)
	router.GET("/login", apiHandler.SigninEndpoint)

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
