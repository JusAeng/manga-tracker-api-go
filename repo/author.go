package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetAuthors() ([]*models.Author, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, name, created_at, updated_at FROM authors ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*models.Author
	for rows.Next() {
		var a models.Author
		if err := rows.Scan(&a.ID, &a.Name, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		authors = append(authors, &a)
	}
	return authors, rows.Err()
}

func GetAuthorById(id uuid.UUID) (*models.Author, error) {
	var a models.Author
	err := db.Pool.QueryRow(context.Background(), `
		SELECT id, name, created_at, updated_at FROM authors WHERE id = $1
	`, id).Scan(&a.ID, &a.Name, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func AddAuthor(a *models.Author) (*models.Author, error) {
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO authors (name) VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`, a.Name).Scan(&a.ID, &a.Name, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func UpdateAuthor(a *models.Author) (*models.Author, error) {
	err := db.Pool.QueryRow(context.Background(), `
		UPDATE authors SET name = $1, updated_at = now() WHERE id = $2
		RETURNING id, name, created_at, updated_at
	`, a.Name, a.ID).Scan(&a.ID, &a.Name, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func DeleteAuthorById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM authors WHERE id = $1`, id)
	return err
}

// GetAuthorsByMangaId returns the manga's credited authors with their
// role, for the manga-detail read.
func GetAuthorsByMangaId(mangaId uuid.UUID) ([]*models.MangaAuthor, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT a.id, a.name, ma.role
		FROM manga_authors ma
		JOIN authors a ON a.id = ma.author_id
		WHERE ma.manga_id = $1
		ORDER BY a.name
	`, mangaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*models.MangaAuthor
	for rows.Next() {
		var ma models.MangaAuthor
		if err := rows.Scan(&ma.AuthorID, &ma.Name, &ma.Role); err != nil {
			return nil, err
		}
		authors = append(authors, &ma)
	}
	return authors, rows.Err()
}

func AddMangaAuthor(mangaId, authorId uuid.UUID, role string) error {
	_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO manga_authors (manga_id, author_id, role) VALUES ($1, $2, $3)
		ON CONFLICT (manga_id, author_id, role) DO NOTHING
	`, mangaId, authorId, role)
	return err
}

func RemoveMangaAuthor(mangaId, authorId uuid.UUID, role string) error {
	_, err := db.Pool.Exec(context.Background(), `
		DELETE FROM manga_authors WHERE manga_id = $1 AND author_id = $2 AND role = $3
	`, mangaId, authorId, role)
	return err
}
