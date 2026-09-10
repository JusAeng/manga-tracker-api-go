package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// RateManga upserts userId's rating for mangaId — one row per (user, manga),
// so re-rating just updates the existing value instead of stacking rows.
func RateManga(userId, mangaId uuid.UUID, rating int) error {
	_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO ratings (user_id, manga_id, rating)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, manga_id) DO UPDATE SET
			rating = EXCLUDED.rating, updated_at = now()
	`, userId, mangaId, rating)
	return err
}

func DeleteRating(userId, mangaId uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM ratings WHERE user_id = $1 AND manga_id = $2`, userId, mangaId)
	return err
}

// GetMangaRatingSummary computes the average/count on demand rather than
// reading a denormalized column — there isn't one, matching how `follows`
// keeps no counter either. userId may be the zero UUID (no authenticated
// user, e.g. an admin token) in which case myRating is always nil.
func GetMangaRatingSummary(mangaId uuid.UUID, userId uuid.UUID) (avg *float64, count int, myRating *int, err error) {
	err = db.Pool.QueryRow(context.Background(),
		`SELECT avg(rating), count(*) FROM ratings WHERE manga_id = $1`, mangaId,
	).Scan(&avg, &count)
	if err != nil {
		return nil, 0, nil, err
	}

	if userId == uuid.Nil {
		return avg, count, nil, nil
	}

	var mine int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT rating FROM ratings WHERE user_id = $1 AND manga_id = $2`, userId, mangaId,
	).Scan(&mine)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return avg, count, nil, nil
		}
		return nil, 0, nil, err
	}
	return avg, count, &mine, nil
}
