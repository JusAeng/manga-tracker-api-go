package models

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	UserID    uuid.UUID `json:"userId"`
	MangaID   uuid.UUID `json:"mangaId"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RatingSummary is what both "read a manga's rating" and "rate/unrate a
// manga" endpoints return, so the frontend can update from either without
// a refetch.
type RatingSummary struct {
	AverageRating *float64 `json:"averageRating"`
	RatingCount   int      `json:"ratingCount"`
	MyRating      *int     `json:"myRating"`
}
