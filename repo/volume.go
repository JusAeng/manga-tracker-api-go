package repo

import (
	"context"
	"errors"
	"time"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const volumeColumns = `
	id, thai_edition_id, volume_number, isbn, publish_date, price, image_url, status, created_at, updated_at
`

func scanVolume(row pgx.Row) (*models.Volume, error) {
	var v models.Volume
	err := row.Scan(
		&v.ID, &v.ThaiEditionID, &v.VolumeNumber, &v.ISBN, &v.PublishDate,
		&v.Price, &v.ImageURL, &v.Status, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func GetVolumesByThaiEditionId(thaiEditionId uuid.UUID) ([]*models.Volume, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT `+volumeColumns+` FROM volumes WHERE thai_edition_id = $1 ORDER BY volume_number
	`, thaiEditionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	volumes := make([]*models.Volume, 0)
	for rows.Next() {
		v, err := scanVolume(rows)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, v)
	}
	return volumes, rows.Err()
}

// GetRecentVolumesForFollowedManga finds volumes published since `since`
// for manga the given user follows — user -> follows -> manga ->
// thai_editions -> volumes. `publish_date >= $2` naturally drops rows with
// a NULL publish_date (the column is nullable): a NULL comparison is never
// true, so Postgres excludes those rows from the WHERE clause on its own.
func GetRecentVolumesForFollowedManga(userId uuid.UUID, since time.Time) ([]*models.MangaVolumeUpdate, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT m.title_en, m.title_original, v.volume_number, v.publish_date
		FROM volumes v
		JOIN thai_editions te ON te.id = v.thai_edition_id
		JOIN manga m ON m.id = te.manga_id
		JOIN follows f ON f.manga_id = m.id
		WHERE f.user_id = $1 AND v.publish_date >= $2
		ORDER BY v.publish_date DESC
	`, userId, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updates := make([]*models.MangaVolumeUpdate, 0)
	for rows.Next() {
		var u models.MangaVolumeUpdate
		if err := rows.Scan(&u.MangaTitleEN, &u.MangaTitleOriginal, &u.VolumeNumber, &u.PublishDate); err != nil {
			return nil, err
		}
		updates = append(updates, &u)
	}
	return updates, rows.Err()
}

func AddVolume(v models.Volume) (*models.Volume, error) {
	row := db.Pool.QueryRow(context.Background(), `
		INSERT INTO volumes (thai_edition_id, volume_number, isbn, publish_date, price, image_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+volumeColumns,
		v.ThaiEditionID, v.VolumeNumber, v.ISBN, v.PublishDate, v.Price, v.ImageURL, v.Status,
	)
	created, err := scanVolume(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, errors.New("already added this volume")
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign_key_violation
			return nil, errors.New("thai edition does not exist")
		}
		return nil, err
	}
	return created, nil
}

// UpdateVolume and DeleteVolumeById address a volume by its own id, not by
// (thai_edition_id, volume_number) — the volumes table has a real primary
// key, so there's no reason to route updates through the parent's scope.
func UpdateVolume(id uuid.UUID, req models.Volume) (*models.Volume, error) {
	row := db.Pool.QueryRow(context.Background(), `
		UPDATE volumes SET
			volume_number = $1, isbn = $2, publish_date = $3, price = $4,
			image_url = $5, status = $6, updated_at = now()
		WHERE id = $7
		RETURNING `+volumeColumns,
		req.VolumeNumber, req.ISBN, req.PublishDate, req.Price, req.ImageURL, req.Status, id,
	)
	return scanVolume(row)
}

func DeleteVolumeById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM volumes WHERE id = $1`, id)
	return err
}
