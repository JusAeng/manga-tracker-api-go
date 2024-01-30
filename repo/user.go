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
func SubscribeMangaById(userId primitive.ObjectID,mangaId primitive.ObjectID) error{
	if !isMangaExist(mangaId){ return errors.New("No manga exist")}

	collection := db.Client.Database("manga-tracker").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"_id":userId}).Decode(&user)
	if err != nil{
		fmt.Println("User not exist")
		return err
	}
	if user.SubscribeList == nil {
		user.SubscribeList = make([]string,0)
		user.SubscribeList = append(user.SubscribeList, mangaId.Hex())
	}else{
		var temp []string
		for _,v := range user.SubscribeList{
			if v != mangaId.Hex() {
				temp = append(temp, v)
			}
		}
		if len(temp) == len(user.SubscribeList){
			temp = append(temp, mangaId.Hex())
		}
		user.SubscribeList = temp
	}
	update := bson.M{
		"$set": bson.M{
			"subscribeList": user.SubscribeList,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": userId}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	return nil
}


func UpdateOwnerList(userId primitive.ObjectID,mangaId primitive.ObjectID,vol int) error {
	if !isMangaExist(mangaId){ return errors.New("No manga exist")}
	collection := db.Client.Database("manga-tracker").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"_id":userId}).Decode(&user)
	if err != nil {
		fmt.Println("User not found eiei")
		return err
	}
	vols := []int{vol}
	if user.OwnerList == nil {
		user.OwnerList = map[string][]int{
			mangaId.Hex(): vols,
		}
	}else{
		allvols ,exist := user.OwnerList[mangaId.Hex()]
		if !exist {
			user.OwnerList[mangaId.Hex()] = vols
		}else{
			temp := []int{}
			for _,v := range allvols{
				if v != vol{
					temp = append(temp, v)
				}
			}
			if len(temp) == len(allvols){
				temp = append(temp, vol)
			}
			user.OwnerList[mangaId.Hex()] = temp
		}
	}
	update := bson.M{
		"$set": bson.M{
			"ownerList": user.OwnerList,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": userId}, update)
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