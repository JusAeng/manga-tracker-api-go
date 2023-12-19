package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manga struct {
	ID      primitive.ObjectID 	`json:"_id" bson:"_id"`
	Title   string  			`json:"title"`
}