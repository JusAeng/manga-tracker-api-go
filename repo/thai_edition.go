package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const thaiEditionColumns = `id, manga_id, publisher_id, title_th, first_date_th, created_at, updated_at`

func scanThaiEdition(row pgx.Row) (*models.ThaiEdition, error) {
	var e models.ThaiEdition
	err := row.Scan(&e.ID, &e.MangaID, &e.PublisherID, &e.TitleTH, &e.FirstDateTH, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetThaiEditionsByMangaId(mangaId uuid.UUID) ([]*models.ThaiEdition, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT `+thaiEditionColumns+` FROM thai_editions WHERE manga_id = $1 ORDER BY first_date_th
	`, mangaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var editions []*models.ThaiEdition
	for rows.Next() {
		e, err := scanThaiEdition(rows)
		if err != nil {
			return nil, err
		}
		editions = append(editions, e)
	}
	return editions, rows.Err()
}

func GetThaiEditionById(id uuid.UUID) (*models.ThaiEdition, error) {
	row := db.Pool.QueryRow(context.Background(), `SELECT `+thaiEditionColumns+` FROM thai_editions WHERE id = $1`, id)
	e, err := scanThaiEdition(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

func AddThaiEdition(e *models.ThaiEdition) (*models.ThaiEdition, error) {
	row := db.Pool.QueryRow(context.Background(), `
		INSERT INTO thai_editions (manga_id, publisher_id, title_th, first_date_th)
		VALUES ($1, $2, $3, $4)
		RETURNING `+thaiEditionColumns,
		e.MangaID, e.PublisherID, e.TitleTH, e.FirstDateTH,
	)
	return scanThaiEdition(row)
}

func UpdateThaiEdition(e *models.ThaiEdition) (*models.ThaiEdition, error) {
	row := db.Pool.QueryRow(context.Background(), `
		UPDATE thai_editions SET publisher_id = $1, title_th = $2, first_date_th = $3, updated_at = now()
		WHERE id = $4
		RETURNING `+thaiEditionColumns,
		e.PublisherID, e.TitleTH, e.FirstDateTH, e.ID,
	)
	return scanThaiEdition(row)
}

func DeleteThaiEditionById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM thai_editions WHERE id = $1`, id)
	return err
}
