package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manga struct {
	ID      			primitive.ObjectID 	`json:"_id" bson:"_id"`
	Title   			string  			`json:"title"`
	OtherTitle			[]string			`json:"otherTitles"`
	Author				string				`json:"author"`
	Drawer				string				`json:"drawer"`
	OtherParticipate	[]string			`json:"otherParticipate"`
	Genre				string				`json:"genre"`
	OtherGeres			[]string			`json:"otherGenres"`
	Image				string				`json:"image"`
	Introduction		string				`json:"introduction"`
	Publisher			string				`json:"publisher"`
	FirstDateJP			string				`json:"firstDateJP"`
	FirstDateTH			string				`json:"firstDateTH"`
	Vols				[]int				`json:"vols"`
	LastVol				int					`json:"lastVol"`
	Subscribers			int					`json:"subscribers"`
	Score				float32				`json:"score"`
	TotalVoters			int 				`json:"totalVoters"`
}

type Vols struct {
	ID      		primitive.ObjectID 	`json:"_id" bson:"_id"`
	Vol 			int					`json:"vol"`
	MangaID			string 				`json:"mangaId"`
	Image			string				`json:"image"`
	PublishDate		string				`json:"publishDate"`
	TotalOwner		int					`json:"totalOwner"`
}