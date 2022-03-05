package authentication

import (
	"context"
	"fmt"
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

type Photo struct {
	Src             string `json:"src,omitempty"`
	AltText         string `json:"altText,omitempty"`
	Photographer    string `json:"photographer,omitempty"`
	PhotographerURL string `json:"photographerURL,omitempty"`
	Id              string `json:"id,omitempty"`
}
type Moodboard struct {
	UserName       string  `json:"username,omitempty" validate:"required"`
	Name           string  `json:"name,omitempty" validate:"required"`
	Images         []Photo `json:"images,omitempty" validate:"required"`
	DefaultImageId string  `json:"defaultImageId,omitempty"`
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
		c.JSON(400, gin.H{"error": "Username already exists."})
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
		c.JSON(400, gin.H{"error": "No user found."})
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

func (this *APIHandler) SaveEndpoint(c *gin.Context) {
	type RequestJSON struct {
		UserName       string  `json:"username,omitempty" validate:"required"`
		Name           string  `json:"name,omitempty" validate:"required"`
		Images         []Photo `json:"images,omitempty" validate:"required"`
		DefaultImageId string  `json:"defaultImageId,omitempty"`
	}
	var newMoodboardJSON RequestJSON
	var newMoodboard Moodboard
	err := c.BindJSON(&newMoodboardJSON)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	newMoodboard = Moodboard{
		UserName:       newMoodboardJSON.UserName,
		Name:           newMoodboardJSON.Name,
		Images:         newMoodboardJSON.Images,
		DefaultImageId: newMoodboardJSON.DefaultImageId,
	}
	moodboards := this.MongoClient.Database("photoInspo").Collection("Moodboards")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := moodboards.CountDocuments(ctx, bson.M{"username": newMoodboard.UserName, "name": newMoodboard.Name})
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	if count >= 1 {
		c.JSON(400, gin.H{"error": "Moodboard already exists."})
		c.Abort()
		return
	}

	_, err = moodboards.InsertOne(ctx, newMoodboard)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	c.JSON(200, gin.H{"message": "success!"})
}

func (this *APIHandler) UpdateEndpoint(c *gin.Context) {
	type RequestJSON struct {
		UserName       string  `json:"username,omitempty" validate:"required"`
		OldName        string  `json:"oldName,omitempty" validate:"required"`
		Name           string  `json:"name,omitempty" validate:"required"`
		Images         []Photo `json:"images,omitempty" validate:"required"`
		DefaultImageId string  `json:"defaultImageId,omitempty"`
	}
	var newMoodboardJSON RequestJSON
	var newMoodboard Moodboard
	err := c.BindJSON(&newMoodboardJSON)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	newMoodboard = Moodboard{
		UserName:       newMoodboardJSON.UserName,
		Name:           newMoodboardJSON.Name,
		Images:         newMoodboardJSON.Images,
		DefaultImageId: newMoodboardJSON.DefaultImageId,
	}
	moodboards := this.MongoClient.Database("photoInspo").Collection("Moodboards")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Must update the document.
	filter := bson.M{"name": newMoodboardJSON.OldName, "username": newMoodboard.UserName}
	update := bson.D{{"$set", bson.D{
		{"name", newMoodboard.Name},
		{"images", newMoodboard.Images},
		{"defaultimageid", newMoodboard.DefaultImageId},
	}}}
	_, err = moodboards.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "success!"})
}

func (this *APIHandler) GetMoodboardsEndpoint(c *gin.Context) {
	type ResponseJSON struct {
		Count      int         `json:"count"`
		Moodboards []Moodboard `json:"moodboards"`
	}

	username := c.Query("username")
	moodboards := this.MongoClient.Database("photoInspo").Collection("Moodboards")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var resultMoodboards []Moodboard
	filter := bson.M{"username": username}

	cursor, err := moodboards.Find(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	for cursor.Next(ctx) {
		var moodboard Moodboard
		err := cursor.Decode(&moodboard)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		resultMoodboards = append(resultMoodboards, moodboard)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	err = cursor.Close(ctx)
	if err != nil {
		return
	}
	c.JSON(200, ResponseJSON{
		Count:      len(resultMoodboards),
		Moodboards: resultMoodboards,
	})
}

func (this *APIHandler) DeleteEndpoint(c *gin.Context) {
	type RequestJSON struct {
		UserName string `json:"username,omitempty" validate:"required"`
		Name     string `json:"name,omitempty" validate:"required"`
	}
	var request RequestJSON
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	moodboards := this.MongoClient.Database("photoInspo").Collection("Moodboards")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"name": request.Name, "username": request.UserName}
	res, err := moodboards.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if res.DeletedCount == 0 {
		c.JSON(400, gin.H{"error": "Moodboard not found."})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message": fmt.Sprintf("Deleted %d moodboard.", res.DeletedCount)})
}
