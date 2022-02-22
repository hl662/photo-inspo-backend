package authentication

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type APIHandler struct {
	MongoClient *mongo.Client
}

func (this *APIHandler) SignupEndpoint(c *gin.Context) {
	var newUser User

	type responseJSON struct {
		UserToken primitive.ObjectID `json:"userToken,omitempty" validate:"required"`
		Username  string             `json:"username,omitempty" validate:"required"`
	}

	err := c.BindJSON(&newUser)
	if err != nil {
		// Send back internal server error with a message
		return
	}
	newUser.Password = EncryptAES(newUser.Password)
	userCollection := this.MongoClient.Database("photoInspo").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		// return bad json.
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, responseJSON{UserToken: result.InsertedID.(primitive.ObjectID),
		Username: newUser.Username})
}
