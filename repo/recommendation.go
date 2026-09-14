package repo

import (
	"context"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
)

// GetUserRecommendedManga reads whatever the recommender service (Python,
// via Pub/Sub) last computed for this user, ordered by rank. Owns no
// write path here — that table is populated entirely by the recommender
// service, this is read-only. Empty (not an error) means nothing's been
// computed yet — the caller decides what to show instead.
func GetUserRecommendedManga(userId uuid.UUID) ([]*models.Manga, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT `+mangaColumns+`
		FROM manga
		JOIN recommendations r ON r.manga_id = manga.id
		WHERE r.user_id = $1
		ORDER BY r.rank
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
