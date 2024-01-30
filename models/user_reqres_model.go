package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddOwnerListReq struct{
	MangaId			primitive.ObjectID	`json:"mangaId"`
	Vol			int						`json:"vol"`	
}