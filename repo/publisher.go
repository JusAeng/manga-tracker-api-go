package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const publisherColumns = `id, name, website_url, logo_url, created_at, updated_at`

func scanPublisher(row pgx.Row) (*models.Publisher, error) {
	var p models.Publisher
	err := row.Scan(&p.ID, &p.Name, &p.WebsiteURL, &p.LogoURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPublishers() ([]*models.Publisher, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT `+publisherColumns+` FROM publishers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publishers := make([]*models.Publisher, 0)
	for rows.Next() {
		p, err := scanPublisher(rows)
		if err != nil {
			return nil, err
		}
		publishers = append(publishers, p)
	}
	return publishers, rows.Err()
}

func GetPublisherById(id uuid.UUID) (*models.Publisher, error) {
	row := db.Pool.QueryRow(context.Background(), `SELECT `+publisherColumns+` FROM publishers WHERE id = $1`, id)
	p, err := scanPublisher(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func AddPublisher(p *models.Publisher) (*models.Publisher, error) {
	row := db.Pool.QueryRow(context.Background(), `
		INSERT INTO publishers (name, website_url, logo_url)
		VALUES ($1, $2, $3)
		RETURNING `+publisherColumns,
		p.Name, p.WebsiteURL, p.LogoURL,
	)
	return scanPublisher(row)
}

func UpdatePublisher(p *models.Publisher) (*models.Publisher, error) {
	row := db.Pool.QueryRow(context.Background(), `
		UPDATE publishers SET name = $1, website_url = $2, logo_url = $3, updated_at = now()
		WHERE id = $4
		RETURNING `+publisherColumns,
		p.Name, p.WebsiteURL, p.LogoURL, p.ID,
	)
	return scanPublisher(row)
}

func DeletePublisherById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM publishers WHERE id = $1`, id)
	return err
}
