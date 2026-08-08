package models

import (
	"github.com/google/uuid"
)

type Manga struct {
	ID                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	OtherTitles       []string  `json:"otherTitles"`
	Author            string    `json:"author"`
	OtherParticipate  []string  `json:"otherParticipate"`
	Genre             string    `json:"genre"`
	OtherGenres       []string  `json:"otherGenres"`
	Image             string    `json:"image"`
	Introduction      string    `json:"introduction"`
	Publisher         string    `json:"publisher"`
	FirstDateJP       string    `json:"firstDateJP"`
	FirstDateTH       string    `json:"firstDateTH"`
	LastVol           int       `json:"lastVol"`
	IsHighlight       bool      `json:"isHighlight"`
	SubscribersCount  int       `json:"subscribers"`
	Score             float64   `json:"score"`
	TotalVoters       int       `json:"totalVoters"`
}

type Vol struct {
	ID              uuid.UUID `json:"id"`
	MangaID         uuid.UUID `json:"mangaId"`
	VolNumber       int       `json:"vol"`
	Image           string    `json:"image"`
	PublishDate     string    `json:"publishDate"`
	TotalOwnerCount int       `json:"totalOwner"`
}
