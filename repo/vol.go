package repo

import (
	"context"
	"errors"
	"fmt"

	// "strconv"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddMangaVol(vol models.Vol) (*models.Vol, error) {
	var manga models.Manga
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	err := collection.FindOne(context.TODO(),bson.M{"_id":vol.MangaID}).Decode(&manga)
	if err != nil{
		return nil,err
	}
	lastest := vol.Vol
	if manga.Vols == nil{
		manga.Vols = make([]models.Vol,0)
	} else{
		for _,e := range manga.Vols {
			if e.Vol == vol.Vol {
				return nil,errors.New("already add this vol")
			}
			if e.Vol > lastest{
				lastest = e.Vol
			}
		}
	}
	manga.Vols = append(manga.Vols, vol)
	manga.LastVol = lastest
	update := bson.M{
		"$set": bson.M{
			"vols": manga.Vols,
			"lastVol": manga.LastVol,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": vol.MangaID}, update)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	return &vol, nil
}

func UpdateVol(mangaId primitive.ObjectID,req models.Vol) (error){
	collection := db.Client.Database("manga-tracker").Collection("mangas")
	filter := bson.M{
        "_id": 		mangaId,
        "vols.vol": req.Vol,
    }
    update := bson.M{
        "$set": bson.M{
            "vols.$.image":       req.Image,
            "vols.$.publishDate": req.PublishDate,
        },
    }
	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return errors.New("can't update this vol")
	}
	fmt.Printf("Update %d documents\n", result.ModifiedCount)
	return nil
}

func DeleteManyVols(mangaID primitive.ObjectID, volNumbers []int) ([]int,error) {
    collection := db.Client.Database("manga-tracker").Collection("mangas")
	fmt.Println("look: ",mangaID,volNumbers)
	filter := bson.M{"_id": mangaID}
	update := bson.M{"$pull": bson.M{"vols": bson.M{"vol": bson.M{"$in": volNumbers}}}}

    result,err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return nil,err
	}
	fmt.Printf("Update %d documents\n", result.ModifiedCount)
    return volNumbers,err
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