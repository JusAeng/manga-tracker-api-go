package models

import (
	"time"

	"github.com/google/uuid"
)

type Volume struct {
	ID            uuid.UUID  `json:"id"`
	ThaiEditionID uuid.UUID  `json:"thaiEditionId"`
	VolumeNumber  int        `json:"volumeNumber"`
	ISBN          *string    `json:"isbn"`
	PublishDate   *time.Time `json:"publishDate"`
	Price         *float64   `json:"price"`
	ImageURL      string     `json:"imageUrl"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// MangaVolumeUpdate is a narrow, purpose-built result shape for "which of
// this user's followed manga got a new volume recently" — not a table row,
// so it doesn't follow the id/timestamps shape the rest of this file does.
type MangaVolumeUpdate struct {
	MangaTitleEN       string
	MangaTitleOriginal string
	VolumeNumber       int
	PublishDate        time.Time
}
