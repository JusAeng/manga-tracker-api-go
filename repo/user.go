package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// create
func CreateUser(user *models.User) (*models.User, error) {
	_, err := db.Client.Database("manga-tracker").Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Couldn't create user : %v", err)
		return user, err
	}

	return user, nil
}

// delete
func DeleteUserById(userId primitive.ObjectID) error{
	collection := db.Client.Database("manga-tracker").Collection("users")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": userId})
	if err != nil {
		log.Fatalf("[User] Error to delete user: %v", err)
	}
	if result.DeletedCount == 1 {
		fmt.Println("[User] Delete userId:",userId)
	}
	return nil
}

// read
func GetUserProfileById(userId primitive.ObjectID) (*models.User) {
	var userProfile *models.User
	collection := db.Client.Database("manga-tracker").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&userProfile)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return userProfile
}

// patch
func UpdateUserProfile(userId primitive.ObjectID,key string,newValue string) error{
	fmt.Println("value",key,newValue)
	if (key != "name" && key != "image") {
		return errors.New("This key is not allow to change")
	}
	collection := db.Client.Database("manga-tracker").Collection("users")
	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{key: newValue}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil{
		log.Fatal("Repository: update user name to DB error ",err)
		return err
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
		user.SubscribeList = make([]string,0)
	}
	user.SubscribeList = append(user.SubscribeList, mangaId)
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

func AddMangaVols(userId primitive.ObjectID,mangaId primitive.ObjectID,vols []int) error {
	collection := db.Client.Database("manga-tracker").Collection("users")

	var user *models.User
	err := collection.FindOne(context.TODO(), bson.M{"_id":userId}).Decode(&user)
	if err != nil {
		log.Fatalln("Find user error:",err)
		return err
	}
	return err
	
}