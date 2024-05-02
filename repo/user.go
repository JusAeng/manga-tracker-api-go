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
func RegisterUser(userId primitive.ObjectID,name string,picture string) (*models.User, error){
	collection := db.Client.Database("manga-tracker").Collection("users")
	var user models.User
	user.ID = userId
	user.Name = name
	user.Image = picture
	insertResult,err := collection.InsertOne(context.TODO(),user)
	if err != nil{
		return nil,err
	}
	err = collection.FindOne(context.Background(), bson.M{"_id": insertResult.InsertedID}).Decode(&user)
	if err != nil {
		return nil,err
	}
	
	return &user,nil
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
		return nil
	}
	return userProfile
}

func GetAllUsers() ([]*models.User,error) {
	
	var users []*models.User

	collection := db.Client.Database("manga-tracker").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &users)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
	}
	return users, err
}

// patch
func UpdateUserProfile(userId primitive.ObjectID,key string,newValue string) error{
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
func SubscribeMangaById(userId primitive.ObjectID,mangaId primitive.ObjectID) ([]string,error){
	isMangaExist(mangaId)
	collection := db.Client.Database("manga-tracker").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"_id":userId}).Decode(&user)
	if err != nil{
		fmt.Println("User not exist")
		return []string{},err
	}
	temp := []string{}
	if user.SubscribeList == nil {
		user.SubscribeList = []string{mangaId.Hex()}
	}else{
		for _,e := range user.SubscribeList{
			if e != mangaId.Hex() {
				temp = append(temp, e)
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
			"totalSubscribe": len(user.SubscribeList),
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": userId}, update)
	if err != nil {
		fmt.Println(err)
		return []string{},err
	}
	if len(temp) == len(user.SubscribeList){
		UpdateMangaSubscriber(mangaId,1)
	}else{
		UpdateMangaSubscriber(mangaId,-1)
	}
	
	return user.SubscribeList,nil
}


func UpdateOwnerList(userId primitive.ObjectID,mangaId primitive.ObjectID,vol int) ([]int,error) {
	if !isMangaExist(mangaId){ 
		return []int{}, errors.New("no manga exist")
	}
	collection := db.Client.Database("manga-tracker").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"_id":userId}).Decode(&user)
	if err != nil {
		fmt.Println("User not found eiei")
		return []int{}, err
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
			user.TotalBooks = 1
		}else{
			temp := []int{}
			for _,v := range allvols{
				if v == vol {
					vol = 9999
					user.TotalBooks -= 1
					continue
				}else{
					if vol < v{
						temp = append(temp, vol)
						user.TotalBooks += 1
						vol = 9999
					}
					temp = append(temp, v)
				}
			}
			if vol < 9999 {
				temp = append(temp, vol)
				user.TotalBooks += 1
			}
			
			user.OwnerList[mangaId.Hex()] = temp
		}
	}
	update := bson.M{
		"$set": bson.M{
			"ownerList": user.OwnerList,
			"totalBooks": user.TotalBooks,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": userId}, update)
	if err != nil {
		fmt.Println(err)
		return []int{},err
	}

	return user.OwnerList[mangaId.Hex()],nil
}

func UpdateRateList(userId primitive.ObjectID,mangaId primitive.ObjectID,score int) (error){
	if !isMangaExist(mangaId){ return errors.New("no manga exist")}
	user := GetUserProfileById(userId)
	if user.RateList == nil{
		user.RateList = make(map[string]int)
	}

	user.RateList[mangaId.Hex()] = score
	if score == 0 {
		delete(user.RateList,mangaId.Hex())
	}
	update := bson.M{
		"$set": bson.M{
			"rateList": user.RateList,
		},
	}

	collection := db.Client.Database("manga-tracker").Collection("users")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": userId}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// err = UpdateMangaScore(mangaId,score)
	// if err != nil {
	// 	return errors.New("update fail")
	// }

	return err
}