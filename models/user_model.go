package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	LineUserID  string    `json:"-"`
	DisplayName string    `json:"displayName"`
	PictureURL  string    `json:"pictureUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
