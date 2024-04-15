package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      		primitive.ObjectID 	`json:"_id" bson:"_id"`
	Name			string				`json:"name"`
	Image			string				`json:"image"`
	TotalSubscribe	int					`json:"totalSubscribe" bson:"totalSubscribe"`
	TotalBooks		int					`json:"totalBooks" bson:"totalBooks"`
	SubscribeList	[]string			`json:"subscribeList" bson:"subscribeList"`
	OwnerList		map[string][]int	`json:"ownerList" bson:"ownerList"`
	RateList		map[string]int		`json:"rateList" bson:"rateList"`
}