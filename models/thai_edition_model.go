package models

import (
	"time"

	"github.com/google/uuid"
)

type ThaiEdition struct {
	ID          uuid.UUID  `json:"id"`
	MangaID     uuid.UUID  `json:"mangaId"`
	PublisherID uuid.UUID  `json:"publisherId"`
	TitleTH     string     `json:"titleTh"`
	FirstDateTH *time.Time `json:"firstDateTh"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
