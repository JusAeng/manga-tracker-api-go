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
