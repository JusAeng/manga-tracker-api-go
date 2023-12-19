package repo

import (
	"context"
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

// Update
func UpdateMangaByTitle(manga *models.Manga) (*models.Manga, error) {
	return manga, nil
}
