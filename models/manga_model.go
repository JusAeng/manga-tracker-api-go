package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manga struct {
	ID      			primitive.ObjectID 	`json:"_id" bson:"_id"`
	Title   			string  			`json:"title"`
	OtherTitle			[]string			`json:"otherTitles" bson:"otherTitles"`
	Author				string				`json:"author"`
	OtherParticipate	[]string			`json:"otherParticipate" bson:"otherParticipate"`
	Genre				string				`json:"genre"`
	OtherGeres			[]string			`json:"otherGenres" bson:"otherGenres"`
	Image				string				`json:"image"`
	Introduction		string				`json:"introduction"`
	Publisher			string				`json:"publisher"`
	FirstDateJP			string				`json:"firstDateJP" bson:"firstDateJP"`
	FirstDateTH			string				`json:"firstDateTH" bson:"firstDateTH"`
	Vols				map[string]Vol		`json:"vols"`
	LastVol				string				`json:"lastVol" bson:"lastVol"`
	Subscribers			int					`json:"subscribers"`
	Score				float32				`json:"score"`
	TotalVoters			int 				`json:"totalVoters" bson:"totalVoters"`
}

type Vol struct {
	MangaID			string 				`json:"mangaId"`
	Vol 			string				`json:"vol"`
	Image			string				`json:"image"`
	PublishDate		string				`json:"publishDate"`
	TotalOwner		int					`json:"totalOwner"`
}