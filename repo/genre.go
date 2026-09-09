package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetGenres() ([]*models.Genre, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT id, name FROM genres ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]*models.Genre, 0)
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}
	return genres, rows.Err()
}

func GetGenreById(id uuid.UUID) (*models.Genre, error) {
	var g models.Genre
	err := db.Pool.QueryRow(context.Background(), `SELECT id, name FROM genres WHERE id = $1`, id).Scan(&g.ID, &g.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func AddGenre(g *models.Genre) (*models.Genre, error) {
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO genres (name) VALUES ($1) RETURNING id, name
	`, g.Name).Scan(&g.ID, &g.Name)
	return g, err
}

func UpdateGenre(g *models.Genre) (*models.Genre, error) {
	err := db.Pool.QueryRow(context.Background(), `
		UPDATE genres SET name = $1 WHERE id = $2 RETURNING id, name
	`, g.Name, g.ID).Scan(&g.ID, &g.Name)
	return g, err
}

func DeleteGenreById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM genres WHERE id = $1`, id)
	return err
}

// GetGenresByMangaId returns the manga's genres, for the manga-detail read.
func GetGenresByMangaId(mangaId uuid.UUID) ([]*models.Genre, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT g.id, g.name
		FROM manga_genres mg
		JOIN genres g ON g.id = mg.genre_id
		WHERE mg.manga_id = $1
		ORDER BY g.name
	`, mangaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]*models.Genre, 0)
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}
	return genres, rows.Err()
}

func AddMangaGenre(mangaId, genreId uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO manga_genres (manga_id, genre_id) VALUES ($1, $2)
		ON CONFLICT (manga_id, genre_id) DO NOTHING
	`, mangaId, genreId)
	return err
}

func RemoveMangaGenre(mangaId, genreId uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `
		DELETE FROM manga_genres WHERE manga_id = $1 AND genre_id = $2
	`, mangaId, genreId)
	return err
}
