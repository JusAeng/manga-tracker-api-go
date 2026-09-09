package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const mangaColumns = `
	id, title_original, title_en, introduction, image_url, first_date_jp,
	status, created_at, updated_at
`

func scanManga(row pgx.Row) (*models.Manga, error) {
	var m models.Manga
	err := row.Scan(
		&m.ID, &m.TitleOriginal, &m.TitleEN, &m.Introduction, &m.ImageURL,
		&m.FirstDateJP, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMangas lists manga, optionally filtered by a case-insensitive
// substring match against either title. Pass "" for no filter.
func GetMangas(search string) ([]*models.Manga, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT `+mangaColumns+` FROM manga
		WHERE $1 = '' OR title_original ILIKE '%' || $1 || '%' OR title_en ILIKE '%' || $1 || '%'
		ORDER BY title_original
	`, search)
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

func GetMangaById(id uuid.UUID) (*models.Manga, error) {
	row := db.Pool.QueryRow(context.Background(), `SELECT `+mangaColumns+` FROM manga WHERE id = $1`, id)
	m, err := scanManga(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func AddManga(manga *models.Manga) (*models.Manga, error) {
	row := db.Pool.QueryRow(context.Background(), `
		INSERT INTO manga (title_original, title_en, introduction, image_url, first_date_jp, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+mangaColumns,
		manga.TitleOriginal, manga.TitleEN, manga.Introduction, manga.ImageURL, manga.FirstDateJP, manga.Status,
	)
	return scanManga(row)
}

func UpdateManga(manga *models.Manga) (*models.Manga, error) {
	row := db.Pool.QueryRow(context.Background(), `
		UPDATE manga SET
			title_original = $1, title_en = $2, introduction = $3, image_url = $4,
			first_date_jp = $5, status = $6, updated_at = now()
		WHERE id = $7
		RETURNING `+mangaColumns,
		manga.TitleOriginal, manga.TitleEN, manga.Introduction, manga.ImageURL,
		manga.FirstDateJP, manga.Status, manga.ID,
	)
	return scanManga(row)
}

func DeleteMangaById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM manga WHERE id = $1`, id)
	return err
}
