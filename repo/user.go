package repo

import (
	"context"
	"errors"

	"github.com/JusAeng/manga-tracker-api-go/db"
	"github.com/JusAeng/manga-tracker-api-go/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetOrCreateUserByLineID looks up a user by their LINE sub, creating one
// (and refreshing name/image) if this is their first login. Upserting in
// a single statement avoids the find-then-insert race two concurrent first
// logins from the same account used to be able to hit.
func GetOrCreateUserByLineID(lineUserID, name, image string) (*models.User, error) {
	var user models.User
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO users (line_user_id, name, image)
		VALUES ($1, $2, $3)
		ON CONFLICT (line_user_id) DO UPDATE SET name = EXCLUDED.name, image = EXCLUDED.image
		RETURNING id, line_user_id, name, image, created_at
	`, lineUserID, name, image).Scan(&user.ID, &user.LineUserID, &user.Name, &user.Image, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserProfileById(userId uuid.UUID) (*models.User, error) {
	var user models.User
	err := db.Pool.QueryRow(context.Background(), `
		SELECT id, line_user_id, name, image, created_at FROM users WHERE id = $1
	`, userId).Scan(&user.ID, &user.LineUserID, &user.Name, &user.Image, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]*models.User, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, line_user_id, name, image, created_at FROM users ORDER BY created_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.LineUserID, &user.Name, &user.Image, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

func DeleteUserById(userId uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userId)
	return err
}

// UpdateUserProfile only ever touches "name" or "image" — the caller
// (user_handlers.UpdateUserProfile) already allowlists the key, but SQL
// column names can't be bind parameters anyway, so this switch is the
// actual enforcement, not just a formality.
func UpdateUserProfile(userId uuid.UUID, key string, newValue string) error {
	var query string
	switch key {
	case "name":
		query = `UPDATE users SET name = $1 WHERE id = $2`
	case "image":
		query = `UPDATE users SET image = $1 WHERE id = $2`
	default:
		return errors.New("unsupported profile field")
	}
	_, err := db.Pool.Exec(context.Background(), query, newValue, userId)
	return err
}

// SubscribeMangaById toggles a subscription and keeps manga.subscribers_count
// in sync atomically in the same transaction — no read-modify-write race.
func SubscribeMangaById(userId uuid.UUID, mangaId uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `DELETE FROM subscriptions WHERE user_id = $1 AND manga_id = $2`, userId, mangaId)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE manga SET subscribers_count = subscribers_count - 1 WHERE id = $1`, mangaId); err != nil {
			return nil, err
		}
	} else {
		if _, err := tx.Exec(ctx, `INSERT INTO subscriptions (user_id, manga_id) VALUES ($1, $2)`, userId, mangaId); err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, `UPDATE manga SET subscribers_count = subscribers_count + 1 WHERE id = $1`, mangaId); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return GetUserSubscriptionIDs(userId)
}

func GetUserSubscriptionIDs(userId uuid.UUID) ([]uuid.UUID, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT manga_id FROM subscriptions WHERE user_id = $1`, userId)
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

// UpdateOwnerList toggles ownership of a specific volume. Looking the
// volume up by (manga_id, vol_number) first, then inserting owned_volumes
// against its real vol id, means the foreign key rejects ownership of a
// volume that doesn't exist — no more hand-rolled isMangaExist check.
func UpdateOwnerList(userId uuid.UUID, mangaId uuid.UUID, volNumber int) ([]int, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var volID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM vols WHERE manga_id = $1 AND vol_number = $2`, mangaId, volNumber).Scan(&volID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("no such volume")
		}
		return nil, err
	}

	tag, err := tx.Exec(ctx, `DELETE FROM owned_volumes WHERE user_id = $1 AND vol_id = $2`, userId, volID)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE vols SET total_owner_count = total_owner_count - 1 WHERE id = $1`, volID); err != nil {
			return nil, err
		}
	} else {
		if _, err := tx.Exec(ctx, `INSERT INTO owned_volumes (user_id, vol_id) VALUES ($1, $2)`, userId, volID); err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, `UPDATE vols SET total_owner_count = total_owner_count + 1 WHERE id = $1`, volID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return GetOwnedVolumeNumbers(userId, mangaId)
}

func GetOwnedVolumeNumbers(userId uuid.UUID, mangaId uuid.UUID) ([]int, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT v.vol_number FROM owned_volumes ov
		JOIN vols v ON v.id = ov.vol_id
		WHERE ov.user_id = $1 AND v.manga_id = $2
		ORDER BY v.vol_number
	`, userId, mangaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nums := []int{}
	for rows.Next() {
		var n int
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, rows.Err()
}

// UpdateRateList treats score == 0 as "remove my rating", matching the
// existing API contract. manga.score/total_voters are kept correct by the
// trg_ratings_refresh_stats trigger — no averaging math here anymore.
func UpdateRateList(userId uuid.UUID, mangaId uuid.UUID, score int) error {
	ctx := context.Background()
	if score == 0 {
		_, err := db.Pool.Exec(ctx, `DELETE FROM ratings WHERE user_id = $1 AND manga_id = $2`, userId, mangaId)
		return err
	}
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO ratings (user_id, manga_id, score, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (user_id, manga_id) DO UPDATE SET score = EXCLUDED.score, updated_at = now()
	`, userId, mangaId, score)
	return err
}
