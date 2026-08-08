package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func refreshLastVol(ctx context.Context, tx pgx.Tx, mangaId uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		UPDATE manga SET last_vol = COALESCE((SELECT MAX(vol_number) FROM vols WHERE manga_id = $1), 0)
		WHERE id = $1
	`, mangaId)
	return err
}

func AddMangaVol(vol models.Vol) (*models.Vol, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var created models.Vol
	err = tx.QueryRow(ctx, `
		INSERT INTO vols (manga_id, vol_number, image, publish_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, manga_id, vol_number, image, publish_date, total_owner_count
	`, vol.MangaID, vol.VolNumber, vol.Image, vol.PublishDate).Scan(
		&created.ID, &created.MangaID, &created.VolNumber, &created.Image, &created.PublishDate, &created.TotalOwnerCount,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, errors.New("already add this vol")
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign_key_violation
			return nil, errors.New("manga does not exist")
		}
		return nil, err
	}

	if err := refreshLastVol(ctx, tx, vol.MangaID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &created, nil
}

func UpdateVol(mangaId uuid.UUID, req models.Vol) error {
	tag, err := db.Pool.Exec(context.Background(), `
		UPDATE vols SET image = $1, publish_date = $2
		WHERE manga_id = $3 AND vol_number = $4
	`, req.Image, req.PublishDate, mangaId, req.VolNumber)
	if err != nil {
		return errors.New("can't update this vol")
	}
	if tag.RowsAffected() == 0 {
		return errors.New("vol not found")
	}
	return nil
}

func DeleteManyVols(mangaId uuid.UUID, volNumbers []int) ([]int, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		DELETE FROM vols WHERE manga_id = $1 AND vol_number = ANY($2)
	`, mangaId, volNumbers); err != nil {
		return nil, err
	}

	if err := refreshLastVol(ctx, tx, mangaId); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return volNumbers, nil
}

func DeleteAllVolsByMangaId(mangaId uuid.UUID) error {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM vols WHERE manga_id = $1`, mangaId); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE manga SET last_vol = 0 WHERE id = $1`, mangaId); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
