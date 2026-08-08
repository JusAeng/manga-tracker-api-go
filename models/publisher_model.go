package models

import (
	"time"

	"github.com/google/uuid"
)

type Publisher struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	WebsiteURL string    `json:"websiteUrl"`
	LogoURL    string    `json:"logoUrl"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
