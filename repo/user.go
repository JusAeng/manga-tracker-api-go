package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// read
func GetUserProfile(userId string) (*models.User) {
	var userProfile *models.User
	objectID, err := primitive.ObjectIDFromHex(userId)
	
	collection := db.Client.Database("manga-tracker").Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&userProfile)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return userProfile
}

// create
func AddUser(user *models.User) (*models.User, error) {
	_, err := db.Client.Database("manga-tracker").Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Couldn't create user : %v", err)
		return user, err
	}

	return user, nil
}

// delete
func DeleteUserById(userId string) error{
	collection := db.Client.Database("manga-tracker").Collection("users")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": userId})
	if err != nil {
		log.Printf("Error : %v", err)
	}
	return nil
}