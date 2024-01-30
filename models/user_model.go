package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      		primitive.ObjectID 	`json:"_id" bson:"_id"`
	Name			string				`json:"name"`
	Image			string				`json:"image"`
	TotalSubscribe	int					`json:"totalSubscribe"`
	TotalBooks		int					`json:"totalBooks"`
	SubscribeList	[]string			`json:"subscribeList"`
	OwnerList		map[string][]int	`json:"ownerList"`
	RateList		map[string]int		`json:"rateList"`
}