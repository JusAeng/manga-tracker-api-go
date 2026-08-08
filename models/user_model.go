package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	LineUserID string    `json:"-"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	CreatedAt  time.Time `json:"createdAt"`
}
