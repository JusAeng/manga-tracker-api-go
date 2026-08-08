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
	id, title, other_titles, author, other_participate, genre, other_genres,
	image, introduction, publisher, first_date_jp, first_date_th, last_vol,
	is_highlight, subscribers_count, score, total_voters
`

func scanManga(row pgx.Row) (*models.Manga, error) {
	var m models.Manga
	err := row.Scan(
		&m.ID, &m.Title, &m.OtherTitles, &m.Author, &m.OtherParticipate, &m.Genre, &m.OtherGenres,
		&m.Image, &m.Introduction, &m.Publisher, &m.FirstDateJP, &m.FirstDateTH, &m.LastVol,
		&m.IsHighlight, &m.SubscribersCount, &m.Score, &m.TotalVoters,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Read
func GetMangas() ([]*models.Manga, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT `+mangaColumns+` FROM manga ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangas []*models.Manga
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

func GetMangaByTitle(title string) ([]*models.Manga, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT `+mangaColumns+` FROM manga WHERE title = $1`, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangas []*models.Manga
	for rows.Next() {
		m, err := scanManga(rows)
		if err != nil {
			return nil, err
		}
		mangas = append(mangas, m)
	}
	return mangas, rows.Err()
}

func GetHighlightManga() (*models.Manga, error) {
	row := db.Pool.QueryRow(context.Background(), `SELECT `+mangaColumns+` FROM manga WHERE is_highlight LIMIT 1`)
	m, err := scanManga(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func GetMangaFromSubscribeList(userId uuid.UUID) ([]*models.Manga, error) {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT `+mangaColumns+` FROM manga m
		JOIN subscriptions s ON s.manga_id = m.id
		WHERE s.user_id = $1
		ORDER BY m.title
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangas []*models.Manga
	for rows.Next() {
		m, err := scanManga(rows)
		if err != nil {
			return nil, err
		}
		mangas = append(mangas, m)
	}
	return mangas, rows.Err()
}

// Create
func AddManga(manga *models.Manga) (*models.Manga, error) {
	row := db.Pool.QueryRow(context.Background(), `
		INSERT INTO manga (title, other_titles, author, other_participate, genre, other_genres, image, introduction, publisher, first_date_jp, first_date_th)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING `+mangaColumns,
		manga.Title, manga.OtherTitles, manga.Author, manga.OtherParticipate, manga.Genre, manga.OtherGenres,
		manga.Image, manga.Introduction, manga.Publisher, manga.FirstDateJP, manga.FirstDateTH,
	)
	return scanManga(row)
}

// Delete
func DeleteMangaByTitle(title string) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM manga WHERE title = $1`, title)
	return err
}

func DeleteMangaById(id uuid.UUID) error {
	_, err := db.Pool.Exec(context.Background(), `DELETE FROM manga WHERE id = $1`, id)
	return err
}

// Update
func UpdateManga(manga *models.Manga) (*models.Manga, error) {
	row := db.Pool.QueryRow(context.Background(), `
		UPDATE manga SET
			title = $1, other_titles = $2, author = $3, other_participate = $4,
			genre = $5, other_genres = $6, image = $7, introduction = $8,
			publisher = $9, first_date_jp = $10, first_date_th = $11
		WHERE id = $12
		RETURNING `+mangaColumns,
		manga.Title, manga.OtherTitles, manga.Author, manga.OtherParticipate,
		manga.Genre, manga.OtherGenres, manga.Image, manga.Introduction,
		manga.Publisher, manga.FirstDateJP, manga.FirstDateTH, manga.ID,
	)
	return scanManga(row)
}

func SetMangaHighlight(mangaId uuid.UUID) error {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE manga SET is_highlight = false WHERE is_highlight`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE manga SET is_highlight = true WHERE id = $1`, mangaId); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
