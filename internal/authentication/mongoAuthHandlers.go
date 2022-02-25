package authentication

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserToken struct {
	TokenID  primitive.ObjectID `json:"tokenID,omitempty" validate:"required"`
	Username string             `json:"username,omitempty" validate:"required"`
}
type APIHandler struct {
	MongoClient *mongo.Client
}

func (this *APIHandler) SignupEndpoint(c *gin.Context) {
	var newUser User

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	userCollection := this.MongoClient.Database("photoInspo").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Checks if there already contains a user with this username.
	foundUserCursor := userCollection.FindOne(ctx, bson.M{"username": newUser.Username})
	if foundUserCursor.Err() == nil {
		c.JSON(500, gin.H{"error": "Username already exists."})
	}

	newUser.Password = EncryptAES(newUser.Password)
	result, err := userCollection.InsertOne(ctx, newUser)
	c.JSON(200, UserToken{TokenID: result.InsertedID.(primitive.ObjectID),
		Username: newUser.Username})
}

func (this *APIHandler) SigninEndpoint(c *gin.Context) {
	type mongoResult struct {
		ID       primitive.ObjectID `bson:"_id, omitempty"`
		Username string             `json:"username,omitempty" validate:"required"`
		Password string             `json:"password,omitempty" validate:"required"`
	}
	var requestUser = User{
		Username: c.Query("username"),
		Password: c.Query("password"),
	}
	userCollection := this.MongoClient.Database("photoInspo").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	foundUserCursor := userCollection.FindOne(ctx, bson.M{"username": requestUser.Username})
	if foundUserCursor.Err() != nil {
		c.JSON(500, gin.H{"error": foundUserCursor.Err()})
	}
	var foundUser mongoResult
	err := foundUserCursor.Decode(&foundUser)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}
	if DecryptAES(foundUser.Password) != requestUser.Password {
		c.JSON(400, gin.H{"error": "Error, Incorrect Password"})
	}
	c.JSON(200, UserToken{TokenID: foundUser.ID,
		Username: foundUser.Username})
}
