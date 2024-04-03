package repo

import (
	"context"
	"errors"
	"strconv"

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

// func AddMangaVol(mangaId primitive.ObjectID,vol models.Vol) (*models.Vol, error) {
// 	var manga models.Manga
// 	collection := db.Client.Database("manga-tracker").Collection("mangas")
// 	err := collection.FindOne(context.TODO(),bson.M{"_id":mangaId}).Decode(&manga)
// 	if err != nil{
// 		return nil,err
// 	}
// 	vol.ID = primitive.NewObjectID()
// 	vol.MangaID = manga.ID.Hex()
// 	temp := []models.Vol{}
// 	if manga.Vols == nil{
// 		manga.Vols = []models.Vol{}
// 		temp = append(temp, vol)
// 	}else{
// 		for idx,v := range manga.Vols{
// 			if v.Vol == vol.Vol{
// 				return nil,errors.New("Already Added")
// 			}
// 			if vol.Vol < v.Vol{
// 				temp = append(temp, vol)
// 				temp = append(temp, manga.Vols[idx:]...)
// 				break
// 			} 
// 			temp = append(temp, v)
// 		}
// 	}
// 	if len(manga.Vols) == len(temp){
// 		temp = append(temp, vol)
// 	}
// 	manga.Vols = temp
// 	manga.LastVol = manga.Vols[len(manga.Vols)-1].Vol
// 	update := bson.M{
// 		"$set": bson.M{
// 			"vols": manga.Vols,
// 			"lastVol": manga.LastVol,
// 		},
// 	}
// 	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": mangaId}, update)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil,err
// 	}
// 	return &vol, nil
// }

func AddMangaVol(mangaId primitive.ObjectID,vol models.Vol) (*models.Vol, error) {
	var manga models.Manga
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	err := collection.FindOne(context.TODO(),bson.M{"_id":mangaId}).Decode(&manga)
	if err != nil{
		return nil,err
	}
	if manga.Vols == nil{
		manga.Vols = make(map[string]models.Vol)
	}
	if _, exist := manga.Vols[vol.Vol]; exist {
		return nil,errors.New("this volumn already exist")
    } else {
        manga.Vols[vol.Vol] = vol
    }
	vol.MangaID = manga.ID.Hex()
	lastest := float64(0)
	for key := range manga.Vols {        
		floatValue,err := strconv.ParseFloat(key,64)
		if err != nil{
			return nil,errors.New("can't parse to float")
		}
		if floatValue > lastest{
			lastest = floatValue
		}
    }
	fmt.Println(lastest)
	manga.LastVol = fmt.Sprintf("%.1f", lastest)

	update := bson.M{
		"$set": bson.M{
			"vols": manga.Vols,
			"lastVol": manga.LastVol,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": mangaId}, update)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	return &vol, nil
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

func DeleteAllVolsById(mangaId primitive.ObjectID) error{
	var manga models.Manga
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	err := collection.FindOne(context.TODO(),bson.M{"_id":mangaId}).Decode(&manga)
	if err != nil{
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"vols": nil,
			"lastVol": "0",
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": mangaId}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Update
func UpdateManga(manga *models.Manga) (*models.Manga, error) {
	var existManga models.Manga
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	err := collection.FindOne(context.TODO(), bson.M{"_id":manga.ID}).Decode(&existManga)
	if err != nil{
		return nil,err
	}
	update := bson.M{
		"$set": bson.M{
            "title":          manga.Title,
            "otherTitles":    manga.OtherTitle,
            "author":         manga.Author,
            "otherParticipate": manga.OtherParticipate,
            "genre":          manga.Genre,
            "otherGenres":    manga.OtherGeres,
            "image":          manga.Image,
            "introduction":   manga.Introduction,
            "publisher":      manga.Publisher,
            "firstDateJP":    manga.FirstDateJP,
            "firstDateTH":    manga.FirstDateTH,
            "vols":           manga.Vols,
            "lastVol":        manga.LastVol,
        },
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": manga.ID}, update)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
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
