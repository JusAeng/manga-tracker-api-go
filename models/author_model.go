package models

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MangaAuthor is a manga_authors row joined with the author's name, as
// returned by manga-detail reads — not a standalone CRUD entity.
type MangaAuthor struct {
	AuthorID uuid.UUID `json:"authorId"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
}
