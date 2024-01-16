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

// put
func SubscribeMangaById(mangaId string) error{
	objectID, err := primitive.ObjectIDFromHex(mangaId)
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	_,err = collection.Find(context.TODO(), bson.M{"_id": objectID})
	if (err != nil){
		fmt.Println(err)
		return err
	}

	collection = db.Client.Database("manga-tracker").Collection("users")
	userId := "5f563a9da793b25a09529234"
	objectID, err = primitive.ObjectIDFromHex(userId)

	var user *models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id":objectID}).Decode(&user)
	if user.SubscribeList == nil {
		user.SubscribeList = make(map[string][]int)
	}
	user.SubscribeList[mangaId] = []int{}
	update := bson.M{
		"$set": bson.M{
			"subscribeList": user.SubscribeList,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}