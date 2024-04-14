package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func UpdateVol(mangaId primitive.ObjectID,vol models.Vol) (*models.Vol,error){
	return nil,nil
}

func DeleteManyVols(mangaID primitive.ObjectID, vols []models.Vol) ([]models.Vol,error) {
    collection := db.Client.Database("manga-tracker").Collection("mangas")

	var volNumbers []string
	for _,element := range vols {
		volNumbers = append(volNumbers, element.Vol)
	}
    filter := bson.M{
        "_id": mangaID,
        "vols": bson.M{"$elemMatch": bson.M{"vol": bson.M{"$in": volNumbers}}},
    }
    result,err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil,err
	}
	fmt.Printf("Deleted %d documents\n", result.DeletedCount)
    return vols,err
}

func DeleteAllVolsByMangaId(mangaId primitive.ObjectID) error{
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