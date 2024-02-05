package repo

import (
	"context"
	// "errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
)

// Read
func GetMangas() ([]*models.Manga, error) {
	var mangas []*models.Manga

	collection := db.Client.Database("manga-tracker").Collection("mangas")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &mangas)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
	}
	return mangas, err
}

func GetMangaById(id primitive.ObjectID) ([]*models.Manga, error) {
	var manga []*models.Manga
	
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	cursor, err := collection.Find(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &manga)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
	}
	return manga, err
}

func GetMangaByTitle(title string) ([]*models.Manga, error) {
	var manga []*models.Manga
	
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	cursor, err := collection.Find(context.TODO(), bson.M{"title": title})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &manga)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
	}
	return manga, err
}

func isMangaExist(mangaId primitive.ObjectID) bool {
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	var result *models.Manga
	err := collection.FindOne(context.TODO(), bson.M{"_id": mangaId}).Decode(&result)
	if (err != nil){
		fmt.Println(err)
		return false
	}
	return true
}

// Create
func AddManga(manga *models.Manga) (*models.Manga, error) {
	manga.ID = primitive.NewObjectID()
	_, err := db.Client.Database("manga-tracker").Collection("mangas").InsertOne(context.TODO(), manga)
	if err != nil {
		log.Printf("Couldn't add : %v", err)
		return manga, err
	}
	return manga, nil
}

// Delete
func DeleteMangaByTitle(title string) error {
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"title": title})
	if err != nil {
		log.Printf("Error : %v", err)
	}
	return nil
}

func DeleteMangaById(id primitive.ObjectID) error {
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		log.Printf("Error : %v", err)
	}
	return nil
}

// Update
func UpdateMangaByTitle(manga *models.Manga) (*models.Manga, error) {
	return manga, nil
}

// put
func UpdateMangaSubscriber(mangaId primitive.ObjectID,n int) error {
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	var manga models.Manga
	err := collection.FindOne(context.TODO(), bson.M{"_id":mangaId}).Decode(&manga)
	if err != nil{
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"subscribers": manga.Subscribers + n,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": mangaId}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	return nil
}

func UpdateMangaScore(mangaId primitive.ObjectID,score int) error {
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	var manga models.Manga
	err := collection.FindOne(context.TODO(), bson.M{"_id":mangaId}).Decode(&manga)
	if err != nil{
		return err
	}
	if score == 0{
		manga.TotalVoters -= 1
	}else{
		manga.TotalVoters += 1
	}
	update := bson.M{
		"$set": bson.M{
			"score": (manga.Score+float32(score))/(float32(manga.TotalVoters)),
			"totalVoters": manga.TotalVoters,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": mangaId}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	return nil
}
