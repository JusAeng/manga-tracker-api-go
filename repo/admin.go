package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/jackc/pgx/v5"
)

const adminColumns = `id, username, password_hash, created_at, updated_at`

// GetAdminByUsername returns nil (not an error) when no admin matches —
// callers should give the same "not found" response for a missing admin
// as for a wrong password, so this shape mirrors GetUserProfileById.
func GetAdminByUsername(username string) (*models.Admin, error) {
	var a models.Admin
	err := db.Pool.QueryRow(context.Background(), `
		SELECT `+adminColumns+` FROM admins WHERE username = $1
	`, username).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// CreateAdmin is used by scripts/seedadmin — there's no HTTP route for
// this, since creating the first admin has no authenticated caller yet.
func CreateAdmin(username, passwordHash string) (*models.Admin, error) {
	var a models.Admin
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO admins (username, password_hash)
		VALUES ($1, $2)
		RETURNING `+adminColumns,
		username, passwordHash,
	).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
