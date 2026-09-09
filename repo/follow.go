package repo

import (
	"context"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
)

// ToggleFollow flips whether userId follows mangaId. Unlike the old
// subscription toggle there's no denormalized counter to keep in sync, so
// this is a plain insert/delete instead of a transaction.
func ToggleFollow(userId uuid.UUID, mangaId uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	tag, err := db.Pool.Exec(ctx, `DELETE FROM follows WHERE user_id = $1 AND manga_id = $2`, userId, mangaId)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		if _, err := db.Pool.Exec(ctx, `INSERT INTO follows (user_id, manga_id) VALUES ($1, $2)`, userId, mangaId); err != nil {
			return nil, err
		}
	}
	return GetUserFollowedMangaIDs(userId)
}

func GetUserFollowedMangaIDs(userId uuid.UUID) ([]uuid.UUID, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT manga_id FROM follows WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func GetUserFollowedManga(userId uuid.UUID) ([]*models.Manga, error) {
	// mangaColumns is unqualified — fine for a plain SELECT FROM manga,
	// but follows also has a created_at column, so the join needs the
	// m. prefix explicitly or Postgres rejects it as ambiguous.
	rows, err := db.Pool.Query(context.Background(), `
		SELECT m.id, m.title_original, m.title_en, m.introduction, m.image_url,
			m.first_date_jp, m.status, m.created_at, m.updated_at
		FROM manga m
		JOIN follows f ON f.manga_id = m.id
		WHERE f.user_id = $1
		ORDER BY m.title_original
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mangas := make([]*models.Manga, 0)
	for rows.Next() {
		m, err := scanManga(rows)
		if err != nil {
			return nil, err
		}
		mangas = append(mangas, m)
	}
	return mangas, rows.Err()
}
