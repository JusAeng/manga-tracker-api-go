package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const userColumns = `id, line_user_id, display_name, picture_url, created_at, updated_at`

func scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.LineUserID, &u.DisplayName, &u.PictureURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetOrCreateUserByLineID looks up a user by their LINE sub, creating one
// (and refreshing display name/picture) if this is their first login.
// Upserting in a single statement avoids the find-then-insert race two
// concurrent first logins from the same account used to be able to hit.
func GetOrCreateUserByLineID(lineUserID, displayName, pictureURL string) (*models.User, error) {
	row := db.Pool.QueryRow(context.Background(), `
		INSERT INTO users (line_user_id, display_name, picture_url)
		VALUES ($1, $2, $3)
		ON CONFLICT (line_user_id) DO UPDATE SET
			display_name = EXCLUDED.display_name, picture_url = EXCLUDED.picture_url, updated_at = now()
		RETURNING `+userColumns,
		lineUserID, displayName, pictureURL,
	)
	return scanUser(row)
}

func GetUserProfileById(userId uuid.UUID) (*models.User, error) {
	row := db.Pool.QueryRow(context.Background(), `SELECT `+userColumns+` FROM users WHERE id = $1`, userId)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func GetAllUsers() ([]*models.User, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT `+userColumns+` FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func DeleteUserById(userId uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userId)
	return err
}

// UpdateUserProfile only ever touches "displayName" or "pictureUrl" — the
// caller (user_handlers.UpdateUserProfile) already allowlists the key, but
// SQL column names can't be bind parameters anyway, so this switch is the
// actual enforcement, not just a formality.
func UpdateUserProfile(userId uuid.UUID, key string, newValue string) error {
	var query string
	switch key {
	case "displayName":
		query = `UPDATE users SET display_name = $1, updated_at = now() WHERE id = $2`
	case "pictureUrl":
		query = `UPDATE users SET picture_url = $1, updated_at = now() WHERE id = $2`
	default:
		return errors.New("unsupported profile field")
	}
	_, err := db.Pool.Exec(context.Background(), query, newValue, userId)
	return err
}
